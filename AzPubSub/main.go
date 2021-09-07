
package main

import (
	"C"
	"fmt"
)


func main() {
	client := NewAzPubSubSimpleClient(false, "127.0.0.1:50458")
	fmt.Println(client)
	client.sendMessage("test","test_message_1")
}