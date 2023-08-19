package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DanmakuURL    = "wss://broadcastlv.chat.bilibili.com/sub"
	DanmakuOrigin = "https://live.bilibili.com"
)

func GetHttpData(url string, hd map[string]string) (*http.Response, error) {
	logrus.Debugf("GetHttpData: %v", url)
	client := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	if hd != nil {
		for k, v := range hd {
			req.Header.Set(k, v)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client Get: %w", err)
	}
	return res, nil
}
