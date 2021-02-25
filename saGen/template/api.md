# go-businessGen

#### 简介
生成一般controller & business接口处理函数




#### 使用方法

1. 下载代码到本地项目

2. 在gen.go中，修改调用生成器方法的参数

    ```
    func main() {
    	err := generate(
    		TblDemo{},
    		"techio"
    	)
    	if err != nil {
    		fmt.Println(err.Error())
    	}
    }
    ```

    ##### 说明：

    1. 数据模型
    2. root路由


3. 生成的代码在output目录下

    ##### 说明：
    1. controller，为API接口函数
    2. business，具体接口实现

