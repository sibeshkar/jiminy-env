package main

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/sibeshkar/jiminy-env/shared"
	"github.com/sibeshkar/vncproxy/logger"
	vnc_rec "github.com/sibeshkar/vncproxy/vnc_rec"
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
		logger.Errorf("unable to open file: %s, error: %v", recPath, err)
		return nil, err
	}

	recorder := &Recorder{
		filepath: recPath,
		writer:   writer,
	}

	method := "v0.env.record"
	timestamp := float32(time.Now().UnixNano() / 1000000)

	recorder.batch = &shared.Message{
		Method:    method,
		Timestamp: timestamp,
	}
	recorder.NewBatch()
	return recorder, err

}

func (r Recorder) NewBatch() {

	// r.msgLock.Lock()
	// defer r.msgLock.Unlock()
	r.writeToDisk()
	method := "v0.env.record"
	timestamp := float32(time.Now().UnixNano() / 1000000)

	r.batch = &shared.Message{
		Method:    method,
		Timestamp: timestamp,
	}

}

// func (r Recorder) StartListening() {
// 	for {
// 		data := <-r.batchChan
// 		r.writeToDisk(data)
// 	}

// }

func (r Recorder) writeToDisk() {
	pbutil.WriteDelimited(r.writer, r.batch)

}

func (r Recorder) GetBatch() *shared.Message {
	r.Lock()
	defer r.Unlock()
	return r.batch

}

func (r Recorder) pushToChannel() {
	r.batchChan <- r.batch
}

func (r Recorder) AddRewardtoBatch(reward float32, done bool, info string) {

	r.Lock()
	defer r.Unlock()

	method := r.batch.GetMethod()

	timestamp := r.batch.GetTimestamp()

	body := &shared.Body{
		Reward: reward,
		Done:   done,
		Info:   info,
	}

	r.batch = &shared.Message{
		Method:    method,
		Body:      body,
		Timestamp: timestamp,
	}

	// r.batch.Body.Done = done
	// r.batch.Body.Info = info
}

func (r Recorder) AddObstoBatch(obstype string, obs string) {

	r.Lock()
	defer r.Unlock()

	method := r.batch.GetMethod()

	timestamp := r.batch.GetTimestamp()

	bodyRef := r.batch.GetBody()

	body := &shared.Body{
		Reward:  bodyRef.GetReward(),
		Done:    bodyRef.GetDone(),
		Info:    bodyRef.GetInfo(),
		Obs:     obs,
		ObsType: obstype,
	}

	r.batch = &shared.Message{
		Method:    method,
		Body:      body,
		Timestamp: timestamp,
	}

}
