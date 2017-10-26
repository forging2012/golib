package ding

import (
	"io/ioutil"
	"time"
	"encoding/json"
	"bytes"
	"net/http"
	"github.com/verystar/golib/logger"
)

var DING_TALK_TOKEN string

type DingTalkRequest struct {
	Msgtype string                 `json:"msgtype"`
	Text    map[string]string      `json:"text"`
	At      map[string]interface{} `json:"at"`
}

type DingTalkResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func Alarm(content string, at ...[]string) error {
	dingtalk := &DingTalkRequest{
		Msgtype: "text",
		Text: map[string]string{
			"content": content,
		},
	}

	if len(at) > 0 {
		dingtalk.At = map[string]interface{}{
			"atMobiles": at,
		}
	}

	buf, err := json.Marshal(dingtalk)
	if err != nil {
		return err
	}

	url := "https://oapi.dingtalk.com/robot/send?access_token=" + DING_TALK_TOKEN
	resp, err := PostJson(url, buf)
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		dingtalk.paserResp(data)
	}

	if err != nil {
		logger.Debug("alarm", err.Error())
		return err
	}
	defer resp.Body.Close()
	return nil
}

func PostJson(url string, data []byte) (*http.Response, error) {
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (this *DingTalkRequest) paserResp(data []byte) {
	ret := new(DingTalkResponse)
	err := json.Unmarshal(data, ret)
	if err != nil || ret.Errcode != 0 {
		logger.Debug("dingtalk", string(data))
		return
	}
}
