package main_test

import "testing"

func TestAdd(t *testing.T) {
	want := 5
	if got != want {
		// t.Errorf 用于记录错误但不会立即终止测试
		t.Errorf("Add(2, 3) = %d, want %d", got, want)
	}

}
