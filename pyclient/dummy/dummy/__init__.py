from jiminy.gym.envs.registration import register

register(
    id='VNC.Core-v0',
    entry_point='dummy.envs:DummyVNCEnv',
)
