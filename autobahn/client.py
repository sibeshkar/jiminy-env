import ujson
import time
from autobahn.twisted.websocket import WebSocketClientProtocol, \
    WebSocketClientFactory

from twisted.internet import defer

from jiminy.rewarder import reward_buffer

class RewarderProtocol(WebSocketClientProtocol):

    def __init__(self):
        super(RewarderProtocol, self).__init__()
        self._closed = False
        self._close_message = None

        self._connected = False
        self._requests = {}

        self._reset = None
        self._launch = None
        self._initial_reset = False
        self._initial_launch = False
        self.label = "random-connection-1234"

        self._connection_result = defer.Deferred()
        
    def send_reset(self, env_id, seed, fps, episode_id):
        self._initial_reset = True
        self._reset = {
            'env_id': env_id,
            #'fps': fps,
            'episode_id': episode_id,
        }

        print("resetting")

        return self.send('v0.env.reset', {
            #'seed': seed,
            'env_id': env_id,
            #'fps': fps,
        }, {'episode_id': episode_id}, expect_reply=False)

    def send_launch(self, env_id, seed, fps, episode_id):
        self._initial_launch = True
        self._launch = {
            'env_id': env_id,
            #'fps': fps,
            'episode_id': episode_id,
        }


        return self.send('v0.env.launch', {
            #'seed': seed,
            'env_id': env_id,
            #'fps': fps,
        }, {'episode_id': episode_id}, expect_reply=False)

    def onConnect(self, response):
        print("Server connected: {0}".format(response.peer))
        self._message_id = 0
        self.reward_buffer = reward_buffer.RewardBuffer(self.label)

    def onConnecting(self, transport_details):
        print("Connecting; transport details: {}".format(transport_details))
        return None  # ask for defaults

    def onOpen(self):
        print("WebSocket connection open.")
        self.send_launch("sibeshkar/wob-v0", None, None, 1 )
        print("launched")
        self.send_reset("sibeshkar/wob-v0/TicTacToe", None, None, 1)

    def onMessage(self, payload, isBinary):
        assert not isBinary
        payload = ujson.loads(payload)
        context = self._make_context()
        #latency = context['start'] - payload['headers']['sent_at']
        self.recv(context, payload)
        print("Text message received: {0}".format(payload['method']))

    def _send(self, method, body, headers=None, expect_reply=False):
        if headers is None:
            headers = {}

        id = self._message_id

        self._message_id += 1
        new_headers = {
            'message_id': id,
            'sent_at': time.time(),
        }
        new_headers.update(headers)

        payload = {
            'method': method,
            'body': body,
        }

        #print(payload)

        # #payload= {
        #     'method' : 'v0.env.launch',
        #     'body' : {
        #         'env_id' : 'sibeshkar/wob-v0',
        #     }
        # }
        #self.sendMessage(ujson.dumps(payload).encode('utf-8'), False)

        # payload_reset = {
        #     'method' : 'v0.env.reset',
            
        #     'body' : {
        #         'seed': 56,
        #         'env_id' : 'sibeshkar/wob-v0/ClickShades',
        #         'fps': 60
        #     },
        #     'headers' : {
        #         'message_id' : '0',
        #         'sent_at': 1560007971.97964,
        #         'episode_id': 0
        #     },
        # }

        payload_reset = {
            'method': 'v0.env.reset', 
            'body': {'seed': 56, 'env_id': 'sibeshkar/wob-v0/ClickShades', 'fps': 60}, 
            'headers': {'message_id': 0, 'sent_at': 1560081689.9166744, 'episode_id': 0}
            }

        #payload_reset_2 = {'method': 'v0.env.reset', 'body': {'seed': 56, 'env_id': 'sibeshkar/wob-v0/ClickShades', 'fps': 60}, 'headers': {'message_id': 0, 'sent_at': 1560082557.4949555, 'episode_id': 0}}


        #{'method': 'v0.env.reset', 'body': {'seed': None, 'env_id': 'sibeshkar/wob-v0', 'fps': 60}, 'headers': {'message_id': 0, 'sent_at': 1560007971.97964, 'episode_id': '0'}}

        #self.sendMessage(ujson.dumps(payload).encode('utf-8'), False)
        self.sendMessage(ujson.dumps(payload_reset).encode('utf-8'), False)
        #self.send()
        # if expect_reply:
        #     d = defer.Deferred()
        #     self._requests[id] = (payload, d)
        #     return d
        # else:
        #     return None
    
    def send(self, method, body, headers=None, expect_reply=False):
         self.factory.reactor.callFromThread(self._send, method, body, headers=None, expect_reply=False)

    def onClose(self, wasClean, code, reason):
        print("WebSocket connection closed: {0}".format(reason))

    def send_reset(self, env_id, seed, fps, episode_id):
        pass

    def _finish_reset(self, episode_id):
        print('[%s] Running finish_reset: %s', self.label, episode_id)
        #self.reward_buffer.reset(episode_id)

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

        if method == 'v0.env.reward':
            episode_id = headers['episode_id']
            reward = body['reward']
            done = body['done']
            info = body['info']
            print('[%s] Received %s: reward=%s done=%s info=%s episode_id=%s', self.label, method, reward, done, info, episode_id)
            #pyprofile.incr('rewarder_client.reward', reward)
            # if done:
            #     pyprofile.incr('rewarder_client.done')
            #self.reward_buffer.push(episode_id, reward, done, info)
        elif method == 'v0.env.observation':
            episode_id = headers['episode_id']
            jsonable = body['observation']
            print('[%s] Received %s: observation=%s episode_id=%s', self.label, method, jsonable, episode_id)
            #self.reward_buffer.set_observation(episode_id=episode_id, observation=jsonable)
        elif method == 'v0.env.describe':
            episode_id = headers['episode_id']
            env_id = body['env_id']
            env_state = body['env_status'] ##note this change, original has ['env_state']
            fps = body['fps']
            print('[%s] Received %s: env_id=%s env_state=%s episode_id=%s',
                              self.label, method, env_id, env_state, episode_id)
            #self.reward_buffer.set_env_info(env_state, env_id=env_id, episode_id=episode_id, fps=fps)
        elif method == 'v0.reply.env.reset':
            episode_id = headers['episode_id']
            self._finish_reset(episode_id)
        elif method in ['v0.reply.error', 'v0.reply.control.ping']:
            assert headers.get('parent_message_id') is not None
        else:
            print('Unrecognized websocket method: method=%s body=%s headers=%s (consider adding to rewarder_state.py)', method, body, headers)
            return

    
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

    factory = WebSocketClientFactory(u"ws://127.0.0.1:15900")
    factory.protocol = RewarderProtocol

    reactor.connectTCP("127.0.0.1", 15900, factory)
    reactor.run()