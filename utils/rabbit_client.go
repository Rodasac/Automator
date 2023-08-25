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

func StartClient(log *otelzap.LoggerWithCtx, consumerName string) (*Consumer, error) {
	uri := os.Getenv("RABBITMQ_URI")
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	exchangeType := os.Getenv("RABBITMQ_EXCHANGE_TYPE")
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	bindingKey := os.Getenv("RABBITMQ_BINDING_KEY")
	if uri == "" || exchange == "" || exchangeType == "" || queueName == "" || bindingKey == "" {
		return nil, fmt.Errorf("environment variables for rabbit not set")
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

	log.Debug("declared Exchange, declaring Queue", zap.String("queue", queueName))
	queue, err := c.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %s", err)
	}

	log.Debug(
		"declared Queue, binding to Exchange",
		zap.String("exchange", exchange),
		zap.String("queue", queueName),
		zap.String("bindingKey", bindingKey),
	)

	if err = c.Channel.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue bind: %s", err)
	}

	return c, nil
}
