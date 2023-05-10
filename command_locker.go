package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"time"

	"github.com/hashicorp/consul/api"
)

var errLockWaitTimeExpired = errors.New("wait time for acquiring the lock expired")

// ConsulCommandLocker is an implementation of a command locker for consul,
// responsible for acquiring a distributed lock and executing a command
type ConsulCommandLocker struct {
	apiClient                *api.Client
	lockWaitTime             time.Duration
	minLockTime              time.Duration
	maxExecTime              time.Duration
	failOnLockWaitExpiration bool
}

// NewConsulCommandLocker initializes a new ConsulCommandlocker
func NewConsulCommandLocker(
	endpoint string,
	token string,
	lockWaitTime time.Duration,
	minLockTime time.Duration,
	maxExecTime time.Duration,
	failOnLockWaitExpiration bool,
) (*ConsulCommandLocker, error) {
	ccl := &ConsulCommandLocker{
		lockWaitTime:             lockWaitTime,
		minLockTime:              minLockTime,
		maxExecTime:              maxExecTime,
		failOnLockWaitExpiration: failOnLockWaitExpiration,
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	address := url.Hostname() + ":" + url.Port()

	config := &api.Config{
		Address:  address,
		Scheme:   url.Scheme,
		WaitTime: time.Second,
		Token:    token,
	}

	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	ccl.apiClient = apiClient

	return ccl, nil
}

// LockAndExecute takes a lock key and executes the command
func (ccl *ConsulCommandLocker) LockAndExecute(key, command string) (string, error) {
	lockOpts := &api.LockOptions{
		Key:          key,
		LockWaitTime: ccl.lockWaitTime,
		LockTryOnce:  true,
	}

	lock, err := ccl.apiClient.LockOpts(lockOpts)
	if err != nil {
		return "", err
	}
	// Unlock can return an error but it will be unlocked anyways if the connection is lost
	// so we only want to make sure here, that we can return early
	defer lock.Unlock()

	lockCh, err := lock.Lock(nil)
	if err != nil {
		return "", err
	}
	// The lock was not acquired if lock channel is empty
	// Therefore we can simply return
	if lockCh == nil {
		if ccl.failOnLockWaitExpiration {
			return "", errLockWaitTimeExpired
		}

		return "Nothing was executed\n", nil
	}

	ctx := context.Background()
	var cancel func()
	if ccl.maxExecTime != 0 {
		ctx, cancel = context.WithTimeout(ctx, ccl.maxExecTime)
		defer cancel()
	}

	targetTime := time.Now().Add(ccl.minLockTime)

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	// Ensure to wait at least the minimum lock time
	if remainingTime := targetTime.Sub(time.Now()); remainingTime > 0 {
		time.Sleep(remainingTime)
	}

	resultOutput := stdout.String()
	if stderr.String() != "" {
		resultOutput = fmt.Sprintf("%s\nstderr: %s", resultOutput, stderr.String())
	}

	if err != nil {
		return resultOutput, err
	}

	return resultOutput, nil
}
