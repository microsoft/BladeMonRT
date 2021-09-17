
package main

import (
	"log"
	"fmt"
)

const (
	PfAzPubSubVipFile = "test_PF_Pub_Sub_VIP_file.txt" //`D:\data\NodeServiceSettings.ini.flattened.ini`
)

func main() {
	vip, err := NewUtils().FetchAzPubSubPfVIP(PfAzPubSubVipFile)
	if (err != nil) {
		log.Fatal("Failed to get VIP", err)
	}

	client := NewAzPubSubSimpleClient(false, vip)
	fmt.Println(client)
	client.SendMessage("test","test_message_1")

	// globalClient := NewAzPubSubGlobalClient(true, "10.0.0.155:9092")
	// fmt.Println(globalClient)
	// client.sendMessage("test","test_message_1")
}