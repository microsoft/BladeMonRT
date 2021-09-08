
package main

import (
	"C"
	"fmt"
)


func main() {
	client := NewAzPubSubSimpleClient(true, "10.0.0.155:9092")
	fmt.Println(client)
	client.sendMessage("test","test_message_1")
}