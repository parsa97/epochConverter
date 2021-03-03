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

var errCH chan error

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go exporter.Exporter()
	go consumer.StartConsumer(ctx)
	go producer.StartProducer(ctx)
	/* 	go func() {
		for {
			msg, isMore := <-consumer.EpochTimes
			if isMore != true {
				return
			}
			msgi, err := strconv.Atoi(msg)
			if err != nil {
				errCH <- err
			}
			log.Println(time.Unix(0, int64(msgi)*int64(time.Millisecond)).Format(time.RFC3339Nano), isMore)
		}
	}() */
	go func() {
		select {
		case err := <-consumer.ConsumerErrCH:
			log.Println(err)
			c <- os.Interrupt
		}
	}()
	oscall := <-c
	log.Debug("system call: ", oscall)
	cancel()
}
