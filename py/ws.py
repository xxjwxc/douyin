# pip install websocket-server
from websocket_server import WebsocketServer
import dy
import logging
import json
 
map  = {}
 
# 当新的客户端连接时会提示
def new_client(client, server):
    douyin = dy.DouYing("",server,client)
    map[client['id']] = douyin
    logging.info('[当新的客户端连接时会提示：' + json.dumps(client['id']) + ']')
 
 
# 当旧的客户端离开
def client_left(client, server):
    logging.info('[客户端断开:' + json.dumps(client['id']) + ']')
    if client['id'] in map:
        map[client['id']].stop()
        map.pop(client['id'])
    logging.info('[客户端断开:')
 
 
# 接收客户端的信息。
def message_received(client, server, message):
    logging.info('[收到客户端消息:' + json.dumps(client['id']) + '][消息内容:'+message+']')
    if client['id'] in map:
        map[client['id']].dealMsg(message)
 
 
def Start(port):
    LOG_FORMAT = "%(asctime)s - %(levelname)s - %(message)s"
    logging.basicConfig(level=logging.DEBUG, format=LOG_FORMAT)
    logging.getLogger().setLevel(logging.DEBUG)
    logging.info('[ws start on 0.0.0.0:%d]',port)
    server = WebsocketServer( "0.0.0.0",port)
    server.set_fn_new_client(new_client)
    server.set_fn_client_left(client_left)
    server.set_fn_message_received(message_received)
    server.run_forever()
