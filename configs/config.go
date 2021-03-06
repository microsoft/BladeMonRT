package configs

const (
	WORKFLOW_FILE                       = "configs/workflows.json"
	SCHEDULE_FILE                       = "configs/schedule.json"
	LOGGING_FILE                        = "log"
	MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS = 1
	ENABLE_BOOKMARK_FEATURE      = true
	BOOKMARK_DATABASE_FILE       = "./BookmarkStore.sqlite"
	BOOKMARK_DATABASE_TABLE_NAME = "BookmarkTable"
)

type Config struct {
	WorkflowFile                 string
	ScheduleFile                 string
	LoggingFile                  string
	MaxAgeToProcessWinEvtsInDays int
	EnableBookmarkFeature        bool
	BookmarkDatabaseFile         string
	BookmarkDatabaseTableName    string
}

func NewConfig() Config {
	var config Config = Config{}
	config.WorkflowFile = WORKFLOW_FILE
	config.ScheduleFile = SCHEDULE_FILE
	config.MaxAgeToProcessWinEvtsInDays = MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS
	config.EnableBookmarkFeature = ENABLE_BOOKMARK_FEATURE
	config.BookmarkDatabaseFile = BOOKMARK_DATABASE_FILE
	config.BookmarkDatabaseTableName = BOOKMARK_DATABASE_TABLE_NAME
	return config
}
