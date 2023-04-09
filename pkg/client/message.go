package client

import (
	"fmt"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/sirupsen/logrus"
)

func (c *Client) handleMessage(d []byte) {
	md, err := data.NewMessageData(d)
	if err != nil {
		logrus.Errorf("handleMessage: %v", err)
		return
	}
	switch md.CMD {
	case data.DMK_DANMU_MSG:
		dmk, err := data.NewDanmuMsgData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_DANMU_MSG: %v", err)
			break
		}
		logrus.Debugf("Danmaku uid: %v", dmk.GetUID())
		logrus.Infof("用户 %q 说: %q", dmk.GetUsername(), dmk.GetContent())
	case data.DMK_INTERACT_WORD:
		msg, err := data.NewInteractWordData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_INTERACT_WORD: %v", err)
			break
		}
		logrus.Debugf("Interact word uid %v", msg.Data.UID)
		logrus.Infof("欢迎用户 %q 进入直播间", msg.Data.Uname)
	case data.DMK_ONLINE_RANK_COUNT:
		msg, err := data.NewOnlineRankCountData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_ONLINE_RANK_COUNT: %v", err)
			break
		}
		logrus.Infof("高能榜用户数: %v", msg.Data.Count)
	case data.DMK_ONLINE_RANK_V2:
		msg, err := data.NewOnlineRankV2Data(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_ONLINE_RANK_V2: %v", err)
			break
		}
		logrus.Infof("高能榜用户:")
		for _, v := range msg.Data.List {
			fmt.Printf("\tRank: %v, Name: %q, Score: %v\n",
				v.Rank, v.Uname, v.Score)
		}
	case data.DMK_WATCHED_CHANGE:
		msg, err := data.NewWatchedChangeData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_WATCHED_CHANGE: %v", err)
			break
		}
		logrus.Infof("%s", msg.Data.TextLarge)
	case data.DMK_ENTRY_EFFECT:
	case data.DMK_LIKE_INFO_V3_CLICK:
		msg, err := data.NewLikeInfoV3ClickData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_LIKE_INFO_V3_CLICK: %v", err)
			break
		}
		logrus.Infof("%q %s", msg.Data.Uname, msg.Data.LikeText)
	case data.DMK_LIKE_INFO_V3_UPDATE:
		msg, err := data.NewLikeInfoV3Update(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_LIKE_INFO_V3_UPDATE: %v", err)
			break
		}
		logrus.Infof("点赞数: %v", msg.Data.ClickCount)
	case data.DMK_SEND_GIFT:
		msg, err := data.NewSendGiftData(d)
		if err != nil {
			logrus.Errorf("handleMessage: DMK_SEND_GIFT: %v", err)
			break
		}
		logrus.Infof("感谢 %q 赠送的 %d 个%s",
			msg.Data.Uname, msg.Data.Num, msg.Data.GiftName)
	case data.DMK_STOP_LIVE_ROOM_LIST:
	case data.DMK_NOTICE_MSG:
	default:
		logrus.Debugf("Unrecognized message type: %v", md.CMD)
	}
}
