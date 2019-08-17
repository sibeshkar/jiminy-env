package main

import (
	"fmt"
	"os"
	"time"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	pb_demo "github.com/sibeshkar/demoparser/proto"
	log "github.com/sirupsen/logrus"
)

func main() {

	filename := "/Users/sibeshkar/go/src/github.com/sibeshkar/demoparser/demo/recording_1565742906/record.rbs"

	reader, err := os.OpenFile(filename, os.O_RDWR, 0644)

	if err != nil {
		log.Error("unable to open file: %s, error: %v", filename, err)
	}

	for {

		msg := &pb_demo.Message{}

		pbutil.ReadDelimited(reader, msg)

		fmt.Println(msg)

		time.Sleep(1 * time.Millisecond)

	}

}
