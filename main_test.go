package main

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/test_configs"
	"io/ioutil"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	// Assume
	workflowsJson, err := ioutil.ReadFile(test_configs.TEST_WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	// Create the schedules JSON from a file with a single schedule.
	schedulesJson, err := ioutil.ReadFile(test_configs.TEST_SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWorkflowScheduler := NewMockWorkflowSchedulerInterface(ctrl)

	// Assert
	// Set up assertion.
	var cpuMonitoringExpectedEventQuery WinEventSubscribeQuery = WinEventSubscribeQuery{channel: "Application", query: "*[System[Provider[@Name='CpuSpeedMonitoring']]]"}
	cpuMonitoringScheduleExpectedEventQueries := []WinEventSubscribeQuery{cpuMonitoringExpectedEventQuery}
	mockWorkflowScheduler.EXPECT().addWinEventBasedSchedule(gomock.Any(), cpuMonitoringScheduleExpectedEventQueries)
	var disk7ExpectedEventQuery WinEventSubscribeQuery = WinEventSubscribeQuery{channel: "System", query: "*[System[Provider[@Name='disk'] and EventID=7]]"}
	var disk8ExpectedEventQuery WinEventSubscribeQuery = WinEventSubscribeQuery{channel: "System", query: "*[System[Provider[@Name='disk'] and EventID=8]]"}
	diskScheduleExpectedEventQueries := []WinEventSubscribeQuery{disk7ExpectedEventQuery, disk8ExpectedEventQuery}
	mockWorkflowScheduler.EXPECT().addWinEventBasedSchedule(gomock.Any(), diskScheduleExpectedEventQueries)

	// Assume
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Main")
	var mainObj *Main = &Main{WorkflowScheduler: mockWorkflowScheduler, logger: logger}

	// Action
	mainObj.setupWorkflows(schedulesJson, workflowFactory)
}
