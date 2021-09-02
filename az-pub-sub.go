package main

import "C"
import (
	"syscall"
	"log"
	"fmt"
)

type (
	BOOL          uint32
	BOOLEAN       byte
	BYTE          byte
	DWORD         uint32
	DWORD64       uint64
	PVOID uintptr
	HANDLE        PVOID
	HLOCAL        uintptr
	LARGE_INTEGER int64
	LONG          int32
	LPVOID        uintptr
	SIZE_T        uintptr
	UINT          uint32
	ULONG_PTR     uintptr
	ULONGLONG     uint64
	WORD          uint16
)

var (
	wevtapi = syscall.NewLazyDLL("azpubsub.dll")
	AzPubSubSendMessageEx              = wevtapi.NewProc("AzPubSubSendMessageEx")
	AzPubSubOpenSimpleProducer              = wevtapi.NewProc("AzPubSubOpenSimpleProducer")
	AzPubSubCreateConfiguration              = wevtapi.NewProc("AzPubSubCreateConfiguration")
	AzPubSubClientInitialize              = wevtapi.NewProc("AzPubSubClientInitialize")
	NULL = HANDLE(0)
)

const (
	AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_NONE = 0

	AZPUBSUB_CONFIGURATION_TYPE_SIMPLE = 2
)

type AzPubSubClient struct {
	isTestInstance bool
	apsSecurityType int
	apsConnectionFlags int
	hclient HANDLE
	hconfig HANDLE
	endpoint string
	hproducer uintptr
	apsConfigType int
}

type loggerCallback func(level int, message, context uintptr) uintptr

func pLoggerCallback(level int, message, context uintptr) uintptr {
	fmt.Println(fmt.Sprintf("Log: msg=%s at level=%s", message, level))
	return uintptr(0)
}

func (client *AzPubSubClient) InitClient() {
	fmt.Println("Initializing client.")
	var err error
	client.hclient, err = initClient(pLoggerCallback)
	if (client.hclient == NULL) {
		log.Fatal(err)
	}
}

type AzPubSubSimpleClient struct {
	AzPubSubClient
}

func NewAzPubSubSimpleClient() *AzPubSubSimpleClient {
	var client AzPubSubSimpleClient = AzPubSubSimpleClient{AzPubSubClient: AzPubSubClient{isTestInstance: true,
		endpoint : "127.0.0.1"}}
	client.apsConfigType = AZPUBSUB_CONFIGURATION_TYPE_SIMPLE
	client.InitClient()
	client.InitConfig()

	return &client
}

func (client *AzPubSubClient) InitConfig() {
	fmt.Println("Initializing config.")
	var err error
	// Python version
	// self.hconfig = AzPubSub.AzPubSubCreateConfiguration(self.hclient, self.aps_config_type, AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES.NONE)
	// ========================
	// C version
	// AzPubSubCreateConfiguration(HCLIENT hClient, AZPUBSUB_CONFIGURATION_TYPE type, 
	// AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES globalConfigTemplate)
	// ============
	// Constants
	// AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES.NONE = 0
	// AZPUBSUB_CONFIGURATION_TYPE.SIMPLE = 2	
	client.hconfig, err = createConfiguration(client.hclient, AZPUBSUB_CONFIGURATION_TYPE_SIMPLE, AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_NONE)
	if (client.hconfig == NULL) {
		log.Fatal("Initializing config error: ", err)
	}
}

func initClient(callback loggerCallback) (HANDLE, error) {
	hclient, _, err := AzPubSubClientInitialize.Call(0, syscall.NewCallback(pLoggerCallback), uintptr(0))
	return HANDLE(hclient), err
}

func createConfiguration(client HANDLE, configType int, globalConfigTemplate int) (HANDLE, error) {
	fmt.Println(fmt.Sprintf("The value of parameters is client=%d, configType=%d, globalConfigTemplate=%d", client, AZPUBSUB_CONFIGURATION_TYPE_SIMPLE, AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_NONE))
	hconfig, _, err := AzPubSubCreateConfiguration.Call(0, uintptr(client), uintptr(configType), uintptr(globalConfigTemplate))
	return HANDLE(hconfig), err
}