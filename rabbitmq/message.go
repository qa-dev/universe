package rabbitmq

// Информация о задаче запуска.
type TaskInfo struct {
	ID        int    `json:"id"`
	TaskGroup string `json:"task_group"`
}

// Сообщение для выполнения.
type Message struct {
	TaskInfo  *TaskInfo `json:"task"`
	Cmd       []string  `json:"cmd"`
	Env       string    `json:"env"`
	TryNumber int       `json:"try_number"`
}
