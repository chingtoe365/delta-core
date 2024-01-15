package controller

import (
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/mqttutil"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	Env *bootstrap.Env
}

// PingExample godoc
// @Summary Test send message to MQTT
// @Schemes
// @Description Test message sending to MQTT mosquitto
// @Tags Test
// @Param msg query string true "Message to be sent"
// @Param topic query string true "Topic where the msg is sent to in MQTT (eg. topic/test)"
// @Accept json
// @Produce json
// @Success 200
// @Router /test [get]
func (tc *TestController) Test(c *gin.Context) {
	var env *bootstrap.Env = c.MustGet("env").(*bootstrap.Env)
	p := domain.Profile{
		Name: "foo", Email: "bar",
	}
	client := mqttutil.NewMqttClient(env, &p)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	msg := c.Query("msg")
	topic := c.Query("topic")
	fmt.Printf(">> Publishing to topic: %s\n", topic)
	token := client.Publish(topic, 2, false, msg)
	token.Wait()
	client.Disconnect(250)
}
