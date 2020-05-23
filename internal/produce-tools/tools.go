package produce_tools

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type statistic struct {
	success int
	fail    int
	lock    sync.RWMutex
}

type Producer struct {
	cmd       *kafka.Producer
	ConfigMap map[string]interface{}
	codec     string
	lock      sync.Mutex
	Stat      statistic
}

func NewKafkaProducer(config map[string]interface{}) (*Producer, error) {
	kafkaCfgMap := kafka.ConfigMap{}
	for k, v := range config {
		_ = kafkaCfgMap.SetKey(k, v)
	}
	p, err := kafka.NewProducer(&kafkaCfgMap)
	return &Producer{
		cmd:       p,
		ConfigMap: config,
		lock:      sync.Mutex{},
		codec:     config["compression.type"].(string),
	}, err
}

func (p *Producer) GetConfig() map[string]interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.ConfigMap
}

func (p *Producer) Push(b []byte, topic string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	_ = p.cmd.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:     b,
		Timestamp: time.Now(),
	}, nil)
}

func (p *Producer) RunTimer() {
	go func() {
		for e := range p.cmd.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					p.Stat.IncFail(1)
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					p.Stat.IncSuccess(1)
				}
			}
		}
	}()

	t := *time.NewTicker(60 * time.Second)
	go func() {
		for range t.C {
			fmt.Println("Success:", p.Stat.GetSuccess(), "Fail:", p.Stat.GetFail())
			p.Stat.Flush()
		}
	}()
}

func (p *Producer) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cmd.Close()
}

// -------------------------------------------------------------------------------------------  Stat structure functions
func (s *statistic) IncFail(deltaVal int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.fail += deltaVal
}

func (s *statistic) IncSuccess(deltaVal int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.success += deltaVal
}

func (s *statistic) GetSuccess() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.success
}

func (s *statistic) GetFail() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.fail
}

func (s *statistic) Flush() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.success = 0
	s.fail = 0
}
