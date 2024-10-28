package model

type PublishMessage struct {
	Topic   string
	Message interface{}
}

type CronJobCommand struct {
	Job         string
	TriggerTime int64
}

type DoSmallTaskCommand struct {
	Name        string
	TriggerTime int64
}
