package data

import (
	"encoding/json"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/config"
	"github.com/sirupsen/logrus"
)

type firstData struct {
	Uid       int    `json:"uid"`
	Roomid    int    `json:"roomid"`
	Protover  int    `json:"protover"`
	Platform  string `json:"platform"`
	Clientver string `json:"clientver"`
	Type      int    `json:"type"`
}

func GetFirstData(roomID int) []byte {
	p := &firstData{
		Uid:      config.GetInt("uid"),
		Roomid:   roomID,
		Protover: 3,
		Platform: "web",
		Type:     2,
	}
	d, err := json.Marshal(p)
	if err != nil {
		logrus.Warn(err)
		return nil
	}
	return d
}
