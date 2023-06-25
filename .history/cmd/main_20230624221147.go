package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	kafka "github.com/odanaraujo/golang/stock-exchange/internal/infra"
	"github.com/odanaraujo/golang/stock-exchange/internal/market/dto"
	"github.com/odanaraujo/golang/stock-exchange/internal/market/entity"
	"github.com/odanaraujo/golang/stock-exchange/internal/transformer"
	"github.com/pkg/errors"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}
	producer := kafka.NewKafkaProducer(configMap)
	consume := kafka.NewConsumer(configMap, []string{"input"})

	go consume.Consume(kafkaMsgChan)

	//recebe do canal do kafka, joga no input, processa, joga no output e depois publica no kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)

	go book.tr

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			if err := json.Unmarshal(msg.Value, &tradeInput); err != nil {
				errors.New("error unmarshal value of trade input in the kafka")
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.Marshal(output)
		fmt.Println(outputJson)
		if err != nil {
			errors.New("error marshal value of trade input in the kafka")
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
