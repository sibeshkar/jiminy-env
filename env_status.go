package main

//refer : https://gist.github.com/ImJasonH/6713e855b431a0c067afea4b74cbf504

import (
	"sync"
)

type EnvState struct {
	//access lock for envID
	EnvId string `json:"env_id"`

	//access lock for status
	sync.Mutex
	EnvStatus string `json:"env_status"`

	EpisodeId int `json:"episode_id"`

	Fps float32 `json:"fps"`
}

//Creates new EnvState struct type to store, environment ID, env status,
//episode ID and global FPS.
func NewEnvState(EnvId string, EnvStatus string, EpisodeId int, Fps float32) (*EnvState, error) {
	var err error //replace with proper function for errors later
	envStatus := EnvState{
		EnvId:     EnvId,
		EnvStatus: EnvStatus,
		EpisodeId: EpisodeId,
		Fps:       Fps,
	}

	return &envStatus, err
}

func (e *EnvState) SetEnvId(EnvId string) bool {
	e.Lock()
	defer e.Unlock()
	e.EnvId = EnvId
	return true
}

func (e *EnvState) SetEnvStatus(EnvStatus string) bool {
	e.Lock()
	defer e.Unlock()
	e.EnvStatus = EnvStatus
	return true
}

func (e *EnvState) SetEpisodeId(EpisodeId int) bool {
	e.Lock()
	defer e.Unlock()
	e.EpisodeId = EpisodeId
	return true
}

func (e *EnvState) SetFps(Fps float32) bool {
	e.Lock()
	defer e.Unlock()
	e.Fps = Fps
	return true
}

func (e *EnvState) GetEnvId() string {
	e.Lock()
	defer e.Unlock()
	return e.EnvId
}

func (e *EnvState) GetEnvStatus() string {
	e.Lock()
	defer e.Unlock()
	return e.EnvStatus
}

func (e *EnvState) GetEpisodeId() int {
	e.Lock()
	defer e.Unlock()
	return e.EpisodeId
}

func (e *EnvState) GetFps() float32 {
	e.Lock()
	defer e.Unlock()
	return e.Fps
}
