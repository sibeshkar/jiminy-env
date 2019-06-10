package main

//refer : https://gist.github.com/ImJasonH/6713e855b431a0c067afea4b74cbf504

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/sibeshkar/jiminy-env/utils"
	"github.com/spf13/viper"
)

type EnvState struct {
	//access lock for envID
	EnvId string `json:"env_id"`

	//access lock for status
	sync.Mutex
	EnvStatus string `json:"env_status"`

	EpisodeId int64 `json:"episode_id"`

	Fps float32 `json:"fps"`
}

//Creates new EnvState struct type to store, environment ID, env status,
//episode ID and global FPS.
func NewEnvState(EnvId string, EnvStatus string, EpisodeId int64, Fps float32) (*EnvState, error) {
	var err error //replace with proper function for errors later
	envStatus := EnvState{
		EnvId:     EnvId,
		EnvStatus: EnvStatus,
		EpisodeId: EpisodeId,
		Fps:       Fps,
	}

	envStatus.WriteEnvState()

	return &envStatus, err
}

func (e *EnvState) LoadEnvState() {
	e.Lock()
	defer e.Unlock()
	viper.SetConfigName("env")
	viper.AddConfigPath(utils.TempDir())
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(e)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func (e *EnvState) WriteEnvState() {

	file, _ := json.MarshalIndent(e, "", " ")

	_ = ioutil.WriteFile(utils.TempDir()+"env.json", file, 0644)

}

func (e *EnvState) SetEnvId(EnvId string) bool {
	e.Lock()
	defer e.Unlock()
	e.EnvId = EnvId
	e.WriteEnvState()
	return true
}

func (e *EnvState) SetEnvStatus(EnvStatus string) bool {
	e.Lock()
	defer e.Unlock()
	e.EnvStatus = EnvStatus
	e.WriteEnvState()
	return true
}

func (e *EnvState) SetEpisodeId(EpisodeId int64) bool {
	e.Lock()
	defer e.Unlock()
	e.EpisodeId = EpisodeId
	e.WriteEnvState()
	return true
}

func (e *EnvState) SetFps(Fps float32) bool {
	e.Lock()
	defer e.Unlock()
	e.Fps = Fps
	e.WriteEnvState()
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

func (e *EnvState) GetEpisodeId() int64 {
	e.Lock()
	defer e.Unlock()
	return e.EpisodeId
}

func (e *EnvState) GetFps() float32 {
	e.Lock()
	defer e.Unlock()
	return e.Fps
}
