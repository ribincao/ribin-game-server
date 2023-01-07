package timer

import "time"

type TaskID int64
type Task struct {
	delay    time.Duration
	id       TaskID
	round    int
	callback func()

	async  bool
	stop   bool
	circle bool
}
