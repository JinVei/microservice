package kafka

import (
	"context"
	"os"
	"testing"

	"time"

	"github.com/jinvei/microservice/base/framework/configuration"
)

func TestProducer(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("11001")
	kpconf, err := MakeProducerConfig(conf, "")
	if err != nil {
		t.Fatal(err)
	}
	producer := NewProducer(*kpconf)
	flog.Info("test", "producer", producer)
	err = producer.Send(context.Background(), "", "TestFromMc"+time.Now().String())
	if err != nil {
		t.Fatal(err)
	}
}

type testReceiver struct {
}

func (c *testReceiver) OnReceive(topic string, partition int, Offset int64, key, value []byte) error {
	flog.Info("OnReceive:", "topic", topic, "partition", partition, "Offset", Offset, "key", string(key), "value", string(value))
	return nil
}

func TestConsumer(t *testing.T) {
	os.Setenv("MICROSERVICE_CONFIGURATION_TOKEN", "e30K")
	conf := configuration.DefaultOrDie()
	conf.SetSystemID("11001")
	kpconf, err := MakeConsumerConfig(conf, "")
	if err != nil {
		t.Fatal(err)
	}

	receiver := new(testReceiver)

	consumer := NewConsumer(*kpconf, receiver)

	ctx := context.TODO()
	err = consumer.Start(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Second)

	ctx.Done()

	consumer.Close()

	time.Sleep(3 * time.Second)
}

func TestConsumerA(t *testing.T) {

	// test without consumer group
	receiver := new(testReceiver)

	kpconf := &ConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test",
		// GroupID: "",
	}

	consumer := NewConsumer(*kpconf, receiver)

	ctx := context.TODO()
	err := consumer.Start(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Second)

	ctx.Done()

	consumer.Close()

	time.Sleep(3 * time.Second)
}
