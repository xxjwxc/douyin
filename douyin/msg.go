package douyin

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := string(body)
	re := regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	match := re.FindStringSubmatch(res)
	if match == nil || len(match) < 2 {
		return nil, err
	}

	douyinMsg.LiveRoomId = match[1]

	re = regexp.MustCompile(`live-room-nickname">(.*?)</div>`)
	match = re.FindStringSubmatch(res)
	if match != nil || len(match) > 1 {
		douyinMsg.LiveRoomTitle = match[1]
	}

	re = regexp.MustCompile(`\\"user_unique_id\\\":\\\"(\d+)\\\"}`)
	match = re.FindStringSubmatch(res)
	if match != nil || len(match) > 1 {
		douyinMsg.UserUniqueId = match[1]
	} else {
		douyinMsg.UserUniqueId = douyinMsg.LiveRoomId
	}

	return douyinMsg, nil
}

// 开始执行
func (d *DouyinMsg) Start(onMessage func(*MessageInfo)) error {
	// webSocketUrl := fmt.Sprintf(`wss://webcast3-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:%v|wss_push_did:7186637706334963260&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=%v&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=%v&signature=00000000`,
	// 	d.LiveRoomId, d.UserUniqueId, d.LiveRoomId)
	webSocketUrl := fmt.Sprintf(`wss://webcast3-ws-web-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:%v|wss_push_did:%v|dim_log_id:202302171547011A160A7BAA76660E13ED|fetch_time:1676620021641|seq:1|wss_info:0-1676620021641-0-0|wrds_kvs:WebcastRoomStatsMessage-1676620020691146024_WebcastRoomRankMessage-1676619972726895075_AudienceGiftSyncData-1676619980834317696_HighlightContainerSyncData-2&cursor=t-1676620021641_r-1_d-1_u-1_h-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id='+user_unique_id+'&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%%20(Macintosh;%%20Intel%%20Mac%%20OS%%20X%%2010_15_7)%%20AppleWebKit/537.36%%20(KHTML,%%20like%%20Gecko)%%20Chrome/110.0.0.0%%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=%v&heartbeatDuration=0&signature=00000000`,
		d.LiveRoomId, d.UserUniqueId, d.LiveRoomId)

	head := http.Header{
		"cookie":     {fmt.Sprintf("ttwid=%v", d.Ttwid)},
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
	}

	var err error
	d.Wss, err = mywebsocket.NewWebSocket(webSocketUrl, func(messageType int, p []byte, err error) {
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
		buf, err := io.ReadAll(zr)
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
			onMessage(&MessageInfo{
				Method:  msg.Method,
				Payload: msg.Payload,
			})
		}
	}, head, 30*time.Second)

	if err == nil {
		go func() {
			for {
				pingPack := &PushFrame{
					PayloadType: "bh",
				}
				data, _ := proto.Marshal(pingPack)
				err := d.Wss.SendMessage(websocket.BinaryMessage, data)
				if err != nil {
					break
				}
				// log.Println("发送心跳")
				time.Sleep(time.Second * 10)
			}
		}()
	}
	return err

}

// 停止
func (d *DouyinMsg) Stop() {
	if d.Wss != nil {
		d.Wss.Close()
		d.Wss = nil
	}
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
	d.Wss.SendMessage(websocket.BinaryMessage, b)
	//mylog.Info(`[sendAck] [🌟发送Ack] [房间Id：' + liveRoomId + '] ====> 房间🏖标题【' + liveRoomTitle + '】`)
	return nil
}
