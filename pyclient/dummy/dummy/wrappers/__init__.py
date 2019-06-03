from jiminy.wrappers.recording import Recording
from jiminy.wrappers.render import Render
from jiminy.wrappers.throttle import Throttle

def wrap(env):
    return Timer(Render(Throttle(env)))

def WrappedVNCEnv():
    return wrap(envs.DummyVNCEnv())