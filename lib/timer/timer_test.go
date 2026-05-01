package timer

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewTimerValidation(t *testing.T) {
	tests := []struct {
		name     string
		minTick  time.Duration
		maxDelay time.Duration
		wantErr  bool
	}{
		{name: "minTick must be positive", minTick: 0, maxDelay: time.Second, wantErr: true},
		{name: "maxDelay must be >= minTick", minTick: time.Second, maxDelay: 500 * time.Millisecond, wantErr: true},
		{name: "sub-second minTick must divide one second", minTick: 333 * time.Millisecond, maxDelay: time.Second, wantErr: true},
		{name: "valid one second precision", minTick: time.Second, maxDelay: time.Hour, wantErr: false},
		{name: "valid sub-second precision", minTick: 10 * time.Millisecond, maxDelay: 2 * time.Second, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTimer(tt.minTick, tt.maxDelay)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewTimer() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddTaskValidation(t *testing.T) {
	timer, err := NewTimer(time.Second, time.Minute)
	if err != nil {
		t.Fatalf("创建定时器失败: %v", err)
	}

	if err := timer.AddTask(nil); err == nil {
		t.Fatal("AddTask(nil) 应返回错误")
	}
	if err := timer.AddTask(&Task{delay: 0}); err == nil {
		t.Fatal("delay=0 应返回错误")
	}
	if err := timer.AddTask(&Task{delay: -time.Second}); err == nil {
		t.Fatal("负 delay 应返回错误")
	}
}

func TestTickExecutesAtExpectedDelays(t *testing.T) {
	timer, err := NewTimer(time.Second, time.Hour)
	if err != nil {
		t.Fatalf("创建定时器失败: %v", err)
	}

	executed := make(chan string, 3)
	for _, tc := range []struct {
		name  string
		delay time.Duration
	}{
		{name: "1s", delay: 1 * time.Second},
		{name: "59s", delay: 59 * time.Second},
		{name: "61s", delay: 61 * time.Second},
	} {
		if err := timer.AddTask(&Task{
			Callback: func() { executed <- tc.name },
			delay:    tc.delay,
		}); err != nil {
			t.Fatalf("添加任务 %s 失败: %v", tc.name, err)
		}
	}

	for i := 1; i <= 61; i++ {
		timer.Tick()

		var got []string
		for {
			select {
			case name := <-executed:
				got = append(got, name)
			default:
				goto drained
			}
		}

	drained:
		var want []string
		switch i {
		case 1:
			want = []string{"1s"}
		case 59:
			want = []string{"59s"}
		case 61:
			want = []string{"61s"}
		default:
			want = nil
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("tick=%d, got=%v, want=%v", i, got, want)
		}
	}
}

func TestRunExecutesTasksAndCanStop(t *testing.T) {
	timer, err := NewTimer(10*time.Millisecond, time.Second)
	if err != nil {
		t.Fatalf("创建定时器失败: %v", err)
	}

	executed := make(chan string, 2)
	if err := timer.AddTask(&Task{
		Callback: func() { executed <- "30ms" },
		delay:    30 * time.Millisecond,
	}); err != nil {
		t.Fatalf("添加 30ms 任务失败: %v", err)
	}
	if err := timer.AddTask(&Task{
		Callback: func() { executed <- "50ms" },
		delay:    50 * time.Millisecond,
	}); err != nil {
		t.Fatalf("添加 50ms 任务失败: %v", err)
	}

	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		timer.Run(done)
	}()

	first := <-executed
	second := <-executed
	if !reflect.DeepEqual([]string{first, second}, []string{"30ms", "50ms"}) {
		t.Fatalf("Run 执行顺序错误，got=%v", []string{first, second})
	}

	close(done)
	select {
	case <-stopped:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("Run 在 close(done) 后未及时退出")
	}
}

func TestConcurrentAddTask(t *testing.T) {
	timer, err := NewTimer(time.Second, time.Minute)
	if err != nil {
		t.Fatalf("创建定时器失败: %v", err)
	}

	const n = 200
	var wg sync.WaitGroup
	wg.Add(n)
	errCh := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			errCh <- timer.AddTask(&Task{Callback: func() {}, delay: time.Second})
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("并发 AddTask 返回错误: %v", err)
		}
	}

	if got := timer.timeWheels[0].slots[0].Len(); got != n {
		t.Fatalf("并发添加后任务数不匹配，got=%d, want=%d", got, n)
	}
}
