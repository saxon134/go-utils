package saTrigger

import (
	"encoding/json"
	"github.com/saxon134/go-utils/saBroker"
)

func Remote(name string, params interface{}) error {
	bAry, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = saBroker.Manager.Switch(name).SetParams(bAry).Do()
	return err
}

func Local(job *saBroker.LocalJob) error {
	saBroker.LocalDo(job)
	return nil
}
