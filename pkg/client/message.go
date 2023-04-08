package client

import (
	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/sirupsen/logrus"
)

func (c *Client) handleMessage(d []byte) {
	md, err := data.NewDanmakuMessageData(d)
	if err != nil {
		logrus.Error(err)
	}
	switch md.CMD {
	case data.DMK_DANMU_MSG:
		logrus.Debugf("DANMU_MSG")
		logrus.Debugf("%v", string(d))
	case data.DMK_INTERACT_WORD:
		logrus.Debugf("INTERACT_WORD")
		logrus.Debugf("%v", string(d))
	case data.DMK_ONLINE_RANK_COUNT:
		logrus.Debugf("ONLINE_RANK_COUNT")
		logrus.Debugf("%v", string(d))
	case data.DMK_ONLINE_RANK_V2:
		logrus.Debugf("ONLINE_RANK_V2")
		logrus.Debugf("%v", string(d))
	case data.DMK_WATCHED_CHANGE:
		logrus.Debugf("WATCHED_CHANGE")
		logrus.Debugf("%v", string(d))
	case data.DMK_ENTRY_EFFECT:
	case data.DMK_LIKE_INFO_V3_CLICK:
	case data.DMK_STOP_LIVE_ROOM_LIST:
	case data.DMK_LIKE_INFO_V3_UPDATE:
	default:
		logrus.Debugf("Unrecognized: %v", md.CMD)
	}
}
