package pods

import (
	"errors"
	"os"
	"path"
	"syscall"
	"time"
)

type PodState int

const (
	NotExist PodState = iota
	Embryo
	Prepare
	Prepared
	Run
	Exited
	ExitedGarbage
	Garbage
	PrepareFailed
	Invalid
)

var (
	ErrPodNotFound    = errors.New("rkt-inspect: Pod not found")
	ErrPodWaitTimeout = errors.New("rkt-inspect: Timeout waiting for pod")
	ErrPodExited      = errors.New("rkt-inspect: Pod failed or exited")
)

func canLock(path string) (bool, error) {
	fh, err := os.OpenFile(path, syscall.O_RDONLY|syscall.O_NOCTTY, os.FileMode(0600))
	if err != nil {
		return false, err
	}
	err = syscall.Flock(int(fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK {
		return false, nil
	} else if err == nil {
		syscall.Flock(int(fh.Fd()), syscall.LOCK_UN)
	}

	return true, err
}

func getPodDir(root string, uuid string, state string) string {
	if state == "" {
		state = "run"
	}
	return path.Join(root, "pods", state, uuid)
}

func GetPodState(root string, uuid string) (PodState, error) {
	lockable, err := canLock(getPodDir(root, uuid, "embryo"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		return Embryo, nil
	}

	lockable, err = canLock(getPodDir(root, uuid, "prepare"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		if lockable {
			return PrepareFailed, nil
		} else {
			return Prepare, nil
		}
	}

	lockable, err = canLock(getPodDir(root, uuid, "prepared"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		return Prepared, nil
	}

	lockable, err = canLock(getPodDir(root, uuid, "run"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		if lockable {
			return Exited, nil
		} else {
			return Run, nil
		}
	}

	lockable, err = canLock(getPodDir(root, uuid, "ExitedGarbage"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		return ExitedGarbage, nil
	}

	lockable, err = canLock(getPodDir(root, uuid, "Garbage"))
	if err != nil && !os.IsNotExist(err) {
		return Invalid, err
	}
	if err == nil {
		return Garbage, nil
	}

	return Invalid, nil
}

func isPodRunning(root string, uuid string) (bool, error) {
	podState, err := GetPodState(root, uuid)
	if err != nil {
		return false, err
	} else if podState > Run {
		return false, ErrPodExited
	} else if podState < Run {
		return false, nil
	} else {
		return true, nil
	}
}

func WaitPodRun(root string, uuid string, timeout time.Duration) error {
	c := make(chan error, 1)
	go func() {
		for true {
			running, err := isPodRunning(root, uuid)
			if err != nil {
				c <- err
				return
			} else if running {
				c <- nil
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(timeout):
		return ErrPodWaitTimeout
	}
}
