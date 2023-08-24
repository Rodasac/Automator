package utils

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"os"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	log     *otelzap.LoggerWithCtx
	tag     string
	done    chan error
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer c.log.Debug("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func StartConsumer(log *otelzap.LoggerWithCtx, consumerName string) (*Consumer, error) {
	uri := os.Getenv("RABBITMQ_URI")
	if uri == "" {
		log.Fatal("RABBITMQ_URI is required")
	}
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	if exchange == "" {
		log.Fatal("RABBITMQ_EXCHANGE is required")
	}
	exchangeType := os.Getenv("RABBITMQ_EXCHANGE_TYPE")
	if exchangeType == "" {
		log.Fatal("RABBITMQ_EXCHANGE_TYPE is required")
	}

	c := &Consumer{
		Conn:    nil,
		Channel: nil,
		log:     log,
		tag:     consumerName,
		done:    make(chan error),
	}
	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(consumerName)

	log.Debug("dialing rabbitmq", zap.String("uri", uri))
	c.Conn, err = amqp.DialConfig(uri, config)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	go func() {
		log.Debug("closing rabbit: %s", zap.NamedError("reason", <-c.Conn.NotifyClose(make(chan *amqp.Error))))
	}()

	log.Debug("got Connection, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	log.Debug("got Channel, declaring Exchange", zap.String("exchange", exchange))
	if err = c.Channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("exchange declare: %s", err)
	}

	return c, nil
}
