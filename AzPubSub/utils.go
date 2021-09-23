package azpubsub

import (
"C"
"fmt"
"regexp"
"log"
"io/ioutil"
"github.com/microsoft/BladeMonRT/logging"
)

/** Class that contains utilities to communicate with AzPubSub. */
// TODO: Move the utils into a wrapper.
type Utils struct {
	logger *log.Logger
}

func NewUtils() *Utils {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Utils")
	return &Utils{logger: logger}
}

func (utils *Utils) FetchAzPubSubPfVIP(pfAzPubSubVipFile string) (string, error) {
	utils.logger.Println("Fetching the AzPubSubPfVIP.")

	contentBytes, err := ioutil.ReadFile(pfAzPubSubVipFile)
    if err != nil {
        utils.logger.Println(fmt.Sprintf("Failed to read %s", pfAzPubSubVipFile))
		return "", err
    }
	var content string = string(contentBytes)
  
	regexExpression, err := regexp.Compile(`.*KafkaClusterEndpoint=([0-9\.]+).*`)
	if err != nil {
        utils.logger.Println(fmt.Sprintf("Failed to create regex object to find KafkaClusterEndpoint."))
		return "", err
    }
	regexMatches := regexExpression.FindStringSubmatch(content)
    var kafkaClusterEndpoint string = regexMatches[1]
   
	return kafkaClusterEndpoint, nil
}