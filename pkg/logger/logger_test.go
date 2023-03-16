package logger

import "testing"

func TestLogger(t *testing.T) {
	//panic recover
	defer func() {
		if err := recover(); err != nil {
			t.Logf("panic (by panicf) recover: %v", err)
		}
	}()
	Init()
	Infof("hello world")
	Warnf("warn %s", "world")
	Errorf("error %s", "world")
	Panicf("panic %s", "world")
}
