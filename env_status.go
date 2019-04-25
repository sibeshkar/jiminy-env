package main

//refer : https://gist.github.com/ImJasonH/6713e855b431a0c067afea4b74cbf504

import (
	"fmt"
	"math/rand"
	"sync"
)

type EnvStatus struct {
	//access lock for envID
	envLock sync.Mutex
	EnvId   string

	//access lock for status
	statusLock sync.Mutex
	EnvStatus  string

	episodeLock sync.Mutex
	EpisodeId   int32

	fpsLock sync.Mutex
	Fps     float32
}

// var envstatus EnvStatus

// func main() {
// 	envstatus = EnvStatus{
// 		EnvId:     "wob.mini.TicTacToe",
// 		EnvStatus: "running",
// 		EpisodeId: 45,
// 		Fps:       60,
// 	}
// 	//i := 0
// 	for {
// 		go envstatus.EnvStatus_Update()
// 		go envstatus.EnvStatus_Update2()
// 		envstatus.EpisodeId_Update()
// 		time.Sleep((1 / 4) * time.Second)
// 		//fmt.Println(i)
// 		//i++
// 	}

// }

func (e EnvStatus) EnvStatus_Update() {
	var status string
	if rand.Intn(30)%2 == 0 {
		status = "running"
	} else {
		status = "resetting"
	}
	e.statusLock.Lock()
	e.EnvStatus = status
	e.statusLock.Unlock()
	fmt.Println("The status now due to EnvStatus_Update is:", e.EnvStatus)
}

func (e EnvStatus) EnvStatus_Update2() {
	var status string
	if rand.Intn(30)%3 == 0 {
		status = "running2"
	} else {
		status = "resetting2"
	}
	e.statusLock.Lock()
	e.EnvStatus = status
	e.statusLock.Unlock()
	fmt.Println("The status now due to EnvStatus_Update2 is:", e.EnvStatus)
}

func (e EnvStatus) EpisodeId_Update() {
	e.episodeLock.Lock()
	e.EpisodeId++
	e.episodeLock.Unlock()
	fmt.Println("Updated episode ID is:", e.EpisodeId)
}
