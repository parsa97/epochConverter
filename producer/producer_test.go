package rfcProducer

import (
	"context"
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

const topic = "test"

func TestSendEpochMessage(t *testing.T) {
	var errCH = make(chan error, 0)
	var msgCH = make(chan string, 10000)
	_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", "9092"), 5*time.Second)
	if err != nil {
		log.Error("Error when connecting to Kafka is there any kafka cluster? ", err)
		t.Fatal()
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = StartProducer(ctx, errCH, msgCH)
	if err != nil {
		t.Fatal()
	}
}
