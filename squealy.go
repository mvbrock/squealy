package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/serf/serf"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"strconv"
	"strings"
	"sync"
)

func receiveFromKafka(client *kgo.Client, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for {
		fetches := client.PollFetches(context.Background())
		if fetches.IsClientClosed() {
			return
		}
		fetches.EachError(func(topic string, part int32, err error) {
			log.Fatalf("kafka error on topic: %s, partition: %d, error: %v\n", topic, part, err)
		})
		fetches.EachRecord(func(record *kgo.Record) {
			fmt.Printf("kafka message at topic: %v, partition: %v, offset:%v  %s = %s\n",
				record.Topic, record.Partition, record.Offset, record.Key, record.Value)
		})
	}
}

func receiveFromSerf(serfChan chan serf.Event, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for {
		event := <-serfChan
		fmt.Printf("serf event: %s\n", event.String())
	}
}

func main() {
	// Kafka options
	kafkaHost := flag.String("kafkaHost", "localhost:9092", "The hostname:port string for the Kafka queue")
	kafkaTopic := flag.String("kafkaTopic", "default", "The Kafa topic to read messages from")
	kafkaConsumerGroup := flag.String("kakfaConsumerGroup", "default", "The Kafka consumer group")

	// Serf options
	serfBind := flag.String("serfBind", "", "The hostname:port for binding the Serf agent")
	serfConnect := flag.String("serfConnect", "", "The hostname:port of an instance to join")
	serfNodeName := flag.String("serfNodeName", "0", "The unique node name of this instance")

	// Parse the command line options
	flag.Parse()

	// Print the command line options
	fmt.Println("kafkaHost: ", *kafkaHost)
	fmt.Println("kafkaTopic: ", *kafkaTopic)
	fmt.Println("kafkaConsumerGroup: ", *kafkaConsumerGroup)
	fmt.Println("serfBind: ", *serfBind)
	fmt.Println("serfConnect: ", *serfConnect)
	fmt.Println("serfNodeName: ", *serfNodeName)

	// Establish the Serf configuration
	serfConfig := serf.DefaultConfig()
	serfChan := make(chan serf.Event)
	serfConfig.EventCh = serfChan
	serfConfig.NodeName = *serfNodeName
	if *serfBind != "" {
		bindParts := strings.Split(*serfBind, ":")
		serfConfig.MemberlistConfig.BindAddr = bindParts[0]
		bindPort, err := strconv.Atoi(bindParts[1])
		if err != nil {
			log.Fatalln(err)
		}
		serfConfig.MemberlistConfig.BindPort = bindPort
	}
	serfAgent, err := serf.Create(serfConfig)
	if err != nil {
		log.Fatalln(err)
	}
	if *serfConnect != "" {
		serfAgent.Join([]string{*serfConnect}, false)
	}

	// Print out the member list
	members := serfAgent.Members()
	for _, member := range members {
		fmt.Println(member)
	}

	// Establish the Kafka reader
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers([]string{*kafkaHost}...),
		kgo.ConsumerGroup(*kafkaConsumerGroup),
		kgo.ConsumeTopics(*kafkaTopic))
	if err != nil {
		log.Fatalln(err)
	}

	// Create the wait group for the Kafka and Serf reader threads
	var waitGroup sync.WaitGroup

	// Start reading messages off the Serf channel
	waitGroup.Add(1)
	go receiveFromSerf(serfChan, &waitGroup)

	// Start reading messages off the Kafka queue
	waitGroup.Add(1)
	go receiveFromKafka(kafkaClient, &waitGroup)

	// Wait for the threads to complete
	waitGroup.Wait()
}
