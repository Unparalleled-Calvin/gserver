package timer

import (
	"fmt"
	"time"
)

type Task struct {
	Callback func()
	delay    time.Duration
}

type TimeWheel struct {
	wheelSize int           // the number of slots in the time wheel
	baseTick  time.Duration // the time duration represented by each slot
	slots     [][]*Task     // slot[i] contains (i, i+1]*baseTick delayed tasks
	cursor    int           // the current position of the time wheel
	maxDelay  time.Duration // the maximum delay that can be handled by this time wheel
}

func NewTimeWheel(wheelSize int, baseTick time.Duration) *TimeWheel {
	return &TimeWheel{
		wheelSize: wheelSize,
		baseTick:  baseTick,
		slots:     make([][]*Task, wheelSize),
		cursor:    0,
		maxDelay:  time.Duration(wheelSize) * baseTick,
	}
}

type Timer struct {
	timeWheels  []*TimeWheel
	currentTime time.Time
}

// create a new timer, if minTick is less than 1 second, the timer will create a time wheel with a base tick of minTick, and the number of slots is 1 second / minTick, then the timer will create time wheels with a base tick of 1 second, 1 minute, 1 hour, etc. until the maxDelay is covered
func NewTimer(minTick time.Duration, maxDelay time.Duration) (*Timer, error) {
	timer := &Timer{
		currentTime: time.Now(),
		timeWheels:  []*TimeWheel{},
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

	for _, timeWheel := range t.timeWheels {
		if task.delay <= timeWheel.maxDelay {
			ticks := int((task.delay+timeWheel.baseTick-1)/timeWheel.baseTick) - 1
			slotIndex := (timeWheel.cursor + ticks) % timeWheel.wheelSize
			timeWheel.slots[slotIndex] = append(timeWheel.slots[slotIndex], task)
			return nil
		}
	}
	return fmt.Errorf("task delay is too large")
}

// Tick advances the timer by one tick (the smallest base tick of the time wheels) and executes all tasks that are due at the current time. It also handles the cascading of tasks from higher-level time wheels to lower-level ones as time advances.
func (t *Timer) Tick() {
	var tasks []*Task
	for i, timeWheel := range t.timeWheels {
		if i == 0 {
			tasks = timeWheel.slots[timeWheel.cursor]
			timeWheel.slots[timeWheel.cursor] = []*Task{} // drain the current slot
			timeWheel.cursor = (timeWheel.cursor + 1) % timeWheel.wheelSize
		} else {
			lastTimeWheel := t.timeWheels[i-1]
			timeWheel.cursor = (timeWheel.cursor + 1) % timeWheel.wheelSize
			for _, task := range timeWheel.slots[timeWheel.cursor] { // cascading tasks to lower-level time wheel
				task.delay = task.delay % timeWheel.baseTick
				if task.delay == 0 {
					task.delay = timeWheel.baseTick
				}
				ticks := int((task.delay+lastTimeWheel.baseTick-1)/lastTimeWheel.baseTick) - 1
				slotIndex := (lastTimeWheel.cursor + ticks) % lastTimeWheel.wheelSize
				lastTimeWheel.slots[slotIndex] = append(lastTimeWheel.slots[slotIndex], task)
			}
			timeWheel.slots[timeWheel.cursor] = []*Task{} // drain the current slot
		}
		if timeWheel.cursor != 0 {
			break
		}
	}
	for _, task := range tasks {
		task.Callback()
	}
}
