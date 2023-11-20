package bootstrap

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var options = client.OptionsReader()
	fmt.Printf("ClientId: %s, Received message: %s, from topic: %s\n", options.ClientID(), msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to mosquitto")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewMqttClient(env *Env) mqtt.Client {
	var broker = env.MqttHost
	var port = env.MqttPort
	var user = env.MqttUser
	var password = env.MqttPass
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messageHandler)
	clientId, _ := randomHex(20)
	opts.SetClientID(clientId)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
