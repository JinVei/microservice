package kafka

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/jinvei/microservice/base/framework/configuration"
	confkeys "github.com/jinvei/microservice/base/framework/configuration/keys"
	"github.com/jinvei/microservice/base/framework/log"
	kg "github.com/segmentio/kafka-go"
)

var flog = log.Default

type ProducerConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type Producer struct {
	w   *kg.Writer
	cfg ProducerConfig
}

type ConsumerConfig struct {
	Brokers []string `json:"brokers"`
	GroupID string   `json:"groupid"`
	Topic   string   `json:"topic"`
}

type Consumer struct {
	r        *kg.Reader
	cfg      ConsumerConfig
	callback Receiver
}

type Receiver interface {
	// Receiver callback
	// NOTE: Consumer would not Commit Messages if return error
	OnReceive(topic string, partition int, Offset int64, key, value []byte) error
}

// MakeProducerConfig make ProducerConfig from Configuration Store.
// It can select config piece by specify `area`. this option also can be omit.
// configuration path in Store: `/microservice/framework/kafka/producer/{systemID}/{area}`
// example value in Configuration Store: `{"brokers":["localhost:9092","localhost:9093"], "topic": "test"}`
func MakeProducerConfig(conf configuration.Configuration, area string) (*ProducerConfig, error) {
	sid := conf.GetSystemID()
	cfg := &ProducerConfig{}
	if err := conf.GetJson(filepath.Join(confkeys.FwkafkaProducer, sid, area), cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// MakeConsumerConfig make ConsumerConfig from Configuration Store.
// It can select config piece by specify `area`. this option also can be omit.
// configuration path in Store: `/microservice/framework/kafka/consumer/{systemID}/{area}`
// example value in Configuration Store: `{"brokers":["localhost:9092"],"groupid": "test_group", "topic": "test"}`
func MakeConsumerConfig(conf configuration.Configuration, area string) (*ConsumerConfig, error) {
	sid := conf.GetSystemID()
	cfg := &ConsumerConfig{}
	if err := conf.GetJson(filepath.Join(confkeys.FwKafkaConsumer, sid, area), cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func NewProducer(cfg ProducerConfig) *Producer {
	writer := &kg.Writer{
		Addr:        kg.TCP(cfg.Brokers...),
		Balancer:    kg.Murmur2Balancer{},
		Compression: kg.Snappy,
		Topic:       cfg.Topic,
	}

	return &Producer{
		w:   writer,
		cfg: cfg,
	}
}

func (p *Producer) Send(ctx context.Context, key, content string) error {
	var keyb []byte = nil
	if key != "" {
		keyb = []byte(key)
	}

	msg := kg.Message{
		Value: []byte(content),
		Key:   keyb,
		Time:  time.Now(),
	}
	flog.Debug("WriteMessages", "msg", msg)
	if err := p.w.WriteMessages(ctx, msg); err != nil {
		flog.Error(err, "w.WriteMessages", "msg", msg)
		return err
	}
	return nil
}

func (p *Producer) GetWriter() *kg.Writer {
	return p.w
}

func NewConsumer(cfg ConsumerConfig, cb Receiver) *Consumer {
	reader := kg.NewReader(kg.ReaderConfig{
		Brokers:               cfg.Brokers,
		GroupID:               cfg.GroupID,
		Topic:                 cfg.Topic,
		WatchPartitionChanges: true,
		MinBytes:              2e3, // 10KB
		MaxBytes:              5e6, // 10MB
	})

	return &Consumer{
		r:        reader,
		cfg:      cfg,
		callback: cb,
	}
}

func (c *Consumer) Start(ctx context.Context, offset int64) error {
	if offset != 0 {
		if err := c.r.SetOffset(offset); err != nil {
			return err
		}
	}

	flog.Info("Starting kafka consumer", "ConsumerConfig", c.cfg, "offset", offset)

	go func() {
		for {
			msg, err := c.r.FetchMessage(ctx)
			if err != nil {
				if err == io.EOF {
					flog.Info("Exit Consumer")
					return
				}
				flog.Error(err, "FetchMessage msg err")
				continue
			}

			if err := c.callback.OnReceive(msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value); err != nil {
				// flog.Error(err, "OnReceive() return err")
				continue
			}
			// Should CommitMessages manually if GroupID is set
			if c.cfg.GroupID != "" {
				if err := c.r.CommitMessages(ctx, msg); err != nil && err != io.EOF {
					flog.Error(err, "CommitMessages err")
				}
			}
		}
	}()

	return nil
}

func (c *Consumer) GetReader() *kg.Reader {
	return c.r
}

// Note that it is important to call Close() on a Reader when a process exits.
// This can result in a delay when a new reader on the same topic connects if not call Close() in process exit
func (c *Consumer) Close() error {
	return c.r.Close()
}
