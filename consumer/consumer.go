package epochConsumer

import (
	"context"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartConsumer(ctx context.Context, errCH chan error, msgCH chan string) error {
	log.Info("Server Start Consuming")
	topic := defaultTopic
	if value, ok := os.LookupEnv("CONSUMER_TOPIC"); ok {
		topic = value
	}
	consumer, err := newConsumer()
	if err != nil {
		errCH <- err
		return err
	}

	go func() {
		for err := range consumer.Errors() {
			log.Error(err)
		}
	}()

	member := ConsumerGroupMember{localMsgCh: msgCH}

	go func() {
		for {
			err := consumer.Consume(ctx, strings.Split(topic, ","), &member)
			if err != nil {
				log.Error(err)
			}
			if ctx.Err() != nil {
				errCH <- err
				return
			}
		}
	}()

	<-ctx.Done()
	log.Info("Consumer Stopped")
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancel()
	}()
	return nil
}
