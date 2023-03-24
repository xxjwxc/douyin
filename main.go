package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/mywebsocket"
)

func main() {
	ws, err := mywebsocket.NewWebSocket("ws://192.155.1.28:8131/", func(messageType int, p []byte, err error) {
		fmt.Println(string(p))
	}, http.Header{}, 30*time.Second)
	if err != nil {
		mylog.Error(err)
		return
	}
	msg := `{"url":"118272276289","proxyIp":""}`
	ws.SendMessage([]byte(msg))
	// wss://webcast3-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:7210729785163582266|wss_push_did:7186637706334963260|dim_log_id:20230315200805A3BB172EA72788121E4F|fetch_time:1678882085886|seq:1|wss_info:0-1678882085886-0-0|wrds_kvs:WebcastRoomRankMessage-1678882083541268494_WebcastRoomStatsMessage-1678882083487411248&cursor=d-1_u-1_h-1_t-1678882085886_r-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7186637706334963260&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows%20NT%2010.0;%20Win64;%20x64)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/111.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7210729785163582266&heartbeatDuration=0&signature=RBDXc0uX0UKp7CqB
	// msg.NewWsClientManager("webcast3-ws-web-lf.douyin.com", "443", `/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:7210729785163582266|wss_push_did:7186637706334963260|dim_log_id:20230315200805A3BB172EA72788121E4F|fetch_time:1678882085886|seq:1|wss_info:0-1678882085886-0-0|wrds_kvs:WebcastRoomRankMessage-1678882083541268494_WebcastRoomStatsMessage-1678882083487411248&cursor=d-1_u-1_h-1_t-1678882085886_r-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7186637706334963260&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows%20NT%2010.0;%20Win64;%20x64)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/111.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7210729785163582266&heartbeatDuration=0&signature=RBDXc0uX0UKp7CqB`, 10)
	// time.Sleep(10 * time.Minute)
	// // wsc.dail()

	// tmp, err := msg.NewRoom("630233272021")
	// tmp.Start(onMessage)

	time.Sleep(10 * time.Minute)
}

// func onMessage(msg *msg.MessageInfo) {
// 	fmt.Println(msg)
// }
