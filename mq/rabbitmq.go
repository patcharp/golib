package mq

import (
	"context"
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	VirtualHost string
	Exchange    string
}

type Client struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel

	vhost    string
	exchange string
	dialUrl  string
	job      *scheduler.Job
}

func New(cfg *Config) Client {
	dialUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	return Client{
		dialUrl:  dialUrl,
		vhost:    cfg.VirtualHost,
		exchange: cfg.Exchange,
	}
}

func (c *Client) Connect() error {
	_ = c.stopKeepAlive()
	conn, err := amqp.DialConfig(c.dialUrl, amqp.Config{
		Vhost:     c.vhost,
		Heartbeat: time.Second * 5,
		Locale:    "en_US",
	})
	if err != nil {
		return err
	}
	c.Connection = conn

	if err := c.createChannel(); err != nil {
		return err
	}
	if err = c.startKeepAlive(); err != nil {
		return err
	}
	return nil
}

func (c *Client) createChannel() error {
	if c.Channel == nil || (c.Channel != nil && c.Channel.IsClosed()) {
		ch, err := c.Connection.Channel()
		if err != nil {
			return err
		}
		c.Channel = ch
	}
	return nil
}

func (c *Client) Close() error {
	if c.Channel != nil && !c.Channel.IsClosed() {
		if err := c.Channel.Close(); err != nil {
			logrus.Warnln("[messageq] closing channel error ->", err)
		}
	}
	if c.Connection != nil && !c.Connection.IsClosed() {
		return c.Connection.Close()
	}
	_ = c.stopKeepAlive()
	return nil
}

func (c *Client) QueueDeclare(name string) (*amqp.Queue, error) {
	if err := c.createChannel(); err != nil {
		return nil, err
	}
	q, err := c.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (c *Client) PublishJSONMessage(qName string, msgId string, body []byte) error {
	return c.PublishMessage(qName, msgId, "", "", fiber.MIMEApplicationJSON, body)
}

func (c *Client) PublishTextMessage(qName string, msgId string, body []byte) error {
	return c.PublishMessage(qName, msgId, "", "", fiber.MIMETextPlain, body)
}

func (c *Client) PublishMessage(qName string, msgId string, appId string, userId string, mimeType string, body []byte) error {
	q, err := c.QueueDeclare(qName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return c.Channel.PublishWithContext(ctx,
		c.exchange,
		q.Name,
		false,
		false,
		amqp.Publishing{
			AppId:       appId,
			UserId:      userId,
			ContentType: mimeType,
			MessageId:   msgId,
			Timestamp:   time.Now(),
			Body:        body,
		},
	)
}

func (c *Client) startKeepAlive() error {
	var err error
	c.job, err = scheduler.Every(15).Seconds().Run(func() {
		if c.Connection == nil || (c.Connection != nil && c.Connection.IsClosed()) {
			logrus.Errorln("MQ keep alive error ->", err)
			if err := c.Connect(); err != nil {
				logrus.Errorln("Trying to reconnect to MQ error ->", err)
			} else {
				logrus.Infoln("MQ reconnect success.")
			}
		}
	})
	return err
}

func (c *Client) stopKeepAlive() error {
	if c.job != nil {
		c.job.Quit <- true
	}
	return nil
}
