package azpubsub

import "C"
import (
	"syscall"
	"log"
	"unsafe"
	//"sysc"
	"golang.org/x/sys/windows"
	"errors"
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

// ======================================= AzPubSubGlobalClient

type AzPubSubGlobalClient struct {
	AzPubSubClient
	openedTopics map[string]HPRODUCERTOPIC
}

func NewAzPubSubGlobalClient(topics []string, isTestInstance bool, endpoint string) *AzPubSubGlobalClient {
	var client AzPubSubGlobalClient = AzPubSubGlobalClient{AzPubSubClient: NewAzPubSubClient(isTestInstance, endpoint)}
	client.apsConfigType = AZPUBSUB_CONFIGURATION_TYPE_GLOBAL
	client.InitClient()
	client.InitConfig()

	for index := range topics {
		client.OpenTopic(topics[index])
	}

	return &client
}

func (client *AzPubSubGlobalClient) OpenTopic(topic string) {
	hproducerTopic, err := CallAzPubSubOpenProducerTopic(client.hproducer, topic)
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubOpenProducerTopic failed with err: ", err)
	}
	if (HPRODUCERTOPIC(hproducerTopic) == HPRODUCERTOPIC(0)) {
		log.Println("AzPubSubOpenSimpleProducer failed with status.")
	} else {
		log.Println("AzPubSubOpenSimpleProducer init successfully")
	}
}

// ======================================= AzPubSubSimpleClient

type AzPubSubSimpleClient struct {
	AzPubSubClient
}

func NewAzPubSubSimpleClient(isTestInstance bool, endpoint string) *AzPubSubSimpleClient {
	var client AzPubSubSimpleClient = AzPubSubSimpleClient{AzPubSubClient: NewAzPubSubClient(isTestInstance, endpoint)}
	client.apsConfigType = AZPUBSUB_CONFIGURATION_TYPE_SIMPLE
	client.InitClient()
	client.InitConfig()
	client.AddConfigKey("azpubsub.security.provider", "ApPki")

	client.OpenSimpleProducer()

	return &client
}

func (client *AzPubSubSimpleClient) AddConfigKey(key string, value string) {
	log.Println("Adding configuration key", key)

	err := CallAzPubSubAddStringConfiguration(client.hconfig, key, value)
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubAddStringConfiguration failed with err: ", err)
	}
}

func (client *AzPubSubSimpleClient) OpenSimpleProducer() {
	var err error
	client.hproducer, err = CallAzPubSubOpenSimpleProducer(client.hconfig, client.apsSecurityType, client.endpoint)
	
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubOpenSimpleProducer failed with err: ", err)
	}
	if (client.hproducer == HPRODUCER(0)) {
		log.Println("AzPubSubOpenSimpleProducer failed with status.")
	} else {
		log.Println("AzPubSubOpenSimpleProducer init successfully")
	}
}


func (client *AzPubSubSimpleClient) SendMessage(topic string, msg string) (SimpleResponse, error) {
	response, status, err := CallAzPubSubSendMessageEx(client.hproducer, topic, msg)
	if (err.Error() != ERR_OK) {
		log.Println("CallAzPubSubSendMessageEx failed with err: ", err)
		return SimpleResponse{}, errors.New("Send message failed.") 
	}
	if (status != STATUS_OK) {
		log.Println("CallAzPubSubSendMessageEx failed with status: ", status)
		return SimpleResponse{}, errors.New("Send message failed.")
	}

	simpleResponse, err := newSimpleResponse(response)
	if (err != nil) {
		log.Println("Send message failed to parse response.")
		return simpleResponse, errors.New("Send message failed.")
	} 
	
	log.Println("Send message parsed response with message, status code, and substatus code:", simpleResponse.message, simpleResponse.statusCode, simpleResponse.subStatusCode)
	return simpleResponse, nil
}

func CallAzPubSubSendMessageEx(hproducer HPRODUCER, topic string, msg string) (HRESPONSE, DWORD, error) {
	msgPtr, err := windows.UTF16PtrFromString(msg)
	if (err != nil) {
		return HRESPONSE(0), DWORD(0), errors.New("Failed to convert message to correct format.")
	}
	topicPtr, err := windows.UTF16PtrFromString(topic)
	if (err != nil) {
		return HRESPONSE(0), DWORD(0), errors.New("Failed to convert topic to correct format.") 
	}
	var msgLength int = len(msg)
	var hresponse HRESPONSE = HRESPONSE(0)

	// Assume default hash based partitioning
	// TODO: Add key parameter so that we can do other types of partitioning.
	status, _, err := AzPubSubSendMessageEx.Call(uintptr(hproducer),
	uintptr(unsafe.Pointer(topicPtr)),
	0,
	0,
	uintptr(unsafe.Pointer(msgPtr)),
	uintptr(msgLength),
	uintptr(unsafe.Pointer(&hresponse)))

	return hresponse, DWORD(status), err
}

func CallAzPubSubClientInitialize(callback AZPUBSUB_LOG_CALLBACK) (HCLIENT, error) {
	hclient, _, err := AzPubSubClientInitialize.Call(syscall.NewCallback(pLoggerCallback), uintptr(LPVOID(0)))
	return HCLIENT(hclient), err
}

func CallAzPubSubCreateConfiguration(client HCLIENT, configType ENUM, globalConfigTemplate UINT) (HCONFIG, error) {
	hconfig, _, err := AzPubSubCreateConfiguration.Call(uintptr(client), uintptr(configType), uintptr(globalConfigTemplate))
	return HCONFIG(hconfig), err
}

func CallAzPubSubAddStringConfiguration(config HCONFIG, key string, value string) error {
	keyCString := C.CString(key)
	valueCString := C.CString(value)
	_, _, err := AzPubSubAddStringConfiguration.Call(uintptr(config), uintptr(unsafe.Pointer(keyCString)), uintptr(unsafe.Pointer(valueCString)))
	return err
}

func CallAzPubSubOpenSimpleProducer(config HCONFIG, apsSecurityType AZPUBSUB_SECURITY_TYPE, endpoint string) (HPRODUCER, error) {
	endpointPtr, err := syscall.UTF16PtrFromString(endpoint)
	if (err != nil) {
		return HPRODUCER(0), err
	}

	hproducer, _, err := AzPubSubOpenSimpleProducer.Call(uintptr(config), uintptr(apsSecurityType), 0, 0, uintptr(unsafe.Pointer(endpointPtr)))
	return HPRODUCER(hproducer), err 
}

type SimpleResponse struct {
	hResponse HRESPONSE
	statusCode INT 
	message string 
	subStatusCode INT 
}

func newSimpleResponse(hresponse HRESPONSE) (SimpleResponse, error) {
	message, err := getResponseMessage(hresponse)
	if (err != nil) {
		return SimpleResponse{}, nil
	}

	statusCode, err := getStatusCode(hresponse)
	if (err != nil) {
		return SimpleResponse{}, nil
	}

	subStatusCode, err := getSubStatusCode(hresponse)
	if (err != nil) {
		return SimpleResponse{}, nil
	}

	return SimpleResponse{message: message, statusCode: statusCode, subStatusCode: subStatusCode}, nil
}

func getResponseMessage(hresponse HRESPONSE) (string, error) {
	bufferLength := DWORD(0)
	responseSpace := make([]byte, 0)
	status, err := CallAzPubSubResponseGetMessage(hresponse, &responseSpace,
		bufferLength,
		&bufferLength)

	if (err.Error() != ERR_OK) {
		log.Println("First AzPubSubResponseGetMessage failed with error:", err)
		return "", errors.New("AzPubSubResponseGetMessage failed.")
	}
	if (status != STATUS_MORE_DATA) {
		log.Println("First AzPubSubResponseGetMessage failed with status:", status)
		return "", errors.New("AzPubSubResponseGetMessage failed.")
	}
	
	responseSpace = make([]byte, bufferLength)
	status, err = CallAzPubSubResponseGetMessage(hresponse, &responseSpace,
		bufferLength,
		&bufferLength)
	
	if (err.Error() != ERR_OK) {
		log.Println("Second AzPubSubResponseGetMessage failed with error:", err)
		return "", errors.New("AzPubSubResponseGetMessage failed.")
	}
	if (status != STATUS_OK) {
		log.Println("Second AzPubSubResponseGetMessage failed with status:", status)
		return "", errors.New("AzPubSubResponseGetMessage failed.")
	}

	return string(responseSpace), nil

}

func getStatusCode(hresponse HRESPONSE) (INT, error) {
	statusCode := INT(-1)
	status, err := CallAzPubSubResponseGetStatusCode(hresponse, &statusCode)
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubResponseGetStatusCode failed with error:", err)
		return INT(-1), errors.New("AzPubSubResponseGetStatusCode failed.")
	}
	if status != STATUS_OK {
		log.Println("AzPubSubResponseGetStatusCode with status:", status)
		return INT(-1), errors.New("AzPubSubResponseGetStatusCode failed.")
	}

	return statusCode, nil
}

func getSubStatusCode(hresponse HRESPONSE) (INT, error) {
	statusCode := INT(-1)
	status, err := CallAzPubSubResponseGetSubStatusCode(hresponse, &statusCode)
	if (err.Error() != ERR_OK) {
		log.Println("AzPubSubResponseGetSubStatusCode failed with error:", err)
		return INT(-1), errors.New("AzPubSubResponseGetSubStatusCode failed.")
	}
	if status != STATUS_OK {
		log.Println("AzPubSubResponseGetSubStatusCode with status:", status)
		return INT(-1), errors.New("AzPubSubResponseGetSubStatusCode failed.")
	}

	return statusCode, nil
}

func CallAzPubSubResponseGetMessage(hresponse HRESPONSE, responseBuff LPPSTR, bufferLength DWORD, pointerToBufferLength PDWORD) (DWORD, error) {
	status, _, err := AzPubSubResponseGetMessage.Call(uintptr(hresponse),
	uintptr(unsafe.Pointer(responseBuff)),
	uintptr(bufferLength),
	uintptr(unsafe.Pointer(pointerToBufferLength)))
	return DWORD(status), err
}

func CallAzPubSubResponseGetStatusCode(hresponse HRESPONSE, statusCode PINT) (DWORD, error) {
	status, _, err := AzPubSubResponseGetStatusCode.Call(uintptr(hresponse),
	uintptr(unsafe.Pointer(statusCode)))
	return DWORD(status), err
}

func CallAzPubSubResponseGetSubStatusCode(hresponse HRESPONSE, subStatusCode PINT) (DWORD, error) {
	subStatus, _, err := AzPubSubResponseGetSubStatusCode.Call(uintptr(hresponse),
	uintptr(unsafe.Pointer(subStatusCode)))
	return DWORD(subStatus), err
}

// Global client wrapper
func CallAzPubSubOpenProducerTopic(hproducer HPRODUCER, topic string) (HPRODUCERTOPIC, error) {
	topicString, err := syscall.UTF16PtrFromString(topic)
	if (err != nil) {
		return HPRODUCERTOPIC(0), err
	}
	hproducerTopic, _, err := AzPubSubOpenProducerTopic.Call(uintptr(hproducer),
	uintptr(unsafe.Pointer(topicString)), uintptr(HCONFIG(0)), uintptr(0)) // How to represent null value of AZPUBSUB_TOPIC_PARTITIONER?
	return HPRODUCERTOPIC(hproducerTopic), err
}