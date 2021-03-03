package epochConsumer

import (
	"context"
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestStartConsumer(t *testing.T) {
	_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", "9092"), 5*time.Second)
	if err != nil {
		log.Error("Error when connecting to Kafka is there any kafka cluster?", err)
		t.Fatal()
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = StartConsumer(ctx)
	if err != nil {
		t.Fatal(err)
	}
}