package test_configs

const (
	ARBITRARY_EVT_XML                        = `<Event xmlns='http://schemas.microsoft.com/win/2004/08/events/event'><System><Provider Name='CpuSpeedMonitoring'/><EventID Qualifiers='0'>999</EventID><Version>0</Version><Level>4</Level><Task>1</Task><Opcode>0</Opcode><Keywords>0x80000000000000</Keywords><TimeCreated SystemTime='2021-08-10T19:10:29.0000000Z'/><EventRecordID>19818</EventRecordID><Correlation/><Execution ProcessID='0' ThreadID='0'/><Channel>Application</Channel><Computer>LAPTOP-0PM4COPH</Computer><Security/></System><EventData><Data>test</Data></EventData></Event>`
	TEST_WORKFLOW_FILE                       = "test_configs/test_workflows.json"
	TEST_SCHEDULE_FILE                       = "test_configs/test_schedule.json"
	TEST_SINGLE_SCHEDULE_FILE                = "test_configs/test_single_schedule.json"
	TEST_MAX_AGE_TO_PROCESS_WIN_EVTS_IN_DAYS = 1
	TEST_LOGGING_FILE                        = "log"
)
