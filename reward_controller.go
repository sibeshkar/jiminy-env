package main

import (
	"log"
	"time"
)

func (c *AgentConn) RewardController() error {
	for {
		c.SendReward(GetReward())
		time.Sleep(time.Duration(1/c.envState.Fps) * time.Second)
	}

}

func (c *AgentConn) SendReward(m *Message) error {
	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

//Random function to generate get reward from environment
func GetReward() *Message {
	//env.GetReward()string
	reward, _ := env.GetReward()

	method := "v0.env.reward"

	headers := Headers{
		Sent_at: time.Now().Unix(),
	}

	body := Body{
		EnvId:  "wob.mini.TicTacToe",
		Reward: reward,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	return &m
}
