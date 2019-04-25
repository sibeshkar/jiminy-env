package main

//refer : https://gist.github.com/ImJasonH/6713e855b431a0c067afea4b74cbf504

import (
	"fmt"
	"math/rand"
	"sync"
)

type EnvState struct {
	//access lock for envID
	EnvId string

	//access lock for status
	sync.Mutex
	EnvStatus string

	EpisodeId int32

	Fps float32
}

//Creates new EnvState struct type to store, environment ID, env status,
//episode ID and global FPS.
func NewEnvState(EnvId string, EnvStatus string, EpisodeId int32, Fps float32) (*EnvState, error) {
	var err error //replace with proper function for errors later
	envStatus := EnvState{
		EnvId:     EnvId,
		EnvStatus: EnvStatus,
		EpisodeId: EpisodeId,
		Fps:       Fps,
	}

	return &envStatus, err
}

func (e EnvState) EnvStatus_Update() {
	var status string
	if rand.Intn(30)%2 == 0 {
		status = "running"
	} else {
		status = "resetting"
	}
	e.Lock()
	e.EnvStatus = status
	e.Unlock()
	fmt.Println("The status now due to EnvStatus_Update is:", e.EnvStatus)
}

func (e EnvState) EpisodeId_Update() {
	e.Lock()
	e.EpisodeId++
	e.Unlock()
	fmt.Println("Updated episode ID is:", e.EpisodeId)
}
