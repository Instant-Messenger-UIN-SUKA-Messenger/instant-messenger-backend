package rabbitmq

import "github.com/streadway/amqp"

var rabbitMQ *amqp.Connection
var rabbitMQChannel *amqp.Channel