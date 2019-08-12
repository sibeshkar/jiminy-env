import ujson
import time
from autobahn.twisted.websocket import WebSocketClientProtocol, \
    WebSocketClientFactory

class RewarderProtocol(WebSocketClientProtocol):

    def onConnect(self, response):
        print("Server connected: {0}".format(response.peer))

    def onConnecting(self, transport_details):
        print("Connecting; transport details: {}".format(transport_details))
        return None  # ask for defaults

    def onOpen(self):
        print("WebSocket connection open.")
        self.send()

    def onMessage(self, payload, isBinary):
        assert not isBinary
        payload = ujson.loads(payload)
        context = self._make_context()
        #latency = context['start'] - payload['headers']['sent_at']
        self.recv(context, payload)
        print("Text message received: {0}".format(payload['method']))

    def _send(self):
        pass
        payload_launch = {
            'method' : 'v0.env.launch',
            'body' : {
                'env_id' : 'sibeshkar/wob-v1',
                'fps' : 30,
                'record': True,
            }
        }
        self.sendMessage(ujson.dumps(payload_launch).encode('utf-8'), False)

        payload_reset = {
            'method' : 'v0.env.reset',
            'body' : {
                'env_id' : 'sibeshkar/wob-v1/ClickButton'
            }
        }

        self.sendMessage(ujson.dumps(payload_reset).encode('utf-8'), False)
        #self.send()
    
    def send(self):
        self.factory.reactor.callFromThread(self._send)

    def onClose(self, wasClean, code, reason):
        print("WebSocket connection closed: {0}".format(reason))

    def send_reset(self, env_id, seed, fps, episode_id):
        pass

    def _finish_reset(self, episode_id):
        pass

    def waitForWebsocketConnection(self):
        pass

    def _manual_recv(self, method, body, headers={}):
        pass

    def recv(self, context, response):
        method = response['method']
        body = response['body']
        headers = response['headers']

        remote_time = headers['sent_at']
        local_time = context['start']
        print("Method: {}, Body: {}, Headers: {}, Remote time: {}".format(method, body, headers, remote_time))
        pass
    
    def _make_context(self):
        return {'start': time.time()}

    def close(self, code=1000, reason=None):
        pass

    

class RewarderSession(object):
    def __init__(self):
        pass

    def close(self, name=None, reason=u'closed by RewarderSession.close'):
        pass
    
    def connect(self, name, address, label, password, env_id=None, seed=None, fps=60,
                start_timeout=None, observer=False, skip_network_calibration=False):
                pass
    
    def _already_closed(self, i):
        pass

    def _connect(self, name, address, env_id, seed, fps, i, network, env_status, reward_buffer,
                 label, password, start_timeout,
                 observer, skip_network_calibration,
                 attempt=0, elapsed_sleep_time=0,):
                 pass
    
    def pop_errors(self):
        pass

    def reset(self, seed=None, env_id=None):
        pass

    def _reset(self, seed=None, env_id=None):
        pass

    def _send_env_reset(self, client, seed=None, episode_id=None, env_id=None):
        pass

    def pop(self, warn=True, peek_d=None):
        pass

    def wait(self, timeout=None):
        pass

    def send_action(self, action_n, env_id):
        pass

    def _send_action(self, env_id, action_n):
        pass

    def _send_env_action(self, client, env_id, action_n):
        pass

    def rewards_count(self):
        pass

    def pop_observation(self):
        pass

class Network(object):
    def __init__(self):
        pass

    def active(self):
        pass

    def reversed_clock_skew(self):
        pass

    def _report(self):
        pass

    def _start(self):
        pass

    def close(self):
        pass

    def calibrate(self, client):
        pass

    def _start_measure_connection_time(self, d):
        pass

    def _measure_connection_time(self, d, connection_time_m, i):
        pass
    
    def _start_measure_application_ping(self, d=None):
        pass

    def _measure_application_ping(self, d, clock_skew_m, request_overhead_m, response_overhead_m, application_rtt_m, i):
        pass

    def _update_exposed_metrics(self):
        pass
    
    def _start_measure_clock_skew(self):
        pass



if __name__ == '__main__':

    import sys

    from twisted.python import log
    from twisted.internet import reactor

    log.startLogging(sys.stdout)

    factory = WebSocketClientFactory(u"ws://127.0.0.1:15901")
    factory.protocol = RewarderProtocol

    reactor.connectTCP("127.0.0.1", 15901, factory)
    reactor.run()