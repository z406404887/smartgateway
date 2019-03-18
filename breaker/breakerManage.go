package breaker

import "time"

var (
	BreakerList = make(map[string]*CircuitBreaker) //熔断器列表
)

func NewBreaker(name, triggerOpen func(Counts) bool, maxRequest int, beHalOpenInterval time.Duration, clearInterval time.Duration) *CircuitBreaker {

}
