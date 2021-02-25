package saGen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func generate1(obj interface{}, rootRoute string) error {
	refectType := reflect.TypeOf(obj)

	pkgName := ""
	titlePkgName := ""
	modelName := ""
	//获取包名、结构体名
	{
		str := refectType.String()
		ary := strings.Split(str, ".")
		if ary != nil && len(ary) >= 2 {
			pkgName = ary[0]
			titlePkgName = strings.Title(pkgName)
			modelName = ary[1]
		}
	}

	//生成controller
	{
		read_f, err := os.OpenFile("template/controller.tpl", os.O_RDONLY, 0600)
		if err != nil {
			return err
		}

		b, _ := ioutil.ReadAll(read_f)
		_ = read_f.Close()

		content := string(b)
		content = strings.Replace(content, "{{TitlePkgName}}", titlePkgName, -1)
		content = strings.Replace(content, "{{PkgName}}", pkgName, -1)
		content = strings.Replace(content, "{{ModelName}}", modelName, -1)
		content = strings.Replace(content, "{{FunDetail}}", titlePkgName+"Detail", -1)
		content = strings.Replace(content, "{{FunAdd}}", titlePkgName+"Add", -1)
		content = strings.Replace(content, "{{FunUpdate}}", titlePkgName+"Update", -1)
		content = strings.Replace(content, "{{FunUpdateStatus}}", titlePkgName+"UpdateStatus", -1)
		content = strings.Replace(content, "{{FunList}}", titlePkgName+"List", -1)
		content = strings.Replace(content, "{{BsPkgName}}", "bs"+titlePkgName, -1)

		data := []byte(content)
		f_n := "output/controller." + pkgName + ".go"
		if ioutil.WriteFile(f_n, data, 0644) != nil {
			return errors.New("出错")
		}
	}

	//生成route配置代码
	confStr := `
	"RootRoute/PkgName.add":           {Method: Post, Check: saHttp.MsCheck, Handle: controller.TitlePkgNameAdd},
	"RootRoute/PkgName.update":        {Method: Post, Check: saHttp.MsCheck, Handle: controller.TitlePkgNameUpdate},
	"RootRoute/PkgName.update.status": {Method: Post, Check: saHttp.MsCheck, Handle: controller.TitlePkgNameUpdateStatus},
	"RootRoute/PkgName.list":          {Method: Get, Check: saHttp.NullCheck, Handle: controller.TitlePkgNameList},
	"RootRoute/PkgName":               {Method: Get, Check: saHttp.NullCheck, Handle: controller.TitlePkgNameDetail},
	`
	confStr = strings.Replace(confStr, "RootRoute", rootRoute, -1)
	confStr = strings.Replace(confStr, "TitlePkgName", titlePkgName, -1)
	confStr = strings.Replace(confStr, "PkgName", pkgName, -1)
	fmt.Println("路由配置：")
	fmt.Println(confStr)

	return nil
}
