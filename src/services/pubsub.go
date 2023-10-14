package pubsub

import (
	"delta-core/domain"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var options = client.OptionsReader()
	fmt.Printf("clientId: %s, Received message: %s, from topic: %s\n", options.ClientID(), msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func CreateSubClient(taskId string, topic string) mqtt.Client {
	var broker = "localhost"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetClientID(taskId)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client, topic)
	// client.Disconnect(250)
	return client
}

// func CancelSubClient(taskId string, topic string) mqtt.Client {
// }

// called when server start up
func InitialiseSubClients(tasks []domain.Task) {
	for _, item := range tasks {
		go CreateSubClient(item.ID.Hex(), item.Title)
	}
}

// to publish a message to this Client
func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

// Sub to all topics, and we filter
func sub(client mqtt.Client, topic string) {
	// if callback to process received message3rd is nill, the DefaultPublishHandler is used.
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", topic)
}
