package main

import (
	"errors"
	"io/ioutil"
	"os"
	"time"
)

var (
	ErrUUIDWaitTimeout = errors.New("rkt-inspect: Timeout waiting for UUID file")
)

func WaitUUIDFile(uuidFile string, timeout time.Duration) error {
	c := make(chan error, 1)
	go func() {
		for true {
			_, err := os.Stat(uuidFile)
			if err == nil {
				c <- nil
				return
			} else if !os.IsNotExist(err) {
				c <- err
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(timeout):
		return ErrUUIDWaitTimeout
	}
}

func ReadUUID(uuidFile string) (string, error) {
	dat, err := ioutil.ReadFile(uuidFile)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
