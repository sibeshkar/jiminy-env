package main

import (
	"math/rand"
	"time"
)

func RewardController(status, trigger chan string, rewardChan chan Body) {
	for {
		if rand.Intn(10)%2 == 0 {
			status <- "running"
		} else {
			status <- "resetting"
		}
	}

}

//Random function to generate get reward from environment
func ConstructRewardMessage(reward Body) Message {
	//env.GetReward()string

	method := "v0.env.reward"

	headers := Headers{
		Sent_at: time.Now().Unix(),
	}

	body := Body{
		EnvId:  "wob.mini.TicTacToe",
		Reward: reward.Reward,
		Done:   reward.Done,
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
func GetReward() Body {
	//env.GetReward()string
	reward, done, _ := env.GetReward()
	reply := Body{
		Reward: reward,
		Done:   done,
	}

	return reply
}
