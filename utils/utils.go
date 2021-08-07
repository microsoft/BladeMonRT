package utils

import (
	"regexp"
	"time"
	"strconv"
  "log"
  "github.com/microsoft/BladeMonRT/logging"
)

/** Class that contains utilities used in BladeMonRT classes. */
type Utils struct {
	logger *log.Logger
}

func NewUtils() *Utils {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Utils")
	return &Utils{logger : logger}
}

/** Parses out the Event Provider, EventID, TimeCreated[SystemTime], eventRecordID (which isdifferent from event ID) from some Event XML. Returns these values as an str, int, date, and int respectively. */
func (utils *Utils) ParseEventXML(eventXML string) (string, int, time.Time, int) {
  re := regexp.MustCompile(`.*Provider *Name=[\"\']([^\"]+)[\"\'].*<EventID[^>]*>([0-9]+)</EventID>.*<TimeCreated +SystemTime=[\"\']([0-9\-]*)T.*<EventRecordID>([0-9]+)</EventRecordID>.*`)
  attributes := re.FindStringSubmatch(eventXML)
  var provider string = attributes[1]

  eventID, err := strconv.Atoi(attributes[2])
  if (err != nil) {
	  utils.logger.Println("Wrong format of event ID.")
  }

  timeCreated, err := time.Parse("2006-01-02", attributes[3])
  if (err != nil) {
	  utils.logger.Println("Wrong format of time.")
  }

  eventRecordID, err := strconv.Atoi(attributes[4])
  if (err != nil) {
	  utils.logger.Println("Wrong format of event record ID.")
  }

  return provider, eventID, timeCreated, eventRecordID
}