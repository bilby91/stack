package models

type TaskType int

const (
	TASK_FETCH_OTHERS TaskType = iota
	TASK_FETCH_ACCOUNTS
	TASK_FETCH_RECIPIENTS
	TASK_FETCH_PAYMENTS
)

type TaskTree struct {
	TaskType TaskType
	Name     string

	NextTasks []TaskTree
}

type Workflow []TaskTree
