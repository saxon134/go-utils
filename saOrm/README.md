# saOrm

#### 介绍

- 支持的tag标签：
  
  index  ->   单字段索引，生成参考建表语句会建立索引
  oss    ->   oss文件，入库、出库前域名、oss type处理；
  varchar(32) ->  字符最大长度为32
  &lt; &lt;= &lt; &gt;= &gt;&lt; ->  大小校验，如果是数组，则校验数组长度；字符串，校验rune长度
  required
  enum(1:激活,2:废弃)
  in(1,2)  in(ms,md)
  int8 int int64 string varchar(128) char(23) index created updated
  decimal(1,5)
  comment
  default
  phone
  updated   -> 时间，更新的时候会自动取当前时间，仅支持时间或者时间指针对象
  created   -> 时间，创建的时候会自动设置当前时间，仅支持时间或者时间指针对象
  
- tag示例：

  type:"varchar(32);created;comment:状态 2-正常 1-取消了;default:12;phone;required;in(ms,md);comment:字段;<=23;phone;required;default;err:缺少参数"

- GenTblSql：

    生成数据库操作SQL，如果表存在则生成建表语句；如果表存在，则生成修改表结构语句

- DB：

    通过反射，根据标签校验数据格式，及oss相关处理

    支持大小校验、oss处理
