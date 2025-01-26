package constants

type Condition int

const (
	RUNNING Condition = iota
	NOT_RUNNING
	EXCEPTION
	PASS
)

type TaskOperation int

const (
	TASK_START TaskOperation = iota
	TASK_STOP
	TASK_START_WAIT_DONE
	TASK_STOP_WAIT_DONE
)
