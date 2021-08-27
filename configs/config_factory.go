package configs

const (
	WORKFLOW_FILE                       = "configs/workflows.json"
	SCHEDULE_FILE                       = "configs/schedule.json"
	LOGGING_FILE                        = "log"
	MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS = 1
)

type ConfigFactory struct{}

func (factory ConfigFactory) GetConfig() Config {
	return Config{WorkflowFile: WORKFLOW_FILE, ScheduleFile: SCHEDULE_FILE, LoggingFile: LOGGING_FILE, MaxAgeToProcessWinEvtsInDays: MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS}
}
