package rfcProducer

import (
	"context"
	consumer "epochConvertor/consumer"
	exporter "epochConvertor/metrics"
	"errors"
	"os"
	"strconv"
	"time"

	cmap "github.com/orcaman/concurrent-map"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func StartProducer(ctx context.Context) {
	topic := "output"
	if value, ok := os.LookupEnv("PRODUCER_TOPIC"); ok {
		topic = value
	}

	producer, err := newProducer()
	if err != nil {
		log.Error("Could not create producer: ", err)
		os.Exit(1)
	}
	if err := sendRFCMessage(ctx, producer, topic); err != nil {
		log.Error("failed to produce: ", err)
	}
}

func sendRFCMessage(ctx context.Context, producer sarama.SyncProducer, topic string) error {
	log.Info("Server Start Producing")
	var partitionProduced = cmap.New()
	go func() error {
		for {
			time.Sleep(time.Millisecond)
			epochMsg, isMore := <-consumer.EpochTimes
			if isMore != true {
				err := errors.New("wait for new messages to consume")
				log.Println(err)
				return err
			}
			msgi, err := strconv.Atoi(epochMsg)
			if err != nil {
				log.Println(err)
			}
			msg := prepareMessage(topic, time.Unix(0, int64(msgi)*int64(time.Millisecond)).Format(time.RFC3339Nano))
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				return err
			}
			counter, _ := partitionProduced.Get(string(partition))
			if counter == nil {
				partitionProduced.Set(string(partition), 1)
			} else {
				counted := counter.(int)
				counted++
				partitionProduced.Set(string(partition), counted)
				log.Println("Counter: ", " Message: ", epochMsg, " topic: ", topic, " partition: ", partition, " offset: ", offset)
				M := exporter.ProducedMessageCounter.WithLabelValues(strconv.Itoa(int(partition)), topic)
				M.Inc()
			}
		}
	}()
	<-ctx.Done()
	log.Info("Server Stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	log.Info("Server Exited Properly")

	return nil
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}
