package main

import (
	kafkaTools "../kafkaTools"
	"github.com/Shopify/sarama"
	"fmt"
)


type MyBusiness struct {
	kafkaTools.KafkaClusterTools
}


func (this *MyBusiness) disposeMsg(*sarama.ConsumerMessage) {
	fmt.Println("Inherit disposeMsg print")
	fmt.Println("Recv msg", string(sarama.ConsumerMessage.Value) )
}


func (this *MyBusiness) endDispose(*sarama.ConsumerMessage) {
	fmt.Println("Inherit endDispose print")
}


func Run() {
	kafkaConfig := kafkaTools.KafkaConfig{
		[]string{"192.168.2.78:9092", "192.168.2.78:9093", "192.168.2.78:9094"},
		[]string{"my-test-1"},
		"my-test-1-group-2",
		[]string{"my-test-2"},
		2,
		-1,
	}

	myBus := new(MyBusiness)
	myBus.SetConfigNormal(kafkaConfig, false)

	myBus.GoRunNUmberInherit(1)
}