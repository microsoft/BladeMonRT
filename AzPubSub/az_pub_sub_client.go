package azpubsub

import "C"
import (
	"syscall"
	"log"
)

type (
	BOOL          uint32
	BOOLEAN       byte
	BYTE          byte
	DWORD         uint32
	DWORD64       uint64
	HANDLE        uintptr
	HLOCAL        uintptr
	LARGE_INTEGER int64
	LONG          int32
	LPVOID        uintptr
	SIZE_T        uintptr
	UINT          uint32
	ULONG_PTR     uintptr
	ULONGLONG     uint64
	WORD          uint16
	LPCSTR        *int8
	INT			  int32

	// SimpleAzPubSubClient specific conversions.
	ENUM 		  int
	AZPUBSUB_SECURITY_TYPE int
	HCLIENT       HANDLE
	HCONFIG		  HANDLE
	HPRODUCER	  HANDLE
	HRESPONSE	  HANDLE
	LOG_LEVEL     int64
	LPPSTR        *[]byte
	PDWORD         *DWORD    
	PINT          *INT  
	
	// GlobalAzPubSubClient specific conversions
	HPRODUCERTOPIC HANDLE
	PBYTE *BYTE
	
)

// TODO: Add use of GetLastError

var (
	// C:\Users\t-nshanker\source\repos\BladeMonRT\azpubsub\azpubsub.dll full path does not work
	wevtapi = syscall.NewLazyDLL(`azpubsub.dll`)
	AzPubSubSendMessageEx              = wevtapi.NewProc("AzPubSubSendMessageEx")
	AzPubSubOpenSimpleProducer              = wevtapi.NewProc("AzPubSubOpenSimpleProducer")
	AzPubSubCreateConfiguration              = wevtapi.NewProc("AzPubSubCreateConfiguration")
	AzPubSubClientInitialize              = wevtapi.NewProc("AzPubSubClientInitialize")
	AzPubSubResponseGetMessage              = wevtapi.NewProc("AzPubSubResponseGetMessage")
	AzPubSubResponseGetStatusCode              = wevtapi.NewProc("AzPubSubResponseGetStatusCode")
	AzPubSubResponseGetSubStatusCode              = wevtapi.NewProc("AzPubSubResponseGetSubStatusCode")
	AzPubSubOpenProducerTopic              = wevtapi.NewProc("AzPubSubOpenProducerTopic")
	AzPubSubAddStringConfiguration = wevtapi.NewProc("AzPubSubAddStringConfiguration")
	NULL = HANDLE(0)
)

const (
	AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_NONE = 0
	AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_LOW_LATENCY = 1
	AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_MEDIUM_LATENCY = 2
	AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_HIGH_LATENCY = 3

	AZPUBSUB_CONFIGURATION_TYPE_GLOBAL = 0
	AZPUBSUB_CONFIGURATION_TYPE_TOPIC = 1
	AZPUBSUB_CONFIGURATION_TYPE_SIMPLE = 2

	AZPUBSUB_SECURITY_TYPE_NONE =  0
	AZPUBSUB_SECURITY_TYPE_SSL = 1
	AZPUBSUB_SECURITY_TYPE_TOKEN = 2

	AZPUBSUB_SECURITY_FLAGS_NONE = 0
	AZPUBSUB_SECURITY_FLAGS_LOCAL = 1
	
	STATUS_OK = 0
 	STATUS_MORE_DATA = 234
	ERR_OK = "The operation completed successfully."

)

type AzPubSubClient struct {
	isTestInstance bool
	apsSecurityType AZPUBSUB_SECURITY_TYPE
	apsConnectionFlags int
	apsConfigType int
	hclient HCLIENT
	hconfig HCONFIG
	hproducer HPRODUCER
	endpoint string
}

func NewAzPubSubClient(isTestInstance bool, endpoint string) AzPubSubClient {
	var client AzPubSubClient = AzPubSubClient{isTestInstance: isTestInstance, endpoint: endpoint}
	if (isTestInstance) {
		log.Println("Initializing client with test globals.")
		client.apsSecurityType = AZPUBSUB_SECURITY_TYPE_NONE
		client.apsConnectionFlags = AZPUBSUB_SECURITY_FLAGS_LOCAL
	} else {
		log.Println("Initializing client with production globals.")
		client.apsSecurityType = AZPUBSUB_SECURITY_TYPE_SSL
		client.apsConnectionFlags = AZPUBSUB_SECURITY_FLAGS_NONE
	}
	return client
}

type AZPUBSUB_LOG_CALLBACK func(level LOG_LEVEL, message LPCSTR, context LPVOID) uintptr

type AZPUBSUB_TOPIC_PARTITIONER func(arg1 DWORD, arg2 HPRODUCERTOPIC, arg3 PBYTE, arg4 DWORD, arg5 DWORD, arg6 LPVOID) 

func pLoggerCallback(level LOG_LEVEL, message LPCSTR, context LPVOID) uintptr {
	// TODO: convert LPCSTR correctly string to be able to read the message not just first character
	// How do we get the size of the message to know how many bytes to read?
	// fmt.Println(fmt.Sprintf("Log: msg=%s at level=%d", string(*message), level))
	return uintptr(0)
}

func (client *AzPubSubClient) InitConfig() {
	var err error
	client.hconfig, err = CallAzPubSubCreateConfiguration(client.hclient, ENUM(client.apsConfigType), AZPUBSUB_GLOBAL_CONFIGURATION_TEMPLATES_NONE)
	
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubCreateConfiguration failed with err: ", err)
	}
	if (client.hconfig == HCONFIG(0)) {
		log.Println("Failed to initialize AzPubSubClient config.")
	} else {
		log.Println("AzPubSubClient config init successfully.")
	}
}

func (client *AzPubSubClient) InitClient() {
	var err error
	client.hclient, err = CallAzPubSubClientInitialize(pLoggerCallback)

	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubClientInitialize failed with err: ", err)
	}
	if (client.hclient == HCLIENT(0)) {
		log.Println("Failed to initialize AzPubSubClient.")
	} else {
		log.Println("AzPubSubClient init successfully.")
	}
}