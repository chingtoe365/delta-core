package mqttutil

import (
	"crypto/rand"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/notificationutil"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func messageHandlerWrapper(env *bootstrap.Env, p *domain.Profile) func(client mqtt.Client, msg mqtt.Message) {
	// return message handler
	return func(client mqtt.Client, msg mqtt.Message) {
		var options = client.OptionsReader()
		log.Printf("ClientId: %s, Received message: %s, from topic: %s\n", options.ClientID(), msg.Payload(), msg.Topic())
		var a domain.Alert
		a.ParseIn(string(msg.Payload()), msg.Topic())
		notificationutil.SendMail(env, p.Email, a.FormatEmail())
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	slog.Info("Connected to mosquitto")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewMqttClient(env *bootstrap.Env, profile *domain.Profile) mqtt.Client {
	var broker = env.MqttHost
	var port = env.MqttPort
	var user = env.MqttUser
	var password = env.MqttPass
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messageHandlerWrapper(env, profile))
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
