# 数据处理

##### sa_check
- 支持的tag标签：
  > >= < <= <> : 大小校验。字符串，校验rune长度；整形则对比数值
  required
  enum(1:激活,2:废弃)
  in(1,2)  in(ms,md)
  phone

- tag示例：
  type:"phone;required;in(ms,md);<=23"

- TODO:
  考虑通过反射，生成校验函数，缓存到内存，下次反射可以直接调用

