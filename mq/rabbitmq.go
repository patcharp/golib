package mq

import (
	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/patcharp/golib/v2/util/httputil"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	VirtualHost string
	Channel     string
}

type Client struct {
	Config  Config
	Ctx     *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
	job     *scheduler.Job
}

func NewMQ(cfg Config) Client {
	return Client{
		Config: cfg,
	}
}

func (c *Client) Connect(qName string) error {
	var err error
	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.Config.Username,
		c.Config.Password,
		c.Config.Host,
		c.Config.Port,
	)
	c.Ctx, err = amqp.DialConfig(connStr, amqp.Config{
		Vhost: c.Config.VirtualHost,
	})
	if err != nil {
		return err
	}
	// Get channel
	c.Channel, err = c.Ctx.Channel()
	if err != nil {
		return err
	}
	// Create queue
	c.Queue, err = c.Channel.QueueDeclare(
		qName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	if err := c.startKeepAlive(qName); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if c.Ctx != nil && !c.Ctx.IsClosed() {
		return c.Ctx.Close()
	}
	_ = c.stopKeepAlive()
	return nil
}

func (c *Client) EnQueue(key string, queueId string, exchange *string, data []byte) error {
	exchangeCfg := ""
	if exchange != nil {
		exchangeCfg = *exchange
	}
	return c.Channel.Publish(
		exchangeCfg,  // exchange
		c.Queue.Name, // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  httputil.MIMEApplicationJSON,
			Body:         data,
			MessageId:    queueId,
			Timestamp:    time.Now(),
			AppId:        key,
		},
	)
}

func (c *Client) startKeepAlive(qName string) error {
	var err error
	c.job, err = scheduler.Every(15).Seconds().Run(func() {
		if c.Ctx == nil || (c.Ctx != nil && c.Ctx.IsClosed()) {
			logrus.Errorln("MQ keep alive error ->", err)
			if err := c.Connect(qName); err != nil {
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
