package epochConsumer

import (
	"context"
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestStartConsumer(t *testing.T) {
	var msgCH = make(chan string, 10000)
	var errCH = make(chan error, 0)
	_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", "9092"), 5*time.Second)
	if err != nil {
		log.Error("Error when connecting to Kafka is there any kafka cluster? ", err)
		t.Fatal()
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = StartConsumer(ctx, errCH, msgCH)
	if err != nil {
		t.Fatal(err)
	}
}
