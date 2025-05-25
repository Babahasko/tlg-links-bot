package event_consumer

import (
	"example/tlgbot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher: fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil { //Сам пример не очень, мы тупо пропускаем итерацию.
		//  Если у нас ошибка с сетью, то былобы хорошо завершить и добавить режим ретрая.
			log.Printf("[ERR] Consumer: %s", err.Error())

			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(time.Second * 1)
		}
		if err:= c.handleEvents(gotEvents); err != nil {
			log.Println(err)
			
			continue
		}

	}
}

// TODO:
/*
1. Проблема потери ивентов: ретраи, возвращение в хранилище, фоллбэк, подтверждение для фетчера
2. Обработка всей пачки
3. Парралельная обработка(потребуется waitgroup)
*/

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			log.Printf("can`t handle event: %s", err.Error())
			continue
		}
	}
	return nil
}