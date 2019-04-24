import ujson
from autobahn.twisted.websocket import WebSocketClientProtocol, \
    WebSocketClientFactory


class MyClientProtocol(WebSocketClientProtocol):

    def onConnect(self, response):
        print("Server connected: {0}".format(response.peer))

    def onConnecting(self, transport_details):
        print("Connecting; transport details: {}".format(transport_details))
        return None  # ask for defaults

    def onOpen(self):
        print("WebSocket connection open.")
        self.send()

    def onMessage(self, payload, isBinary):
        payload = ujson.loads(payload)
        print("Text message received: {0}".format(payload['Name']))

    def _send(self):
        payload = {
            'method' : 'v0.env.launch',
            'body' : {
                'env_id' : 'wob.mini.TicTacToe'
            }
        }
        self.sendMessage(ujson.dumps(payload).encode('utf-8'), False)
        self.send()
    
    def send(self):
        self.factory.reactor.callFromThread(self._send)

    def onClose(self, wasClean, code, reason):
        print("WebSocket connection closed: {0}".format(reason))


if __name__ == '__main__':

    import sys

    from twisted.python import log
    from twisted.internet import reactor

    log.startLogging(sys.stdout)

    factory = WebSocketClientFactory(u"ws://127.0.0.1:15900")
    factory.protocol = MyClientProtocol

    reactor.connectTCP("127.0.0.1", 15900, factory)
    reactor.run()