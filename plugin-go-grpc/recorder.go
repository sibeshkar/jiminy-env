package main

import (
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
