package timer

import (
	"fmt"
	"sort"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
)

type Task struct {
	Callback func()
	delay    time.Duration
	dueTime  time.Time
}

type TimeWheel struct {
	wheelSize int                     // the number of slots in the time wheel
	baseTick  time.Duration           // the time duration represented by each slot
	slots     []*gqueue.TQueue[*Task] // slot[i] contains (i, i+1]*baseTick delayed tasks
	cursor    int                     // the current position of the time wheel
	maxDelay  time.Duration           // the maximum delay that can be handled by this time wheel
}

func NewTimeWheel(wheelSize int, baseTick time.Duration) *TimeWheel {
	timeWheel := &TimeWheel{
		wheelSize: wheelSize,
		baseTick:  baseTick,
		slots:     make([]*gqueue.TQueue[*Task], wheelSize),
		cursor:    0,
		maxDelay:  time.Duration(wheelSize) * baseTick,
	}
	for i := 0; i < wheelSize; i++ {
		timeWheel.slots[i] = gqueue.NewTQueue[*Task]()
	}
	return timeWheel
}

type Timer struct {
	timeWheels []*TimeWheel
}

// create a new timer, if minTick is less than 1 second, the timer will create a time wheel with a base tick of minTick, and the number of slots is 1 second / minTick, then the timer will create time wheels with a base tick of 1 second, 1 minute, 1 hour, etc. until the maxDelay is covered
func NewTimer(minTick time.Duration, maxDelay time.Duration) (*Timer, error) {
	timer := &Timer{
		timeWheels: []*TimeWheel{},
	}

	if minTick <= 0 {
		return nil, fmt.Errorf("minTick must be > 0")
	}
	if maxDelay < minTick {
		return nil, fmt.Errorf("maxDelay must be >= minTick")
	}

	levels := []int{60, 60, 24}
	if minTick < time.Second {
		if time.Second%minTick != 0 {
			return nil, fmt.Errorf("minTick must be a divisor of 1 second")
		} else {
			levels = append([]int{int(time.Second / minTick)}, levels...)
		}
	}

	for _, wheelSize := range levels {
		timer.timeWheels = append(timer.timeWheels, NewTimeWheel(wheelSize, minTick))
		if minTick*time.Duration(wheelSize) >= maxDelay {
			return timer, nil
		}
		minTick *= time.Duration(wheelSize)
	}

	for minTick < maxDelay {
		const overflowWheelSize = 128
		timer.timeWheels = append(timer.timeWheels, NewTimeWheel(overflowWheelSize, minTick))
		minTick *= overflowWheelSize
	}

	return timer, nil
}

// AddTask accepts delays in (0, maxDelay].
func (t *Timer) AddTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("task must not be nil")
	}
	if task.delay <= 0 {
		return fmt.Errorf("task delay must be > 0")
	}
	task.dueTime = time.Now().Add(task.delay)

	for _, timeWheel := range t.timeWheels {
		if task.delay <= timeWheel.maxDelay {
			ticks := int((task.delay+timeWheel.baseTick-1)/timeWheel.baseTick) - 1
			slotIndex := (timeWheel.cursor + ticks) % timeWheel.wheelSize
			timeWheel.slots[slotIndex].Push(task)
			return nil
		}
	}
	return fmt.Errorf("task delay is too large")
}

func PopAll(q *gqueue.TQueue[*Task]) []*Task {
	var tasks []*Task
	for q.Len() > 0 {
		tasks = append(tasks, q.Pop())
	}
	return tasks
}

// Tick advances the timer by one tick (the smallest base tick of the time wheels) and executes all tasks that are due at the current time. It also handles the cascading of tasks from higher-level time wheels to lower-level ones as time advances.
func (t *Timer) Tick() {
	var tasks []*Task
	for i, timeWheel := range t.timeWheels {
		if i == 0 {
			tasks = PopAll(timeWheel.slots[timeWheel.cursor])
			timeWheel.cursor = (timeWheel.cursor + 1) % timeWheel.wheelSize
		} else {
			lastTimeWheel := t.timeWheels[i-1]
			timeWheel.cursor = (timeWheel.cursor + 1) % timeWheel.wheelSize

			for timeWheel.slots[timeWheel.cursor].Len() > 0 { // cascading tasks to lower-level time wheel
				task := timeWheel.slots[timeWheel.cursor].Pop()
				if task.delay == 0 {
					task.delay = timeWheel.baseTick
				}
				ticks := int((task.delay+lastTimeWheel.baseTick-1)/lastTimeWheel.baseTick) - 1
				slotIndex := (lastTimeWheel.cursor + ticks) % lastTimeWheel.wheelSize
				lastTimeWheel.slots[slotIndex].Push(task)

			}
		}
		if timeWheel.cursor != 0 {
			break
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].dueTime.Before(tasks[j].dueTime)
	})
	for _, task := range tasks {
		task.Callback()
	}
}

func (t *Timer) Run(done chan struct{}) {
	tick := t.timeWheels[0].baseTick
	nextTick := time.Now()
	for {
		select {
		case <-done:
			return
		default:
			nextTick = nextTick.Add(tick)
			time.Sleep(time.Until(nextTick))
			t.Tick()
		}
	}
}
