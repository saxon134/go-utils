package brokerTest

import (
	"fmt"
	"github.com/saxon134/go-utils/saBroker"
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

		err := saBroker.RegisterJobs(
			NewPrintJob("print_test"),
		)
		if err != nil {
			fmt.Println("broker jobs init error.")
			return
		}

		fmt.Println("broker init ok.")
	}

	//server层，发送broker消息
	{
		err := saBroker.Do("print_test", "123abc")
		if err != nil {
			fmt.Println("trigger broker error:", err)
			return
		}

		fmt.Println("trigger broker ok!")
	}

	time.Sleep(time.Second * 10)
}
