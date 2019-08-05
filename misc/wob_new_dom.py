import jiminy
import jiminy.gym as gym
import time

if __name__ == "__main__":
    env = gym.make("VNC.Core-v0")
    env = jiminy.wrappers.experimental.SoftmaxClickMouse(env)

    env.configure(env='sibeshkar/wob-v1', task='ClickButton', remotes='vnc://localhost:5901+15901')
    obs = env.reset()

    while True:
        a = env.action_space.sample()
        obs, reward, is_done, info = env.step([a])
        if obs[0]['dom'] is None: #substitute with obs[0]['vision'] for image equivalent
            print("Env is still resetting...")
            continue
        break

    for idx in range(5000):
        time.sleep(0.05)
        #a = env.action_space.sample()
        obs, reward, is_done, info = env.step([a])
        if obs[0] is None:
            print("Env is resetting...")
            continue
        print("Sampled action: ", a)
        print("Response are of index:", idx)
        print("Observation", obs[0]['dom'])
        print("Reward", reward)
        print("Is done", is_done)
        print("Info", info)
        if is_done[0]:
            time.sleep(0.5)
        env.render()


        #im = Image.fromarray(obs[0])
        #im.save("test_frames/frame-%03d.png" % idx)

    env.close()
    pass