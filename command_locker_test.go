package main

import (
	"testing"
	"time"

	"github.com/viafintech/cronlocker/testutils"
)

func TestConsulCommandLockerLockAndExecute(t *testing.T) {
	cases := []struct {
		title                string
		runConcurrent        func(locker *ConsulCommandLocker)
		key                  string
		command              string
		expectedOutputString string
		expectedErrorString  string
	}{
		{
			title:                "success",
			key:                  "test/cron/service/job_name1",
			command:              "echo 1",
			expectedOutputString: "1\n",
		},
		{
			title: "cannot aquire lock",
			runConcurrent: func(locker *ConsulCommandLocker) {
				locker.LockAndExecute("test/cron/service/job_name2", "sleep 2")
			},
			key:                  "test/cron/service/job_name2",
			command:              "echo 1",
			expectedOutputString: "Nothing was executed\n",
		},
		{
			title:               "command fails",
			key:                 "test/cron/service/job_name3",
			command:             "false",
			expectedErrorString: "exit status 1",
		},
	}

	commandLocker, _ := NewConsulCommandLocker(
		testutils.CONSULURI,
		300*time.Millisecond,
		time.Millisecond,
		0,
	)

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			if c.runConcurrent != nil {
				go c.runConcurrent(commandLocker)
				time.Sleep(500 * time.Millisecond)
			}

			outputString, err := commandLocker.LockAndExecute(
				c.key,
				c.command,
			)

			if outputString != c.expectedOutputString {
				t.Errorf(
					"Did not received expected output string:\n%s\n\nReceived:\n%s",
					c.expectedOutputString,
					outputString,
				)
			}

			if err != nil && err.Error() != c.expectedErrorString {
				t.Errorf(
					"Did not received expected error:%#v, Received: %#v",
					c.expectedErrorString,
					err,
				)
			}
		})
	}
}

func TestConsulCommandLockerMinimumLockAndExecuteTime(t *testing.T) {
	commandLocker, _ := NewConsulCommandLocker(
		testutils.CONSULURI,
		300*time.Millisecond,
		500*time.Millisecond,
		0,
	)

	startTime := time.Now()

	commandLocker.LockAndExecute("test/cron/service/min_time_job", "echo 1")

	if time.Since(startTime) <= 500*time.Millisecond {
		t.Errorf("Locker did not wait the minimum time the lock should have been held")
	}
}

func TestConsulCommandLockerMaximumExecutionTime(t *testing.T) {
	commandLocker, _ := NewConsulCommandLocker(
		testutils.CONSULURI,
		100*time.Millisecond,
		300*time.Millisecond,
		500*time.Millisecond,
	)

	startTime := time.Now()

	commandLocker.LockAndExecute("test/cron/service/min_time_job", "sleep 5")

	if time.Since(startTime) > 700*time.Millisecond {
		t.Errorf("Locker did not abort after the maximum execution time was reached")
	}
}
