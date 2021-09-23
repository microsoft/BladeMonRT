package azpubsub

import (
	"C"
	"unsafe"
	"log"
)

type AZPUBSUB_TOPIC_PARTITIONER func(arg1 DWORD, arg2 HPRODUCERTOPIC, arg3 PBYTE, arg4 DWORD, arg5 DWORD, arg6 LPVOID)

type AzPubSubGlobalClient struct {
    AzPubSubClient
    openedTopics map[string]HPRODUCERTOPIC
}
func NewAzPubSubGlobalClient(topics []string, isTestInstance bool, endpoint string) *AzPubSubGlobalClient {
    var client AzPubSubGlobalClient = AzPubSubGlobalClient{AzPubSubClient: NewAzPubSubClient(isTestInstance, endpoint)}
    client.apsConfigType = AZPUBSUB_CONFIGURATION_TYPE_GLOBAL
    client.InitClient()
    client.InitConfig()

	client.AddConfigKey("azpubsub.security.provider", "ApPki")
	// client.SetConnection()
	// open global producer.

    for index := range topics {
        client.OpenTopic(topics[index])
    }
    return &client
}

/**
func (client *AzPubSubGlobalClient) SetConnection() {
    err := CallAzPubSubSetConnection(client.hconfig, client.endpoint, client.apsSecurityType, client.apsConnectionFlags)
    if (err.Error() != ERR_OK) {
        log.Println("AzPubSubSetConnection failed with err: ", err)
    }
}
*/

func (client *AzPubSubGlobalClient) OpenTopic(topic string) {
    hproducerTopic, err := CallAzPubSubOpenProducerTopic(client.Hproducer, topic)
    if (err.Error() != ERR_OK) {
        log.Println("AzPubSubOpenProducerTopic failed with err: ", err)
    }
    if (HPRODUCERTOPIC(hproducerTopic) == HPRODUCERTOPIC(0)) {
        log.Println("AzPubSubOpenSimpleProducer failed with status.")
    } else {
        log.Println("AzPubSubOpenSimpleProducer init successfully")
    }
}

func (client *AzPubSubGlobalClient) AddConfigKey(key string, value string) {
    log.Println("Adding configuration key", key)
    err := CallAzPubSubAddStringConfiguration(client.Hconfig, key, value)
    if (err.Error() != ERR_OK) {
        log.Println("AzPubSubAddStringConfiguration failed with err: ", err)
    }
}

func CallAzPubSubAddStringConfiguration(config HCONFIG, key string, value string) error {
    keyCString := C.CString(key)
    valueCString := C.CString(value)
    _, _, err := AzPubSubAddStringConfiguration.Call(uintptr(config), uintptr(unsafe.Pointer(keyCString)), uintptr(unsafe.Pointer(valueCString)))
    return err
}

func CallAzPubSubOpenProducerTopic(hproducer HPRODUCER, topic string) (HPRODUCERTOPIC, error) {
	topicCString := C.CString(topic)
    hproducerTopic, _, err := AzPubSubOpenProducerTopic.Call(uintptr(hproducer),
    uintptr(unsafe.Pointer(topicCString)), uintptr(HCONFIG(0)), uintptr(0)) // How to represent null value of AZPUBSUB_TOPIC_PARTITIONER?
    return HPRODUCERTOPIC(hproducerTopic), err
}

/*
func CallAzPubSubSetConnection(hconfig HCONFIG, endpoint string, apsSecurityType AZPUBSUB_SECURITY_TYPE, apsConnectionFlags AZPUBSUB_CONNECTION_FLAGS) error {
	topicCString := C.CString(topic)
    hproducerTopic, _, err := CallAzPubSubSetConnection.Call(uintptr(hconfig), uintptr(0), uintptr(0), uintptr(C.Cstring(endpoint))), uintptr(0), uintptr(apsSecurityType), uintptr(apsConnectionFlags))
    return err
}
**/