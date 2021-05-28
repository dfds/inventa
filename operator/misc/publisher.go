package misc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"net"
	"strings"
	"time"
)

type Publisher struct {
	Writer *kafka.Writer
	Topic string
}

type KafkaConfig struct {
	BrokerEndpoint string
	Username string
	Password string
	SASLMechanism string
	TLSEnabled bool
	Topic string
}

func NewPublisher() *Publisher {
	pub := &Publisher{}
	conf := NewKafkaConfig()

	pub.Topic = conf.Topic
	pub.Writer = &kafka.Writer{
		Addr:		kafka.TCP(conf.BrokerEndpoint),
		Async:		false,
		Transport:	MakeKafkaTransport(conf),
		Completion: func(messages []kafka.Message, err error) {
			log.Printf("%v\n", messages)
		},
	}

	return pub
}

func NewKafkaConfig() KafkaConfig {
	conf := KafkaConfig{
		Username: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_USERNAME", CONF_PREFIX), ""),
		Password: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_PASSWORD", CONF_PREFIX), ""),
		BrokerEndpoint: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_ENDPOINT", CONF_PREFIX), ""),
		Topic: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_TOPIC", CONF_PREFIX), ""),
		SASLMechanism: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_SASLMECHANISM", CONF_PREFIX), ""),
		TLSEnabled: GetEnvBool(fmt.Sprintf("%s_PUBLISHER_KAFKA_TLSENABLED", CONF_PREFIX), false),
	}

	return conf
}

func MakeKafkaTransport(conf KafkaConfig) *kafka.Transport {
	var saslMechanism sasl.Mechanism = nil
	var tlsConf *tls.Config = nil

	// Support for SCRAM can be added here.
	switch confSasl := strings.ToLower(conf.SASLMechanism); confSasl {
	case "plain":
		saslMechanism = plain.Mechanism{
			Username: conf.Username,
			Password: conf.Password,
		}
	}

	if conf.TLSEnabled {
		tlsConf = &tls.Config{}
	}

	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
		SASLMechanism: saslMechanism,
		TLS: tlsConf,
	}

	netDialer := &net.Dialer{
		Timeout: dialer.Timeout,
		Deadline: dialer.Deadline,
		LocalAddr: dialer.LocalAddr,
		DualStack: dialer.DualStack,
		FallbackDelay: dialer.FallbackDelay,
		KeepAlive: dialer.KeepAlive,
	}

	transport := &kafka.Transport{
		Dial: netDialer.DialContext,
		SASL: dialer.SASLMechanism,
		TLS: dialer.TLS,
		ClientID: dialer.ClientID,
		IdleTimeout: 9 * time.Minute, // See segmentio/kafka-go/writer.go for why this value was chosen
		MetadataTTL: 15 * time.Second, // See segmentio/kafka-go/writer.go for why this value was chosen
	}

	return transport
}

type PublisherService struct {
	publisher *Publisher
	messageChannel <-chan string
}

func NewPublisherService(messageChannel <-chan string) *PublisherService {
	ps := &PublisherService{
		publisher: NewPublisher(),
		messageChannel: messageChannel,
	}

	return ps
}

func (p *PublisherService) Run() {
	for {
		msg := <-p.messageChannel

		err := p.publisher.Writer.WriteMessages(context.Background(),
			kafka.Message{
				Topic: p.publisher.Topic,
				Value: []byte(fmt.Sprintf(msg)),
			})
		if err != nil {
			log.Println(err)
		}

		fmt.Println(msg)
	}
}

func RunPublisherService(messageChannel <-chan string) {
	ps := NewPublisherService(messageChannel)
	ps.Run()
}
