package msg

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/mywebsocket"
)

// NewRoom 新创建一个房间,(房间id通过抖音官网浏览器中获取:如https://live.douyin.com/533646158504，roomId为533646158504)
func NewRoom(roomId string) (*DouyinMsg, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://live.douyin.com/%v", roomId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Set("cookie", "__ac_nonce=0638733a400869171be51")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	douyinMsg := &DouyinMsg{}
	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "ttwid" {
			douyinMsg.Ttwid = c.Value
			break
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := string(body)
	var routeRegex = regexp.MustCompile(`<script id="RENDER_DATA" type="application/json">(.*?)</script>`)
	matches := routeRegex.FindStringSubmatch(res)
	if len(matches) < 1 {
		return nil, message.GetError(message.NotFindError)
	}

	enEscapeUrl, err := url.QueryUnescape(matches[1]) // url解码
	if err != nil {
		return nil, err
	}

	var tmp douyinResp
	err = json.Unmarshal([]byte(enEscapeUrl), &tmp)
	if err != nil {
		return nil, err
	}
	douyinMsg.LiveRoomId = tmp.App.InitialState.RoomStore.RoomInfo.RoomId
	douyinMsg.LiveRoomTitle = tmp.App.InitialState.RoomStore.RoomInfo.Room.Title
	douyinMsg.UserUniqueId = tmp.App.Odin.UserUniqueId

	return douyinMsg, nil
}

// 开始执行
func (d *DouyinMsg) Start(onMessage func(*MessageInfo)) error {
	// webSocketUrl := fmt.Sprintf(`wss://webcast3-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:%v|wss_push_did:7186637706334963260&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=%v&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=%v&signature=00000000`,
	// 	d.LiveRoomId, d.UserUniqueId, d.LiveRoomId)
	webSocketUrl := fmt.Sprintf(`/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:%v|wss_push_did:7186637706334963260&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=%v&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=%v&signature=00000000`,
		d.LiveRoomId, d.UserUniqueId, d.LiveRoomId)

	u := url.URL{Scheme: "wss", Host: "webcast3-ws-web-lf.douyin.com:443", Path: webSocketUrl}
	// "wss://webcast3-ws-web-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:'+liveRoomId+'|wss_push_did:'+liveRoomId+'|dim_log_id:202302171547011A160A7BAA76660E13ED|fetch_time:1676620021641|seq:1|wss_info:0-1676620021641-0-0|wrds_kvs:WebcastRoomStatsMessage-1676620020691146024_WebcastRoomRankMessage-1676619972726895075_AudienceGiftSyncData-1676619980834317696_HighlightContainerSyncData-2&cursor=t-1676620021641_r-1_d-1_u-1_h-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id='+liveRoomId+'&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/110.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=" + d.LiveRoomId + "&heartbeatDuration=0&signature=Rk7kMWh+wzXKrKP2"
	head := http.Header{
		"cookie":     {fmt.Sprintf("ttwid=%v", d.Ttwid)},
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
	}

	var err error
	d.Wss, err = mywebsocket.NewWebSocket(u.String(), func(messageType int, p []byte, err error) {
		var wssPackage PushFrame
		err = proto.Unmarshal(p, &wssPackage)
		if err != nil {
			mylog.Error(err)
			return
		}
		logId := wssPackage.LogId

		zr, err := gzip.NewReader(bytes.NewReader(wssPackage.Payload))
		if err != nil {
			mylog.Error(err)
			return
		}
		defer zr.Close()
		// 读取gzip对象内容
		buf, err := ioutil.ReadAll(zr)
		if err != nil {
			fmt.Println("[read gzip data err]: ", err)
		}

		var payloadPackage Response
		err = proto.Unmarshal(buf, &payloadPackage)
		if err != nil {
			mylog.Error(err)
			return
		}
		if payloadPackage.NeedAck { // 发送ack包
			d.sendAck(logId, payloadPackage.InternalExt)
		}
		for _, msg := range payloadPackage.MessagesList {
			switch msg.Method {
			case "WebcastMatchAgainstScoreMessage":
				fmt.Println(msg.Payload)
			case "WebcastLikeMessage":
				fmt.Println(msg.Payload)
			case "WebcastMemberMessage":
				fmt.Println(msg.Payload)
			case "WebcastGiftMessage":
				fmt.Println(msg.Payload)
			case "WebcastChatMessage":
				fmt.Println(msg.Payload)
			case "WebcastSocialMessage":
				fmt.Println(msg.Payload)
			case "WebcastRoomUserSeqMessage":
				fmt.Println(msg.Payload)
			case "WebcastUpdateFanTicketMessage":
				fmt.Println(msg.Payload)
			case "WebcastCommonTextMessage":
				fmt.Println(msg.Payload)
			}
		}
	}, head, 30*time.Second)
	return err

}

// 开始执行
func (d *DouyinMsg) sendAck(logId uint64, internalExt string) error {
	obj := PushFrame{}
	obj.PayloadType = "ack"
	obj.LogId = logId
	obj.PayloadType = internalExt
	b, err := proto.Marshal(&obj)
	if err != nil {
		return err
	}
	d.Wss.SendMessage(b)
	mylog.Info(`[sendAck] [🌟发送Ack] [房间Id：' + liveRoomId + '] ====> 房间🏖标题【' + liveRoomTitle + '】`)
	return nil
}
