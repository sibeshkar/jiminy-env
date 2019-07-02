package main

import (
	vncproxy "github.com/amitbet/vncproxy/proxy"
)

func create_vnc_proxy(WsListeningURL string, RecordingDir string, TCPListeningURL string, ProxyVncPassword string, TargetHostname string, TargetPassword string, TargetPort string, ID string) *vncproxy.VncProxy {
	//create default session if required
	proxy := &vncproxy.VncProxy{
		WsListeningURL:  WsListeningURL, // empty = not listening on ws
		RecordingDir:    RecordingDir,   // empty = no recording
		TCPListeningURL: TCPListeningURL,
		//RecordingDir:          "C:\\vncRec", // empty = no recording
		ProxyVncPassword: ProxyVncPassword, //empty = no auth
		SingleSession: &vncproxy.VncSession{
			TargetHostname: TargetHostname,
			TargetPort:     TargetPort,
			TargetPassword: TargetPassword,
			ID:             ID,
			Status:         vncproxy.SessionStatusInit,
			Type:           vncproxy.SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	return proxy
}
