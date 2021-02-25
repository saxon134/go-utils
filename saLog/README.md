# yfLog

#### 介绍
本地日志


- 有info、warn、err三个等级
- 日志较多时，自动从info -> warn -> err降级
- 日志打印通过chan实现
- 目前支持fmt、zap工具