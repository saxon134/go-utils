# 数据处理

##### 注意事项
- 开发过程中，saData包不能引用go-utils任何包


##### sa_check
- 支持的tag标签：

  "> >= < <= <>" 长度校验：字符串，校验rune长度；整形则对比数值<br/>
  "required" 必要参数<br/>
  "enum(1:激活,2:废弃) in(1,2)  in(ms,md)" 枚举<br/>
  "phone" 校验手机格式

- tag示例：
  type:"phone;required;in(ms,md);<=23"

