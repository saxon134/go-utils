# go一些特性封装

##### Pool
本地协程控制，并发控制，不支持多务器多实例控制


##### TokenBucket
令牌桶算法，支持配置Redis，支持多务器多实例控制

##### Limiter
限流算法，支持配置Redis，支持多务器多实例控制
支持限定调用最大、最小间隔，防止并发处理资源

`LimiterTryLock(key, minSecond, LimiterGlobalOption, maxSecond)` 的数值选项
指定 Redis 租约秒数，默认 600 秒。正数向上取整到整秒，例如 `3` 为 3 秒，
`0.5` 为 1 秒，`1800` 为 1800 秒。`LimiterLock` 的 `maxSecond` 使用相同规则。
非正数、无法转换为数字、非有限值或超出 `time.Duration` 秒数范围的租约参数被忽略；
多个有效数值选项以最后一个为准。

动态 key 可以选择以下回收方式：

- 在资源生命周期结束后调用 `LimiterForget(key)`。没有持锁者、等待者或正在获取锁的调用时，
  删除成功并返回 `true`；仍在使用或 key 不存在时返回 `false`。
  此操作会清除该 key 的本地冷却历史，不会删除 Redis 锁，不应作为每次解锁的替代操作。
- 在任何 `LimiterLock`、`LimiterTryLock` 或间接使用它们的调用（例如 `BucketConsume`）之前，
  调用一次 `ConfigureLimiterCleanup(retention)` 启用自动回收。`retention` 必须为正数，
  同时作为所有调用允许的最大 `minSecond`。配置失败返回错误。
  上限按现有的毫秒精度校验，例如 `100*time.Millisecond` 允许 `minSecond=0.1`。
  `LimiterTryLock` 对超限或非有限间隔返回 `false`；`LimiterLock` 对这些参数 panic，
  不会在未加锁时继续执行。有限负数间隔仍按 0 处理。

```go
// 在启动业务 goroutine 前配置；所有 minSecond 必须不超过 3600 秒。
if err := saGo.ConfigureLimiterCleanup(time.Hour); err != nil {
    panic(err)
}
```

自动回收只删除没有使用者、且连续空闲达到 `retention` 的本地条目。
后台每 `retention / 2` 扫描一次，扫描间隔限制在 1 毫秒到 1 分钟之间。
已经获取条目、但尚在等待本地锁的调用也计入使用者，避免回收后同一 key 出现两把锁。
配置只允许在首次使用前成功设置一次，不支持运行中扩大窗口，以免已删除的冷却历史被重新需要。
默认不自动回收，以保留现有 `minSecond` 无上限、每次调用可变化的语义。
