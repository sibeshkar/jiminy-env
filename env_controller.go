package main

import "time"

func (c *AgentConn) EnvController() error {
	for {
		c.SendReward(GetReward())
		time.Sleep(time.Duration(1/c.envState.Fps) * time.Second)
	}

}
