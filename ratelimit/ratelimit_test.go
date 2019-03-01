package ratelimit

import (
	"fmt"
	"testing"
)

func TestRate(t *testing.T) {

	bucket := NewTokenBucket("test", nil)
	for i := 0; i < 500; i++ {
		r := bucket.Take(100000)
		fmt.Printf("结果:%t \r\n", r)
	}
}
