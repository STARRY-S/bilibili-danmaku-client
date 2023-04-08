package data

import (
	"bytes"
	"encoding/json"
)

type DanmakuMessageType string

const (
	DMK_DANMU_MSG           = DanmakuMessageType("DANMU_MSG")         // 弹幕消息
	DMK_ONLINE_RANK_COUNT   = DanmakuMessageType("ONLINE_RANK_COUNT") // 高能榜 (正在观看的人数)
	DMK_ONLINE_RANK_V2      = DanmakuMessageType("ONLINE_RANK_V2")    // 高能榜排名
	DMK_INTERACT_WORD       = DanmakuMessageType("INTERACT_WORD")     // 进房消息
	DMK_WATCHED_CHANGE      = DanmakuMessageType("WATCHED_CHANGE")    // 看过的人数
	DMK_ENTRY_EFFECT        = DanmakuMessageType("ENTRY_EFFECT")
	DMK_LIKE_INFO_V3_CLICK  = DanmakuMessageType("LIKE_INFO_V3_CLICK")
	DMK_LIKE_INFO_V3_UPDATE = DanmakuMessageType("LIKE_INFO_V3_UPDATE")
	DMK_STOP_LIVE_ROOM_LIST = DanmakuMessageType("STOP_LIVE_ROOM_LIST")
)

type DanmakuMessageData struct {
	CMD  DanmakuMessageType `json:"cmd,omitempty"`
	Info []interface{}      `json:"info,omitempty"`
	Data interface{}        `json:"data,omitempty"`
}

func NewDanmakuMessageData(d []byte) (*DanmakuMessageData, error) {
	m := &DanmakuMessageData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
