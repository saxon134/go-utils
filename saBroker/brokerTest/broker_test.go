package brokerTest

import (
	"fmt"
	"gitee.com/go-utils/saBroker"
	"gitee.com/go-utils/saBroker/saTrigger"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	//初始化broker
	{
		var m *saBroker.BrokerManager
		m = saBroker.Init("redis://127.0.0.1:6379", "brokerTest")
		if m == nil {
			fmt.Println("broker init error.")
			return
		}

		err := yfBroker.RegisterRemoteJobs(
			NewPrintJob("print_test"),
		)
		if err != nil {
			fmt.Println("broker jobs init error.")
			return
		}

		_ = yfBroker.RegisterLocalJobs(10, func(j *yfBroker.LocalJob) {
			fmt.Println(j.Type, j.Value)
		})

		fmt.Println("broker init ok.")
	}

	//server层，发送broker消息
	{
		err := yfTrigger.Remote("print_test", "123abc")
		if err != nil {
			fmt.Println("trigger error:", err)
			return
		}

		_ = yfTrigger.Local(&yfBroker.LocalJob{Type: "1", Value: "a"})

		fmt.Println("trigger ok!")
	}

	time.Sleep(time.Second * 10)
}
