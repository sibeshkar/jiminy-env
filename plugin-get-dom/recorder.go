package main

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/sibeshkar/jiminy-env/shared"
	vnc_rec "github.com/sibeshkar/vncproxy/vnc_rec"
	log "github.com/sirupsen/logrus"
)

func NewVncProxy(WsListeningURL string, RecordingDir string, TCPListeningURL string, ProxyVncPassword string, TargetHostname string, TargetPassword string, TargetPort string, ID string) *vnc_rec.VncProxy {
	//create default session if required
	proxy := &vnc_rec.VncProxy{
		WsListeningURL:  WsListeningURL, // empty = not listening on ws
		RecordingDir:    RecordingDir,   // empty = no recording
		TCPListeningURL: TCPListeningURL,
		//RecordingDir:          "C:\\vncRec", // empty = no recording
		ProxyVncPassword: ProxyVncPassword, //empty = no auth
		SingleSession: &vnc_rec.VncSession{
			TargetHostname: TargetHostname,
			TargetPort:     TargetPort,
			TargetPassword: TargetPassword,
			ID:             ID,
			Status:         vnc_rec.SessionStatusInit,
			Type:           vnc_rec.SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	return proxy
}

type Recorder struct {
	filepath  string
	writer    *os.File
	batch     *shared.Message
	batchChan chan *shared.Message
	sync.Mutex
}

func NewRecorder(dirname string) (*Recorder, error) {
	recFile := "record.rbs"
	recPath := path.Join(dirname, recFile)
	if _, err := os.Stat(recPath); err == nil {
		os.Remove(recPath)
	}

	writer, err := os.OpenFile(recPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Error("unable to open file: %s, error: %v", recPath, err)
		return nil, err
	}

	recorder := &Recorder{
		filepath: recPath,
		writer:   writer,
	}

	recorder.NewBatch()

	return recorder, err

}

func (r *Recorder) NewBatch() {

	r.Lock()
	defer r.Unlock()

	method := "v0.env.record"
	timestamp := uint32(time.Now().UnixNano() / int64(time.Millisecond))

	body := &shared.Body{}

	r.batch = &shared.Message{
		Method:    method,
		Body:      body,
		Timestamp: timestamp,
	}

}

// func (r Recorder) StartListening() {
// 	for {
// 		data := <-r.batchChan
// 		r.writeToDisk(data)
// 	}

// }

func (r *Recorder) writeToDisk() {
	pbutil.WriteDelimited(r.writer, r.batch)

}

// func (r Recorder) GetBatch() *shared.Message {
// 	r.Lock()
// 	defer r.Unlock()
// 	return r.batch

// }

// func (r *Recorder) pushToChannel() {
// 	r.batchChan <- r.batch
// }

func (r *Recorder) AddRewardtoBatch(reward float32, done bool, info string) {

	r.Lock()
	defer r.Unlock()

	r.batch.Body.Reward = reward
	r.batch.Body.Done = done
	r.batch.Body.Info = info

	// r.batch.Body.Done = done
	// r.batch.Body.Info = info
}

func (r *Recorder) AddObstoBatch(obstype string, obs string) {

	r.Lock()
	defer r.Unlock()

	r.batch.Body.Obs = obs
	r.batch.Body.ObsType = obstype

}
