package kafkaTools


import (
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
	"syscall"
)

type KafkaConfig struct {
	Hosts     	[]string   	`json:"hosts"`
	Topic     	[]string   	`json:"topic"`
	Consumer  	string   	`json:"consumer"`
	TopicSend 	[]string 	`json:"topic_send"`

	StartNumb	int			`json:"start_number"`
	RunCount	int64		`json:"run_count"`
}

type KafkaClusterTools struct {
	Broker 			[]string
	Topic			[]string
	Group 			string

	logHandle		func(interface{})

	kafkaConfig		*cluster.Config
	kafkaConsumer	*cluster.Consumer

	DisposeFunc		func()
	RunCount		int64
}


func (this *KafkaClusterTools) SetConfigNormal(cfg KafkaConfig, IsUsePartitions bool) {
	config := cluster.NewConfig()

	if IsUsePartitions {
		config.Group.Mode = cluster.ConsumerModePartitions
	}

	//config.Consumer.Offsets.CommitInterval = 1
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // OffsetOldest OffsetNewest
	//config.Version = kafkaVer

	this.kafkaConfig = config
	this.Broker = cfg.Hosts
	this.Topic = cfg.Topic
	this.Group = cfg.Consumer
	this.RunCount = cfg.RunCount

	this.Connect()
}


func (this *KafkaClusterTools) SetConfigSasl(hosts, topic []string, group, saslUser, saslPasswd, certFilePath string, IsUsePartitions bool) {
	config := cluster.NewConfig()

	config.Net.SASL.Enable = true
	config.Net.SASL.User = saslUser
	config.Net.SASL.Password = saslPasswd
	config.Net.SASL.Handshake = true
	certFile := certFilePath + "/ca-cert"

	certBytes, _ := ioutil.ReadFile( certFile )
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("Kafka consumer failed to parse root certificate")
	}

	config.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	config.Net.TLS.Enable = true

	if IsUsePartitions {
		config.Group.Mode = cluster.ConsumerModePartitions
	}


	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // OffsetOldest OffsetNewest
	//config.Version = kafkaVer

	this.Broker = hosts
	this.Topic = topic
	this.Group = group

	this.Connect()
}


func (this *KafkaClusterTools) Connect() {
	var err error
	this.kafkaConsumer, err = cluster.NewConsumer(this.Broker, this.Group, this.Topic, this.kafkaConfig)

	if nil != err {
		panic( err )
	}
}



func (this *KafkaClusterTools) SetLogHandle(funcName func(interface{})) {
	this.logHandle = funcName
}


func (this *KafkaClusterTools) SetLog(msg string) {
	if nil == this.logHandle {
		fmt.Println(msg)
		return
	}

	this.logHandle(msg)
}


func (this *KafkaClusterTools) GoRunAll(start func() (func (*sarama.ConsumerMessage) ), end func ()) {
	// consume errors
	go func() {
		//fmt.Println("biubiu~")
		for err := range this.kafkaConsumer.Errors() {
			//fmt.Printf("%s:Error: %s\n", this.Group, err.Error())
			content := fmt.Sprintf("%s:Error: %s\n", this.Group, err.Error())
			this.SetLog(content)
		}
	}()

	// consume notifications
	go func() {
		//fmt.Println("lalala~")
		for ntf := range this.kafkaConsumer.Notifications() {
			//fmt.Printf("%s:Rebalanced: %+v \n", this.Group, ntf)
			content := fmt.Sprintf("%s:Rebalanced: %+v \n", this.Group, ntf)
			this.SetLog(content)
		}
	}()

	wg := sync.WaitGroup{}

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, os.Interrupt)

	startCount := make(chan int, 200)
	exitAllCorou := make(chan int, 200)

	for {
		select {
		case part, ok := <-this.kafkaConsumer.Partitions():
			//fmt.Println("recv partition is", part, ok)
			if !ok {
				return
			}

			go func(pc cluster.PartitionConsumer) {
				wg.Add(1)
				startCount<- 1

				defer func() {
					if nil != end {
						end()
					}

					pc.Close()
					wg.Done()
					<-startCount
				}()

				dispose := start()

				for {
					select {
					case msg := <- pc.Messages():
						dispose(msg)
						this.kafkaConsumer.MarkOffset(msg, "")
						this.kafkaConsumer.CommitOffsets()
					case ret, ok := <-exitAllCorou:
						fmt.Println("signal recv is", ret, ok, "partition is", part.Partition())
						return
					}
				}

				//for msg := range pc.Messages() {
				//
				//	select {
				//	case <-signals:
				//		return
				//	default:
				//		dispose(msg)
				//		this.kafkaConsumer.MarkOffset(msg, "")
				//		this.kafkaConsumer.CommitOffsets()
				//	}
				//
				//}
			}(part)
		case ret, ok := <-signals:
			fmt.Println("signal recv is", ret, ok)

			for ret := range startCount {
				exitAllCorou <- ret
			}

			wg.Wait()
			return
		}
	}
}


func (this *KafkaClusterTools) GoRunNUmber(total int, start func() (func (*sarama.ConsumerMessage) ), end func () ) {

	// consume errors
	go func() {
		//fmt.Println("biubiu~")
		for err := range this.kafkaConsumer.Errors() {
			//fmt.Printf("%s:Error: %s\n", this.Group, err.Error())
			now := time.Now().Format("2006-01-02 15:04:05")
			content := fmt.Sprintf("%s\t%s:Error: %s\n", now, this.Group, err.Error())
			this.SetLog(content)
		}
	}()

	// consume notifications
	go func() {
		//fmt.Println("lalala~")
		for ntf := range this.kafkaConsumer.Notifications() {
			//fmt.Printf("%s:Rebalanced: %+v \n", this.Group, ntf)
			now := time.Now().Format("2006-01-02 15:04:05")
			content := fmt.Sprintf("%s\t%s:Rebalanced: %+v \n", now, this.Group, ntf)
			this.SetLog(content)
		}
	}()

	wg := sync.WaitGroup{}
	//wgIsStart := sync.WaitGroup{}

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, os.Interrupt)

	startCount := make(chan int, 200)
	exitAllCorou := make(chan int, 200)

	for i := 0; i < total; i++ {
		startCount<- 1
	}

	//id := 0

	for {
		select {
		case <-startCount:

			go func() {
				//fmt.Println("Start go coroutine from Messages()")
				wg.Add(1)

				defer func() {
					if nil != end {
						end()
					}

					wg.Done()
					startCount<- 1
				}()

				dispose := start()

				//temp := id

				var recvCount int64 = 0
				for {
					select {
					case msg := <- this.kafkaConsumer.Messages():
						//fmt.Print( temp, os.Getpid() )
						//fmt.Println("Recv msg", )
						dispose(msg)
						this.kafkaConsumer.MarkOffset(msg, "")
						this.kafkaConsumer.CommitOffsets()

						recvCount++
						if recvCount == this.RunCount {
							//syscall.
							signals<- syscall.SIGUSR1
							<-exitAllCorou
							return
						}
					case ret, ok := <-exitAllCorou:
						fmt.Println("Go coroutine signal recv is", ret, ok)
						return
					}
				}

			}()

			//time.Sleep(time.Second * 10)
			//id++
		case ret, ok := <-signals:
			fmt.Println("signal recv is", ret, ok)

			for i := 0; i < total; i++ {
				exitAllCorou<- 1
			}


			wg.Wait()
			return
		}
	}
}


func (this *KafkaClusterTools) GoRunNUmberInherit(total int) {

	// consume errors
	go func() {
		//fmt.Println("biubiu~")
		for err := range this.kafkaConsumer.Errors() {
			//fmt.Printf("%s:Error: %s\n", this.Group, err.Error())
			now := time.Now().Format("2006-01-02 15:04:05")
			content := fmt.Sprintf("%s\t%s:Error: %s\n", now, this.Group, err.Error())
			this.SetLog(content)
		}
	}()

	// consume notifications
	go func() {
		//fmt.Println("lalala~")
		for ntf := range this.kafkaConsumer.Notifications() {
			//fmt.Printf("%s:Rebalanced: %+v \n", this.Group, ntf)
			now := time.Now().Format("2006-01-02 15:04:05")
			content := fmt.Sprintf("%s\t%s:Rebalanced: %+v \n", now, this.Group, ntf)
			this.SetLog(content)
		}
	}()

	wg := sync.WaitGroup{}
	//wgIsStart := sync.WaitGroup{}

	signals := make(chan os.Signal, 10)
	signal.Notify(signals, os.Interrupt)

	startCount := make(chan int, 200)
	exitAllCorou := make(chan int, 200)

	for i := 0; i < total; i++ {
		startCount<- 1
	}

	//id := 0

	for {
		select {
		case <-startCount:

			go func() {
				//fmt.Println("Start go coroutine from Messages()")
				wg.Add(1)

				defer func() {
					this.endDispose()
					wg.Done()
					startCount<- 1
				}()

				//temp := id

				var recvCount int64 = 0
				for {
					select {
					case msg := <- this.kafkaConsumer.Messages():
						//fmt.Print( temp, os.Getpid() )
						//fmt.Println("Recv msg", )
						this.disposeMsg(msg)
						this.kafkaConsumer.MarkOffset(msg, "")
						this.kafkaConsumer.CommitOffsets()

						recvCount++
						if recvCount == this.RunCount {
							//syscall.
							signals<- syscall.SIGUSR1
							<-exitAllCorou
							return
						}
					case ret, ok := <-exitAllCorou:
						fmt.Println("Go coroutine signal recv is", ret, ok)
						return
					}
				}

			}()

			//time.Sleep(time.Second * 10)
			//id++
		case ret, ok := <-signals:
			fmt.Println("signal recv is", ret, ok)

			for i := 0; i < total; i++ {
				exitAllCorou<- 1
			}


			wg.Wait()
			return
		}
	}
}


type InheritFunc interface {
	disposeMsg(*sarama.ConsumerMessage)
	endDispose()
}


func (this *KafkaClusterTools) disposeMsg(*sarama.ConsumerMessage) {

}


func (this *KafkaClusterTools) endDispose() {

}