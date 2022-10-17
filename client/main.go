package main

import (
	"fmt"
	. "global"
	"log"
	"math/rand"
	"os"

	"github.com/streadway/amqp"
)

// Print usage description
func printUsage() {
	println("Client program send requests to the server via external message queue.")
	println()
	println("Usage:")
	println()
	println("    ./client <command> <key> <value>")
	println()
	println("The commands are:")
	println()
	fmt.Printf("    %s		Add item to the list.", ADD_ITEM)
	fmt.Printf("    %s		Rmove item from the list.", REMOVE_ITEM)
	fmt.Printf("    %s		Read item from the list.", GET_ITEM)
	fmt.Printf("    %s		Read all list.", GET_ALL_ITEMS)
	println("")
}

// Parse argument parameters
func parseArgs() (string, string, string) {
	var command string = ""
	var key string = ""
	var value string = ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 {
		key = os.Args[2]
	}
	if len(os.Args) > 3 {
		value = os.Args[3]
	}
	return command, key, value
}

// Generate random string with given length
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(rand.Intn(26) + 65)
	}
	return string(bytes)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Program entry point
func main() {

	// Get arguments from the command line
	if len(os.Args) == 1 {
		printUsage()
		return
	}
	command, key, value := parseArgs()

	// Create connection to the RabbitMQ server
	conn, err := amqp.Dial(MQ_SERVER_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create channel to communicate with RabbitMQ
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declear queue in the channel
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create consumer against the message queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Publish request to the MQ
	corrId := randomString(32)
	err = ch.Publish(
		"",             // exchange
		MSG_QUEUE_NAME, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          EncodeToBytes(MessageBody{Command: command, Key: key, Value: value}),
		})
	failOnError(err, "Failed to publish a message")

	// Wait for the server response
	var buffer []byte = nil
	for d := range msgs {
		if corrId == d.CorrelationId {
			buffer = d.Body
			break
		}
	}

	if buffer == nil {
		return
	}

	// Process server response
	switch command {
	// Response for add_item command
	case ADD_ITEM:
		resp := DecodeToAddItemResponse(buffer)
		if resp.Success {
			println("Succeed")
		} else {
			println("Failed")
		}
	// Response for remove_item command
	case REMOVE_ITEM:
		resp := DecodeToRemoveItemResponse(buffer)
		if resp.Success {
			println("Succeed")
		} else {
			println("Failed")
		}
	// Response for get_item command
	case GET_ITEM:
		resp := DecodeToGetItemResponse(buffer)
		if resp.Success {
			println(resp.Item)
		} else {
			println("Failed")
		}
	// Response for get_all_items command
	case GET_ALL_ITEMS:
		resp := DecodeToGetAllItemsResponse(buffer)
		log.Print(resp)
	default:

	}
}
