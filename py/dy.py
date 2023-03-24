import _thread
import binascii
import gzip
import json
import logging
import re
import time
import requests
import websocket
import urllib
import threading
from protobuf_inspector.types import StandardParser
from google.protobuf import json_format
from dy_pb2 import PushFrame
from dy_pb2 import Response
from dy_pb2 import MatchAgainstScoreMessage
from dy_pb2 import LikeMessage
from dy_pb2 import MemberMessage
from dy_pb2 import GiftMessage
from dy_pb2 import ChatMessage
from dy_pb2 import SocialMessage
from dy_pb2 import RoomUserSeqMessage
from dy_pb2 import UpdateFanTicketMessage
from dy_pb2 import CommonTextMessage


class DouYing:
    def __init__(this, url, server,client):
        LOG_FORMAT = "%(asctime)s - %(levelname)s - %(message)s"
        logging.basicConfig(level=logging.DEBUG, format=LOG_FORMAT)
        logging.getLogger().setLevel(logging.DEBUG)
        this.url = url
        this.server = server
        this.client = client
        this.isStop = True

    def start(this):
        this.parseLiveRoomUrl(this.url)
    
    def stop(this):
        if this.isStop == False:
            this.ws.keep_running = False
            this.ws.close()
            this.isStop = True
    def dealMsg(this,text):
        if this.isStop == True:
            msg = json.loads(text)
            if len(msg['url']) > 0:
                this.url = 'https://live.douyin.com/' + msg['url']
                this.start()

                logging.info('[开始启动:' + this.url + '] ')

    def onMessage(this,ws: websocket.WebSocketApp, message: bytes):
        wssPackage = PushFrame()
        wssPackage.ParseFromString(message)
        logId = wssPackage.logId
        decompressed = gzip.decompress(wssPackage.payload)
        payloadPackage = Response()
        payloadPackage.ParseFromString(decompressed)
        # 发送ack包
        if payloadPackage.needAck:
            this.sendAck(ws, logId, payloadPackage.internalExt)
        for msg in payloadPackage.messagesList:
            if msg.method == 'WebcastMatchAgainstScoreMessage':
                this.unPackMatchAgainstScoreMessage(msg.payload)
                continue

            if msg.method == 'WebcastLikeMessage':
                this.unPackWebcastLikeMessage(msg.payload)
                continue

            if msg.method == 'WebcastMemberMessage':
                this.unPackWebcastMemberMessage(msg.payload)
                continue
            if msg.method == 'WebcastGiftMessage':
                this.unPackWebcastGiftMessage(msg.payload)
                continue
            if msg.method == 'WebcastChatMessage':
                this.unPackWebcastChatMessage(msg.payload)
                continue

            if msg.method == 'WebcastSocialMessage':
               this.unPackWebcastSocialMessage(msg.payload)
               continue

            if msg.method == 'WebcastRoomUserSeqMessage':
                this.unPackWebcastRoomUserSeqMessage(msg.payload)
                continue

            if msg.method == 'WebcastUpdateFanTicketMessage':
                this.unPackWebcastUpdateFanTicketMessage(msg.payload)
                continue

            if msg.method == 'WebcastCommonTextMessage':
                this.unPackWebcastCommonTextMessage(msg.payload)
                continue

        #logging.info('[onMessage] [⌛️方法' + msg.method + '等待解析～] [房间Id：' + liveRoomId + ']')


    def unPackWebcastCommonTextMessage(this,data):
        commonTextMessage = CommonTextMessage()
        commonTextMessage.ParseFromString(data)
        data = json_format.MessageToDict(commonTextMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "CommonText"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['scene']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastCommonTextMessage] [] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    def unPackWebcastUpdateFanTicketMessage(this,data):
        updateFanTicketMessage = UpdateFanTicketMessage()
        updateFanTicketMessage.ParseFromString(data)
        data = json_format.MessageToDict(updateFanTicketMessage, preserving_proto_field_name=True)
        
        msg  = {}
        msg['type'] = "UpdateFanTicket"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = ""
        msg['data'] = data['roomFanTicketCountText']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastUpdateFanTicketMessage] [] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    def unPackWebcastRoomUserSeqMessage(this,data):
        roomUserSeqMessage = RoomUserSeqMessage()
        roomUserSeqMessage.ParseFromString(data)
        data = json_format.MessageToDict(roomUserSeqMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "RoomUserSeq"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = ""
        msg['data'] = data['totalUserStr']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastRoomUserSeqMessage] [] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    def unPackWebcastSocialMessage(this,data):
        socialMessage = SocialMessage()
        socialMessage.ParseFromString(data)
        data = json_format.MessageToDict(socialMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "Social"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['shareTarget']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastSocialMessage] [➕直播间关注消息] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    # 普通消息
    def unPackWebcastChatMessage(this,data):
        chatMessage = ChatMessage()
        chatMessage.ParseFromString(data)
        data = json_format.MessageToDict(chatMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "Chat"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['content']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # logging.info('[unPackWebcastChatMessage] [📧直播间弹幕消息] [房间Id：' + liveRoomId + '] ｜ ' + data['content'])
        # logging.info('[unPackWebcastChatMessage] [📧直播间弹幕消息] [房间Id：' + liveRoomId + '] ｜ ' + json.dumps(data))
        return data


    # 礼物消息
    def unPackWebcastGiftMessage(this,data):
        giftMessage = GiftMessage()
        giftMessage.ParseFromString(data)
        data = json_format.MessageToDict(giftMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "Gift"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['gift']['name']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastGiftMessage] [🎁直播间礼物消息] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    # xx成员进入直播间消息
    def unPackWebcastMemberMessage(this,data):
        memberMessage = MemberMessage()
        memberMessage.ParseFromString(data)
        data = json_format.MessageToDict(memberMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "Membe"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['anchorDisplayText']['defaultPatter']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastMemberMessage] [🚹🚺直播间成员加入消息] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    # 点赞
    def unPackWebcastLikeMessage(this,data):
        likeMessage = LikeMessage()
        likeMessage.ParseFromString(data)
        data = json_format.MessageToDict(likeMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "Like"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = data['user']['nickName']
        msg['data'] = data['count']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackWebcastLikeMessage] [👍直播间点赞消息] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    # 解析WebcastMatchAgainstScoreMessage消息包体
    def unPackMatchAgainstScoreMessage(this,data):
        matchAgainstScoreMessage = MatchAgainstScoreMessage()
        matchAgainstScoreMessage.ParseFromString(data)
        data = json_format.MessageToDict(matchAgainstScoreMessage, preserving_proto_field_name=True)
        msg  = {}
        msg['type'] = "MatchAgainstScore"
        # msg['room_name'] = this.liveRoomTitle
        msg['user_name'] = ""
        msg['data'] = data['against']['rightName']
        this.server.send_message(this.client,json.dumps(msg, ensure_ascii=False))
        # log = json.dumps(data, ensure_ascii=False)
        # logging.info('[unPackMatchAgainstScoreMessage] [🤷不知道是啥的消息] [房间Id：' + liveRoomId + '] ｜ ' + log)
        return data


    # 发送Ack请求
    def sendAck(this,ws, logId, internalExt):
        obj = PushFrame()
        obj.payloadType = 'ack'
        obj.logId = logId
        obj.payloadType = internalExt
        data = obj.SerializeToString()
        ws.send(data, websocket.ABNF.OPCODE_BINARY)
    # logging.info('[sendAck] [🌟发送Ack] [房间Id：' + liveRoomId + '] ====> 房间🏖标题【' + liveRoomTitle + '】')


    def onError(this,ws, error):
        logging.error('[onError] [webSocket Error事件] [房间Id：' + liveRoomId + ']')


    def onClose(this,ws, a, b):
        logging.info('[onClose] [webSocket Close事件] [房间Id：' + liveRoomId + ']')


    def onOpen(this,ws):
        _thread.start_new_thread(this.ping, (ws,))
        logging.info('[onOpen] [webSocket Open事件] [房间Id：' + liveRoomId + ']')


    # 发送ping心跳包
    def ping(this,ws):
        while this.isStop == False:
            obj = PushFrame()
            obj.payloadType = 'hb'
            data = obj.SerializeToString()
            ws.send(data, websocket.ABNF.OPCODE_BINARY)
            
            logging.info('[ping] [💗发送ping心跳] [房间Id：' + liveRoomId + '] ====> 房间🏖标题【' + liveRoomTitle + '】')
            time.sleep(10)


    def wssServerStart(this,roomId,user_unique_id):
        global liveRoomId
        liveRoomId = roomId
        websocket.enableTrace(False)
        webSocketUrl = 'wss://webcast3-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:'+roomId+'|wss_push_did:'+user_unique_id+'&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id='+user_unique_id+'&device_platform=web&cookie_enabled=true&screen_width=2560&screen_height=1440&browser_language=zh-CN&browser_platform=Win32&browser_name=Mozilla&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id='+roomId+'&signature=00000000'
        h = {
            'cookie': "ttwid=" + ttwid,
            'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36',
        }
        # 创建一个长连接
        this.ws = websocket.WebSocketApp(
            webSocketUrl, on_message=this.onMessage, on_error=this.onError, on_close=this.onClose,
            on_open=this.onOpen,
            header=h
        )
        this.isStop = False
        this.ws.keep_running = True 
        this.wst = threading.Thread(target=this.ws.run_forever)
        this.wst.daemon = True
        this.wst.start()
        # this.ws.run_forever()


    def parseLiveRoomUrl(this,url):
        h = {
            'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36',
            'cookie': '__ac_nonce=0638733a400869171be51',
        }
        res = requests.get(url=url, headers=h)
        global ttwid, roomStore, liveRoomId, liveRoomTitle,user_unique_id
        data = res.cookies.get_dict()
        ttwid = data['ttwid']
        res = res.text
        res = re.search(r'<script id="RENDER_DATA" type="application/json">(.*?)</script>', res)
        res = res.group(1)
        res = urllib.parse.unquote(res, encoding='utf-8', errors='replace')
        res = json.loads(res)
        roomStore = res['app']['initialState']['roomStore']
        liveRoomId = roomStore['roomInfo']['roomId']
        liveRoomTitle = roomStore['roomInfo']['room']['title']
        user_unique_id = res['app']['odin']['user_unique_id']
        this.wssServerStart(liveRoomId,user_unique_id)