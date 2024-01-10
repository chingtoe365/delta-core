package domain

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

type Signal struct {
	Short       string
	Description string
}

type Alert struct {
	TradeItem string
	Signal    Signal
	Time      time.Time
}

func (a *Alert) ParseIn(shortDesc string, topic string) {
	var topicSplits = strings.Split(topic, "/")
	// time, err := time.Parse("%", timePayload)
	// if err != nil {
	// 	panic(err)
	// }
	a.TradeItem = topicSplits[0]
	a.Signal = Signal{
		Short:       topicSplits[1],
		Description: shortDesc,
	}
	a.Time = time.Now()
}

func (a *Alert) FormatEmail() string {
	tmpl, err := template.New("email").Parse(
		`
		Trade Item: {{.TradeItem}}
		Signal Short Name: {{.Signal.Short}}
		Singal Description: {{.Signal.Description}}
		Detected Time: {{.Time}}
		`,
	)
	if err != nil {
		panic(err)
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, a)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
