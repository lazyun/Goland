package main

import (
	kafkaTools "../kafkaTools"
	"fmt"
	"time"
	"sync"
	"github.com/Shopify/sarama"
)

func main() {

	tools := kafkaTools.KafkaClusterTools{}

	tools.SetConfigNormal([]string{"192.168.2.78:9092", "192.168.2.78:9093", "192.168.2.78:9094"}, []string{"my-test-1"}, "my-test-1-group-2", false)

	tools.GoRunNUmber(2, businessInit, nil)
	//wg := sync.WaitGroup{}
	//for i := 0; i < 3; i++ {
	//	wg.Add(1)
	//	go wgTest(i, &wg)
	//}
	//
	//wg.Wait()


}


func wgTest(t int, wg *sync.WaitGroup) {
	time.Sleep( time.Second * time.Duration(t) )
	fmt.Println( t, "run out!" )
	wg.Done()
}


func businessInit() (func (*sarama.ConsumerMessage)) {
	testObj := Business{}

	return testObj.Dispose
}


type Business struct {
	partition 	int32
	offset 		int64
	byteMsg		[]byte
	stringMsg	string
	timestamp 	time.Time
}

func (this *Business) Dispose(msg *sarama.ConsumerMessage) {
	this.partition = msg.Partition
	this.offset = msg.Offset
	this.byteMsg = msg.Value
	this.stringMsg = string(this.byteMsg)

	this.timestamp = msg.Timestamp

	this.save()
}


func (this *Business) save() {
	fmt.Println("recv msg partition is", this.partition,
		"offset is", this.offset,
			"byte message is", this.byteMsg,
				"string message is", this.stringMsg,
					"timestamp is", this.timestamp )
}