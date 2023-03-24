import logging
from websocket_server import WebsocketServer

import dy
import ws
import sys

if __name__ == '__main__':
    # sys.setCharacterEncoding('utf-8') #set default encoding to utf-8
    LOG_FORMAT = "%(asctime)s - %(levelname)s - %(message)s"
    logging.basicConfig(level=logging.DEBUG, format=LOG_FORMAT)
    logger = logging.getLogger(__name__)
    # url = 'https://live.douyin.com/759420337103'
    # douyin = dy.DouYing(url,"server","client")
    # douyin.start()

    ws.Start(8131)
    
    # dy.parseLiveRoomUrl(url)