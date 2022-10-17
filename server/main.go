package main

import (
	"fmt"
	. "global"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

var store = orderedmap.NewOrderedMap[string, string]()
var store_mutex sync.Mutex

func failOnError(err error, msg string) {
	if err != nil {
		log_error(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Process message retrieved from the MQ and return response buffer
func processMessage(msg MessageBody) []byte {
	switch msg.Command {
	// Add item to the ordered map
	case ADD_ITEM:
		store_mutex.Lock()
		store.Set(msg.Key, msg.Value)
		store_mutex.Unlock()
		return EncodeToBytes(AddItemResponse{Success: true})

	// Remove item from the ordered map
	case REMOVE_ITEM:
		store_mutex.Lock()
		store.Delete(msg.Key)
		store_mutex.Unlock()
		return EncodeToBytes(RemoveItemResponse{Success: true})

	// Get item for the key
	case GET_ITEM:
		store_mutex.Lock()
		value, err := store.Get(msg.Key)
		store_mutex.Unlock()
		return EncodeToBytes(GetItemResponse{Success: err, Item: value})

	// Get all items in the store
	case GET_ALL_ITEMS:
		store_mutex.Lock()
		data := make([]KeyValuePair, 0)
		for _, key := range store.Keys() {
			value, _ := store.Get(key)
			data = append(data, KeyValuePair{Key: key, Value: value})
		}
		store_mutex.Unlock()
		return EncodeToBytes(GetAllItemsResponse{Items: data})
	default:
		log_error("Unknown command")
	}
	return nil
}

// Print usage description
func printUsage() {
	println("Server program reads requests and execute them on its data structure.")
	println()
	println("Usage:")
	println()
	println("    ./server [-t] [thread_count] [-d] [delay] [-l] [log_file]")
	println()
	println("The arguments are:")
	println()
	println("    thread_count:	Thread count server can concurrently run, defaults to 4")
	println("    delay:			Response delay for every message in milliseconds, defaults to 0")
	println("    log_file:		Log file path, defaults to stdout")
	println("")
}

// Parse argument parameters
func parseArgs() (string, int, int) {
	var log_file_path string = ""
	var delay int = 0
	var thread_count int = 4
	var err error

	for i := range os.Args {
		if i == 0 {
			continue
		}
		// Parse thread_count from arguments
		if os.Args[i] == "-t" {
			if i < len(os.Args)-1 {
				thread_count, err = strconv.Atoi(os.Args[i+1])
				if err != nil {
					println("Unable to get thread_count.")
				}
			} else {
				println("Thread count not set, set to default.")
			}
		}

		// Parse delay from arguments
		if os.Args[i] == "-d" {
			if i < len(os.Args)-1 {
				delay, err = strconv.Atoi(os.Args[i+1])
				if err != nil {
					println("Unable to get delay.")
				}
			} else {
				println("Delay not set, set to default.")
			}
		}

		// Parse log_file_path from arguments
		if os.Args[i] == "-l" {
			if i < len(os.Args)-1 {
				log_file_path = os.Args[i+1]
			} else {
				println("Log file path not set, set to default.")
			}
		}
	}

	return log_file_path, delay, thread_count
}

// Program entry point
func main() {

	// Get arguments from the command line
	if len(os.Args) == 1 {
		printUsage()
		return
	}
	log_file_path, delay, thread_count := parseArgs()

	// Init logger for stdout and file
	init_logger(log_file_path)

	// output channel for synchronizing routine creation and result retrival
	output_chan := make(chan *any, thread_count)

	// Fill the output channel with arbitrary values for the first routine creation
	for i := 0; i < thread_count; i++ {
		output_chan <- nil
	}

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
		MSG_QUEUE_NAME, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Prevent round robin method and set prefetch count
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

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

	go func() {

		for d := range msgs {
			// Decode message buffer to the struct
			body := []byte(d.Body)
			msg := DecodeToMessage(body)

			// Pop single result from ouput channel so that another routine is able to be created.
			<-output_chan

			go func(d amqp.Delivery, msg MessageBody) {
				// Sleep for given milliseconds to simulate heavy calculation
				time.Sleep(time.Duration(delay) * time.Millisecond)

				// Process message and retrieve response
				log_info(msg)
				response := processMessage(msg)
				log_info(response)

				// Publish response buffer to the channel
				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          response,
					})
				failOnError(err, "Failed to publish a message")

				// Allow other threads to be created
				output_chan <- nil

			}(d, msg)
		}
	}()

	// Wait forever unless compulsary interrupt
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever

}
