package epochConsumer

import (
	"context"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartConsumer(ctx context.Context) error {
	log.Info("Server Start Consuming")
	if value, ok := os.LookupEnv("CONSUMER_TOPIC"); ok {
		consumerTopic = value
	}
	consumer, err := newConsumer()
	if err != nil {
		ConsumerErrCH <- err
		return err
	}

	go func() {
		for err := range consumer.Errors() {
			log.Error(err)
		}
	}()

	member := ConsumerGroupMember{MsgCh: EpochTimes}

	go func() {
		for {
			err := consumer.Consume(ctx, strings.Split(consumerTopic, ","), &member)
			if err != nil {
				log.Error(err)
			}
			if ctx.Err() != nil {
				ConsumerErrCH <- err
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("Consumer Stopped")
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancel()
	}()
	return nil
}
