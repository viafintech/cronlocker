package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	lockWaitTimeMS = flag.Int(
		"lockwaittime",
		500,
		"Configures the wait time for a lock in milliseconds",
	)
	minLockTimeMS = flag.Int(
		"minlocktime",
		5000,
		"Configures the minimum time in milliseconds a lock is held",
	)
	endpoint = flag.String("endpoint", "http://localhost:8500", "endpoint")
	key      = flag.String("key", "none", "key to monitor, e.g. cronjobs/any_service/cron_name")
)

func checkFlag(f *string, name string) {
	if *f == "none" || *f == "" {
		log.Fatalf("Setting %s is mandatory", name)
	}
}

func main() {
	flag.Parse()

	// Ensure flags are given
	checkFlag(endpoint, "endpoint")
	checkFlag(key, "key")

	// original command
	command := strings.Join(flag.Args(), " ")

	// Initiate command locker
	ccl, err := NewConsulCommandLocker(
		*endpoint,
		time.Duration(*lockWaitTimeMS)*time.Millisecond,
		time.Duration(*minLockTimeMS)*time.Millisecond,
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	output, err := ccl.LockAndExecute(*key, command)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Command output:")
	fmt.Print(output)
}
