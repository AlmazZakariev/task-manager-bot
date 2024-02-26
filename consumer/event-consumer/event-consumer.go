package event_consumer

import (
	"log"
	"tasks-manager-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	executor  events.Executor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, executor events.Executor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		executor:  executor,
		batchSize: batchSize,
	}
}
func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			//здесь могут быть проблемы с сетью
			//в fetcher можно встроить механизм retry
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		gotTasksToSend, err := c.executor.Execute()
		if err != nil {
			log.Printf("[ERR] Executor: %s", err.Error())
		}

		if len(gotEvents) == 0 && len(gotTasksToSend) == 0 {
			time.Sleep(1 * time.Second)
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

		}

		if err := c.handleEvents(gotTasksToSend); err != nil {
			log.Print(err)
		}
	}
}
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			//TODO: нужно сохранять событие, если произошла ошибка
			log.Printf("can't handle event %s", err.Error())

			continue
		}
	}
	return nil
}
