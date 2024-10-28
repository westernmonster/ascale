package def

import (
	"ascale/pkg/conf/env"
	"fmt"
)

var TriggerJob = struct {
	CronSendLittleMessage string
	SendHugeMessage       string
}{
	CronSendLittleMessage: "CronSendLittleMessage",
	SendHugeMessage:       "SendHugeMessage",
}

var Topics = struct {
	Trigger    string
	DoTask     string
	DeadLetter string
}{
	Trigger:    fmt.Sprintf(`%s-trigger`, env.DeployEnv),
	DoTask:     fmt.Sprintf(`%s-do-task`, env.DeployEnv),
	DeadLetter: fmt.Sprintf(`%s-deadletter`, env.DeployEnv),
}
