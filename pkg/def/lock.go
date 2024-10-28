package def

import (
	"fmt"
)

func OrderLock(orderID int64) string {
	return fmt.Sprintf("lock_order_%d", orderID)
}

func CronJobLock(jobName string) string {
	return fmt.Sprintf("lock_cron_job_%s", jobName)
}
