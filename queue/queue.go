package queue

type Queue struct {
	Name      string
	State     string
	Consumers int
	Messages  int
}

type QueueWithType struct {
	Queue
	QueueType string
	Stage     string
	SessionId string
}
