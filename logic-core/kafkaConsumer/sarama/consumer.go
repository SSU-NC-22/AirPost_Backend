package sarama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/model"
	"github.com/eunnseo/AirPost/logic-core/setting"
)

var kafkaConsumer *group

type group struct {
	client sarama.ConsumerGroup
	out    chan model.KafkaData
}

func NewKafkaConsumer() *group {
	var err error
	if kafkaConsumer != nil {
		return kafkaConsumer
	}

	outBufSize := setting.Kafkasetting.ChanBufSize
	kafkaConsumer = &group{
		out: make(chan model.KafkaData, outBufSize),
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V0_10_2_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	kafkaConsumer.client, err = sarama.NewConsumerGroup([]string{setting.Kafkasetting.Broker}, setting.Kafkasetting.GroupID, cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	consumer := consumer{
		out:   kafkaConsumer.out,
		ready: make(chan bool),
	}
	go func() {
		for {
			err = kafkaConsumer.client.Consume(ctx, setting.Kafkasetting.Topics, &consumer)
			if err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	return kafkaConsumer
}

func (g *group) GetOutput() <-chan model.KafkaData {
	if g != nil {
		return g.out
	}
	return nil
}

type consumer struct {
	out   chan model.KafkaData
	ready chan bool
}

func (consumer *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		
		fmt.Println("\n\nkafka consumer :", string(message.Value))
		ad := adapter.KafkaData{}
		if err := json.Unmarshal(message.Value, &ad); err != nil {
			log.Print(err)
			continue
		}

		d, err := adapter.KafkaToModel(&ad)
		if err != nil {
			log.Print(err)
			continue
		}
		
		// log.Println("in ConsumeClaim, d(model.KafkaData) = ", d)
		consumer.out <- d // go to "in" in NewLogicCoreUsecase logicCoreUsecase.go
	}

	return nil
}
