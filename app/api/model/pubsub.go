package model

type PublishMessage struct {
	Topic   string
	Message interface{}
}

type TriggerCommand struct {
	Job         string
	TriggerTime int64
}

type DoTaskCommand struct {
	Name        string
	TriggerTime int64
}
