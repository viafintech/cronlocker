package main

import (
	"bytes"
	"context"
	"net/url"
	"os/exec"
	"time"

	"github.com/hashicorp/consul/api"
)

type ConsulCommandLocker struct {
	apiClient    *api.Client
	lockWaitTime time.Duration
	minLockTime  time.Duration
	maxExecTime  time.Duration
}

func NewConsulCommandLocker(
	endpoint string,
	lockWaitTime time.Duration,
	minLockTime time.Duration,
	maxExecTime time.Duration,
) (*ConsulCommandLocker, error) {
	ccl := &ConsulCommandLocker{
		lockWaitTime: lockWaitTime,
		minLockTime:  minLockTime,
		maxExecTime:  maxExecTime,
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
	}

	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	ccl.apiClient = apiClient

	return ccl, nil
}

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
		return "Nothing was executed\n", nil
	}

	ctx := context.Background()
	var cancel func()
	if ccl.maxExecTime == 0 {
		ctx, cancel = context.WithTimeout(ctx, ccl.maxExecTime)
		defer cancel()
	}

	targetTime := time.Now().Add(ccl.minLockTime)

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

	// Ensure to wait at least the minimum lock time
	if remainingTime := targetTime.Sub(time.Now()); remainingTime > 0 {
		time.Sleep(remainingTime)
	}

	if err != nil {
		return "", err
	}

	return out.String(), nil
}
