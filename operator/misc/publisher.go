package misc

import (
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
	writer *kafka.Writer
}

type KafkaConfig struct {
	BrokerEndpoint string
	Username string
	Password string
	SASLMechanism string
	TLSEnabled bool
}

func NewPublisher() *Publisher {
	pub := &Publisher{}
	conf := NewKafkaConfig()

	pub.writer = &kafka.Writer{
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
		SASLMechanism: GetEnvValue(fmt.Sprintf("%s_PUBLISHER_KAFKA_SASLMECHANISM", CONF_PREFIX), ""),
		TLSEnabled: GetEnvBool(fmt.Sprintf("%s_PUBLISHER_KAFKA_SASLMECHANISM", CONF_PREFIX), false),
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
