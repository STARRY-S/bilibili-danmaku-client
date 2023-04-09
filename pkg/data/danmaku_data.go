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
	DMK_LIKE_INFO_V3_CLICK  = DanmakuMessageType("LIKE_INFO_V3_CLICK")  // 用户点赞
	DMK_LIKE_INFO_V3_UPDATE = DanmakuMessageType("LIKE_INFO_V3_UPDATE") // 总点赞数
	DMK_SEND_GIFT           = DanmakuMessageType("SEND_GIFT")           // 送礼物
	DMK_NOTICE_MSG          = DanmakuMessageType("NOTICE_MSG")
	DMK_STOP_LIVE_ROOM_LIST = DanmakuMessageType("STOP_LIVE_ROOM_LIST")
)

type CmdData struct {
	CMD DanmakuMessageType `json:"cmd,omitempty"`
}

type MessageData struct {
	CmdData
	Info []interface{} `json:"info,omitempty"`
	Data interface{}   `json:"data,omitempty"`
}

func NewMessageData(d []byte) (*MessageData, error) {
	m := &MessageData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type DanmuMsgData struct {
	CmdData
	Info []interface{} `json:"info,omitempty"`
}

func NewDanmuMsgData(d []byte) (*DanmuMsgData, error) {
	m := &DanmuMsgData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (dmk *DanmuMsgData) GetUsername() string {
	switch c := dmk.Info[2].(type) {
	case []interface{}:
		switch u := c[1].(type) {
		case string:
			return u
		}
	}
	return ""
}

func (dmk *DanmuMsgData) GetUID() int {
	switch c := dmk.Info[2].(type) {
	case []interface{}:
		switch i := c[0].(type) {
		case float64:
			return int(i)
		}
	}
	return 0
}

func (dmk *DanmuMsgData) GetContent() string {
	switch c := dmk.Info[1].(type) {
	case string:
		return c
	}
	return ""
}

type InteractWordData struct {
	CmdData

	Data struct {
		UID   int    `json:"uid,omitempty"`
		Uname string `json:"uname,omitempty"`
	} `json:"data,omitempty"`
}

func NewInteractWordData(d []byte) (*InteractWordData, error) {
	m := &InteractWordData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type OnlineRankCountData struct {
	CmdData

	Data struct {
		Count int `json:"count,omitempty"`
	}
}

func NewOnlineRankCountData(d []byte) (*OnlineRankCountData, error) {
	m := &OnlineRankCountData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type OnlineRankV2Data struct {
	CmdData

	Data struct {
		List []struct {
			Face  string `json:"face,omitempty"`
			Rank  int    `json:"rank,omitempty"`
			Score string `json:"score,omitempty"`
			Uid   int    `json:"uid,omitempty"`
			Uname string `json:"uname,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

func NewOnlineRankV2Data(d []byte) (*OnlineRankV2Data, error) {
	m := &OnlineRankV2Data{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type WatchedChangeData struct {
	CmdData

	Data struct {
		Num       int    `json:"num,omitempty"`
		TextLarge string `json:"text_large,omitempty"`
		TextSmall string `json:"text_small,omitempty"`
	}
}

func NewWatchedChangeData(d []byte) (*WatchedChangeData, error) {
	m := &WatchedChangeData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type LikeInfoV3ClickData struct {
	CmdData

	Data struct {
		UID      int    `json:"uid,omitempty"`
		Uname    string `json:"uname,omitempty"`
		LikeText string `json:"like_text,omitempty"`
		LikeIcon string `json:"like_icon,omitempty"`
	} `json:"data,omitempty"`
}

func NewLikeInfoV3ClickData(d []byte) (*LikeInfoV3ClickData, error) {
	m := &LikeInfoV3ClickData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type LikeInfoV3Update struct {
	CmdData

	Data struct {
		ClickCount int `json:"click_count,omitempty"`
	} `json:"data,omitempty"`
}

func NewLikeInfoV3Update(d []byte) (*LikeInfoV3Update, error) {
	m := &LikeInfoV3Update{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type SendGiftData struct {
	CmdData

	Data struct {
		Action   string `json:"action,omitempty"`
		Face     string `json:"face,omitempty"`
		GiftName string `json:"giftName,omitempty"`
		UID      int    `json:"uid,omitempty"`
		Uname    string `json:"uname,omitempty"`
		CoinType string `json:"coin_type,omitempty"` // silver / gold
		Num      int    `json:"num,omitempty"`
	} `json:"data,omitempty"`
}

func NewSendGiftData(d []byte) (*SendGiftData, error) {
	m := &SendGiftData{}
	jd := json.NewDecoder(bytes.NewReader(d))
	err := jd.Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
