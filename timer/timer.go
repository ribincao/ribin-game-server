package timer

import (
	"sync"
	"sync/atomic"
	"time"

	errs "github.com/ribincao/ribin-game-server/error"
)

const (
	CIRCLE_MODE     = true
	NOT_CIRCLE_MODE = false
	SYNC_MODE       = true
)

type optionCall func(*TimeWheel) error

func TickSafeMode() optionCall {
	return func(o *TimeWheel) error {
		o.tickQueue = make(chan time.Time, 10)
		return nil
	}
}

type TimeWheel struct {
	tickTime  time.Duration
	ticker    time.Ticker
	tickQueue chan time.Time

	bucketNum    int
	buckets      []map[TaskID]*Task
	bucketIndex  map[TaskID]int
	currentIndex int

	tasks    sync.Map
	taskLock sync.RWMutex

	stopC chan struct{}
	stop  bool
	sync.RWMutex
	startOnce sync.Once
	randomID  int64
}

func (tw *TimeWheel) GenUniqueId() TaskID {
	return TaskID(atomic.AddInt64(&tw.randomID, 1))
}

func NewTimeWheel(tick time.Duration, bucketNum int, options ...optionCall) (*TimeWheel, error) {
	if tick.Milliseconds() < 1 {
		return nil, errs.TimerTickError
	}
	if bucketNum <= 0 {
		return nil, errs.TimerBucketError
	}

	tw := &TimeWheel{
		tickTime:     tick,
		tickQueue:    make(chan time.Time, 10),
		bucketNum:    bucketNum,
		bucketIndex:  make(map[TaskID]int, 1024*100),
		buckets:      make([]map[TaskID]*Task, bucketNum),
		currentIndex: 0,
		stopC:        make(chan struct{}),
	}

	for i := 0; i < bucketNum; i++ {
		tw.buckets[i] = make(map[TaskID]*Task, 16)
	}

	for _, op := range options {
		op(tw)
	}
	return tw, nil
}

func (tw *TimeWheel) Start() {
	tw.startOnce.Do(
		func() {
			tw.ticker = *time.NewTicker(tw.tickTime)
			go tw.tick()
			go tw.schedule()
		},
	)
}

func (tw *TimeWheel) tick() {
	if tw.tickQueue != nil {
		return
	}
	for !tw.stop {
		<-tw.ticker.C
		select {
		case tw.tickQueue <- time.Now():
		default:
			panic("tais long time blocking")
		}
	}
}

func (tw *TimeWheel) schedule() {
	queue := tw.ticker.C
	if tw.tickQueue == nil {
		queue = tw.tickQueue
	}

	for {
		select {
		case <-queue:
			tw.handleTick()
		case <-tw.stopC:
			tw.stop = true
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) handleTick() {
	tw.Lock()
	defer tw.Unlock()

	bucket := tw.buckets[tw.currentIndex]
	for idx, task := range bucket {
		if task.stop {
			tw.deleteTask(task)
			continue
		}

		if bucket[idx].round > 0 {
			bucket[idx].round--
			continue
		}

		if task.async {
			go task.callback()
		} else {
			task.callback()
		}

		if task.circle { // Cron Task
			tw.deleteTask(task)
			tw.store(task, CIRCLE_MODE)
			continue
		}

		tw.deleteTask(task)
	}

	if tw.currentIndex == tw.bucketNum-1 {
		tw.currentIndex = 0
		return
	}
	tw.currentIndex++
}

// RemoveTask
func (tw *TimeWheel) RemoveTask(taskName string) error {
	tw.taskLock.Lock()
	defer tw.taskLock.Unlock()

	taskId, ok := tw.tasks.Load(taskName)
	if !ok {
		return nil
	}

	bucketIndex := tw.bucketIndex[taskId.(TaskID)]
	bucketTasks := tw.buckets[bucketIndex]
	task, ok := bucketTasks[taskId.(TaskID)]
	if !ok {
		tw.tasks.Delete(taskName)
		return nil
	}

	tw.tasks.Delete(taskName)
	tw.remove(task)
	return nil
}

func (tw *TimeWheel) remove(task *Task) {
	tw.Lock()
	defer tw.Unlock()
	tw.deleteTask(task)
}

func (tw *TimeWheel) deleteTask(task *Task) {
	idx := tw.bucketIndex[task.id]
	delete(tw.bucketIndex, task.id)
	delete(tw.buckets[idx], task.id)
}

// AddTask
func (tw *TimeWheel) AddTask(taskName string, delay time.Duration, callback func(), cronTask bool) error {
	if _, ok := tw.tasks.Load(taskName); ok {
		return errs.TimerTaskRepeatError
	}
	task := tw.putTask(delay, callback, cronTask, SYNC_MODE)
	if task == nil {
		tw.tasks.Delete(taskName)
		return errs.TimerTaskAddError
	}
	tw.tasks.Store(taskName, task.id)
	return nil
}

func (tw *TimeWheel) putTask(delay time.Duration, callback func(), circleMode, async bool) *Task {
	if delay <= 0 {
		delay = tw.tickTime
	}

	task := &Task{
		delay:    delay,
		callback: callback,
		circle:   circleMode,
		async:    async,
	}
	task.id = tw.GenUniqueId()

	tw.put(task, task.circle)
	return task
}

func (tw *TimeWheel) put(task *Task, circle bool) {
	tw.Lock()
	defer tw.Unlock()

	tw.store(task, circle)
}

func (tw *TimeWheel) store(task *Task, circleMode bool) {
	round, index := tw.getRoundIndex(task.delay)
	task.round = round

	if round > 0 && circleMode {
		task.round--
	}

	tw.bucketIndex[task.id] = index
	tw.buckets[index][task.id] = task
}

func (tw *TimeWheel) getRoundIndex(delay time.Duration) (int, int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tickTime.Seconds()
	return int(delaySeconds / tickSeconds / float64(tw.bucketNum)),
		(int(float64(tw.currentIndex) + delaySeconds/tickSeconds)) % tw.bucketNum
}

func (tw *TimeWheel) Stop() {
	tw.stopC <- struct{}{}
}
