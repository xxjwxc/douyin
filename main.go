package main

import (
	"douyin/douyin"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
)

func main() {

	tmp, err := douyin.NewRoom("391122936429")
	tmp.Start(onMessage)

	time.Sleep(1 * time.Minute)
	tmp.Stop()
	fmt.Println(tmp, err)
}

func onMessage(msg *douyin.MessageInfo) {
	switch msg.Method {
	case "WebcastMatchAgainstScoreMessage":
		// fmt.Println(msg.Payload)
	case "WebcastLikeMessage":
		var likeMsg douyin.LikeMessage
		_ = proto.Unmarshal(msg.Payload, &likeMsg)
		log.Printf("[点赞] %s 点赞 * %d \n", likeMsg.User.NickName, likeMsg.Count)
	case "WebcastMemberMessage":
		var enterMsg douyin.MemberMessage
		_ = proto.Unmarshal(msg.Payload, &enterMsg)
		log.Printf("[入场] %s 直播间\n", enterMsg.User.NickName)
	case "WebcastGiftMessage":
		var giftMsg douyin.GiftMessage
		_ = proto.Unmarshal(msg.Payload, &giftMsg)
		log.Printf("[礼物] %s : %s * %d \n", giftMsg.User.NickName, giftMsg.Gift.Name, giftMsg.ComboCount)
	case "WebcastChatMessage":
		var chatMsg douyin.ChatMessage
		_ = proto.Unmarshal(msg.Payload, &chatMsg)
		log.Printf("[弹幕] %s : %s\n", chatMsg.User.NickName, chatMsg.Content)
	case "WebcastSocialMessage":
		// var chatMsg douyin.SocialMessage
		// _ = proto.Unmarshal(msg.Payload, &chatMsg)
		// log.Printf("[社交信息] %s : %s\n", chatMsg.User.NickName, chatMsg.ShareTarget)
	case "WebcastRoomUserSeqMessage":
		// var chatMsg douyin.RoomUserSeqMessage
		// _ = proto.Unmarshal(msg.Payload, &chatMsg)
		// log.Printf("[直播间信息] 在线人数: %v  TotalPvForAnchor : %v\n", chatMsg.Total, chatMsg.TotalPvForAnchor)
	case "WebcastUpdateFanTicketMessage":
		// var chatMsg douyin.UpdateFanTicketMessage
		// _ = proto.Unmarshal(msg.Payload, &chatMsg)
		// log.Printf("[粉丝通知] 粉丝数: %v  : %v\n", chatMsg.RoomFanTicketCount, chatMsg.RoomFanTicketCountText)
	case "WebcastCommonTextMessage":
		// fmt.Println(msg.Payload)
	}
}
