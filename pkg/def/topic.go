package def

import (
	"ascale/pkg/conf/env"
	"fmt"
)

var CronJob = struct {
	DoSmallTask string
}{
	DoSmallTask: "DoSmallTask",
}

var Topics = struct {
	CronJob     string
	DoSmallTask string
	DeadLetter  string
}{
	CronJob:     fmt.Sprintf(`%s-cron-job`, env.DeployEnv),
	DoSmallTask: fmt.Sprintf(`%s-do-small-task-send`, env.DeployEnv),
	DeadLetter:  fmt.Sprintf(`%s-deadletter`, env.DeployEnv),
}
