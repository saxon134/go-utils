package saWebhook

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saWebhook/ding"
	"github.com/saxon134/go-utils/saWebhook/feishu"
)

type Channel int

const (
	FeiShu Channel = 1
	Ding   Channel = 2
)

type Txt struct {
	Channel Channel
	Title   string `json:"title"`
	Msg     string `json:"msg"`
	Webhook string `json:"webhook"`
	Secret  string `json:"secret"`
}

func (m *Txt) Send() (err error) {
	if m.Msg == "" {
		return saError.New("msg cannot be empty")
	}

	if m.Channel == FeiShu {
		if m.Title != "" {
			return feishu.New(m.Webhook, m.Secret).SendTxtWithTitle(m.Title, m.Msg)
		} else {
			return feishu.New(m.Webhook, m.Secret).SendTxt(m.Msg)
		}
	} else if m.Channel == Ding {
		return ding.SendTxt(m.Title, m.Msg, m.Webhook)
	}

	return saError.New("channel cannot be empty")
}
