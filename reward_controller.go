package main

import (
	"math/rand"
	"time"
)

func RewardController(status, trigger chan string, reward chan float32) {

	for {
		if rand.Intn(10)%2 == 0 {
			status <- "running"
		} else {
			status <- "resetting"
		}
		time.Sleep(10 * time.Millisecond)
		reward <- GetReward()
	}

}

//Random function to generate get reward from environment
func ConstructRewardMessage(reward float32) Message {
	//env.GetReward()string

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

	return m
}

//---------------------------
// Environment Plugin Interface
//---------------------------
//GetReward: Random function to generate get reward from environment plugin
func GetReward() float32 {
	//env.GetReward()string
	reward, _ := env.GetReward()
	return reward
}
