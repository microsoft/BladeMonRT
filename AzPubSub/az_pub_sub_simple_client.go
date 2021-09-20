package azpubsub

import "C"
import (
	"syscall"
	"log"
	"unsafe"
	"golang.org/x/sys/windows"
	"errors"
)

type AzPubSubSimpleClient struct {
	AzPubSubClient
}

func NewAzPubSubSimpleClient(isTestInstance bool, endpoint string) *AzPubSubSimpleClient {
	var client AzPubSubSimpleClient = AzPubSubSimpleClient{AzPubSubClient: NewAzPubSubClient(isTestInstance, endpoint)}
	client.apsConfigType = AZPUBSUB_CONFIGURATION_TYPE_SIMPLE
	client.InitClient()
	client.InitConfig()

	client.OpenSimpleProducer()

	return &client
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

func CallAzPubSubOpenSimpleProducer(config HCONFIG, apsSecurityType AZPUBSUB_SECURITY_TYPE, endpoint string) (HPRODUCER, error) {
	endpointPtr, err := syscall.UTF16PtrFromString(endpoint)
	if (err != nil) {
		return HPRODUCER(0), err
	}

	hproducer, _, err := AzPubSubOpenSimpleProducer.Call(uintptr(config), uintptr(apsSecurityType), 0, 0, uintptr(unsafe.Pointer(endpointPtr)))
	return HPRODUCER(hproducer), err 
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