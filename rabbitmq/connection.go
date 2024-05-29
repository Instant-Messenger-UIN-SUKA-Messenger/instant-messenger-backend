package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

var rabbitMQ *amqp.Connection

func ConnectToDatabase() error {
	fmt.Println("Try to connect to RabbitMQ!")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	rabbitMQ = conn
	fmt.Println("Successfully connected to RabbitMQ!")

	// Declare the exchange and queue after successful connection
	err = DeclareExchange("DatabaseExchange", "direct") // Replace with desired exchange name and type
	if err != nil {
		return err
	}
	err = DeclareQueue("DatabaseQueue") // Replace with desired queue name
	if err != nil {
		return err
	}
	// Bind the queue to the exchange
	BindQueueToExchange("DatabaseQueue", "DatabaseExchange", "database_key") // Replace with desired queue name, exchange name, and routing key

	// Start consuming messages after connection and setup
	ctx := context.Background() // Create a context
	go ConsumeFromDatabase(ctx) // Start consumer in a separate goroutine

	return nil
}

func ConnectToClient() error {
	fmt.Println("Try to connect to RabbitMQ!")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	rabbitMQ = conn
	fmt.Println("Successfully connected to RabbitMQ!")

	// Declare the exchange and queue after successful connection
	err = DeclareExchange("ClientExchange", "direct") // Replace with desired exchange name and type
	if err != nil {
		return err
	}
	err = DeclareQueue("ClientQueue") // Replace with desired queue name
	if err != nil {
		return err
	}
	// Bind the queue to the exchange
	BindQueueToExchange("ClientQueue", "ClientExchange", "client_key") // Replace with desired queue name, exchange name, and routing key

	// Start consuming messages after connection and setup
	ctx := context.Background() // Create a context
	go ConsumeFromClient(ctx)   // Start consumer in a separate goroutine

	return nil
}
