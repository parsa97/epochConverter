package main

import (
	"context"
	consumer "epochConvertor/consumer"
	exporter "epochConvertor/metrics"
	producer "epochConvertor/producer"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(logLevel())
}

func main() {
	c := make(chan os.Signal, 1)
	var errCH = make(chan error, 0)
	var msgCH = make(chan string, 10000)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go exporter.Exporter(errCH)
	go consumer.StartConsumer(ctx, errCH, msgCH)
	go producer.StartProducer(ctx, errCH, msgCH)
	go func() {
		select {
		case err := <-errCH:
			log.Error(err)
			c <- os.Interrupt
		}
	}()
	oscall := <-c
	log.Debug("system call: ", oscall)
	cancel()
}
