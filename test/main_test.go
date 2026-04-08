package main

import "testing"

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		// t.Errorf 用于记录错误但不会立即终止测试
		t.Errorf("Add(2, 3) = %d, want %d", got, want)
	}

}
