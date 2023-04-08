package data

import (
	"encoding/json"

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
		Uid:       0,
		Roomid:    roomID,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.14.3",
		Type:      2,
	}
	d, err := json.Marshal(p)
	if err != nil {
		logrus.Warn(err)
		return nil
	}
	return d
}

// 00 00 00 14 00 10 00 01
// 00 00 00 03 00 00 00 00
// 00 00 00 0c

//5b 6f 62 6a 65 63 74 20 4f 62 6a 65
// 63 74 5d
