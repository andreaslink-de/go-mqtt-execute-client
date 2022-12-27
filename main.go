/*
Author: Dipl.-Inform. (FH) Andreas Link, andreaslink.de // link-tech.de
Date: 27.12.2022

--------------------------------------------------------------------------------

To modify the key values, compile/build from code, deploy and finally run on a RaspberryPi do the following steps after cloneing the project:

  // Init module mode first:
 	go mod init main
 	go mod tidy

  // Download dependencies:
	go mod download github.com/gorilla/websocket
	go mod download github.com/eclipse/paho.mqtt.golang
	go get github.com/eclipse/paho.mqtt.golang@v1.4.2

  // Adjust topic to subscribe to

  // Build for RaspberryPi ARM using:
    env GOOS=linux GOARCH=arm GOARM=5 go build -o mqtt-hdmi-controller main.go

  // Copy tool to Raspberry Pi '/usr/local/bin/mqtt-hdmi-controller' and setup as a systemd service by creating "/etc/systemd/system/mqtt-hdmi-controller.service":
		[Unit]
		Description=RaspberryPi HDMI-Monitor ON/OFF based on MQTT command
		After=network.target

		[Service]
		Type=simple
		ExecStart=/usr/local/bin/mqtt-hdmi-controller
		WorkingDirectory=/home/pi
		Restart=always

		[Install]
		WantedBy=multi-user.target

	// Enable it: sudo systemctl enable mqtt-hdmi-controller.service
	// And finally start it: sudo systemctl start mqtt-hdmi-controller.service
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Create an MQTT client
	client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://192.168.42.253:1883"))

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to MQTT broker")

	// Subscribe to the topic: zuhause/haus/esszimmer/infomonitor/bildschirm/status --> 0 = OFF; 1 = ON
	if token := client.Subscribe("zuhause/haus/esszimmer/infomonitor/bildschirm/status", 0, func(client mqtt.Client, msg mqtt.Message) {
		// This is the message handler function. It will be called every time a message is received on the subscribed topic.
		fmt.Printf("Received message on topic: %s\n", msg.Topic())
		fmt.Printf("Message payload: %s\n", string(msg.Payload()))

		// Execute a command based on the value of the message payload
		switch string(msg.Payload()) {
		case "1":
			// Execute a system command on Bash:
			// Turn RasPi-Monitor ON:  vcgencmd display_power 1 2>&1 | /usr/bin/logger -t "HDMI-Display ON"
			// Turn Raspi-Monitor OFF: vcgencmd display_power 0 2>&1 | /usr/bin/logger -t "HDMI-Display OFF"
			cmd := exec.Command("bash", "-c", "vcgencmd display_power 1 2>&1 | /usr/bin/logger -t 'HDMI-Display ON'")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Command output: %s\n", string(output))
		case "0":
			// Execute a system command on Bash:
			// Turn RasPi-Monitor ON:  vcgencmd display_power 1 2>&1 | /usr/bin/logger -t "HDMI-Display ON"
			// Turn Raspi-Monitor OFF: vcgencmd display_power 0 2>&1 | /usr/bin/logger -t "HDMI-Display OFF"
			cmd := exec.Command("bash", "-c", "vcgencmd display_power 0 2>&1 | /usr/bin/logger -t 'HDMI-Display OFF'")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Command output: %s\n", string(output))
		}
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Keep the program running indefinitely
	for {
		time.Sleep(1 * time.Second)
	}
}
