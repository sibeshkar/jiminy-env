package main

import (
	"fmt"
	"os"
	"time"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/sibeshkar/jiminy-env/shared"
	log "github.com/sirupsen/logrus"
)

func main() {

	filename := "../record.rbs"

	reader, err := os.OpenFile(filename, os.O_RDWR, 0644)

	if err != nil {
		log.Error("unable to open file: %s, error: %v", filename, err)
	}

	for {

		msg := &shared.Message{}

		pbutil.ReadDelimited(reader, msg)

		fmt.Println(msg)

		time.Sleep(1 * time.Millisecond)

	}

}
