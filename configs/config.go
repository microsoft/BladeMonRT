package configs

type Config struct {
	WorkflowFile                 string
	ScheduleFile                 string
	LoggingFile                  string
	MaxAgeToProcessWinEvtsInDays int
	EnableBookmarkFeature        bool
	BookmarkDatabaseFile         string
	BookmarkDatabaseTableName    string
}
