package breaker

import (
	"sync"
	"time"
)

var (
	mu          = &sync.RWMutex{}
	BreakerList = make(map[string]*CircuitBreaker) //熔断器列表

)

// NewBreaker  获取Breaker
func NewBreaker(name string, triggerOpen func(Counts) bool, maxRequest int, beHalOpenInterval time.Duration, clearInterval time.Duration) *CircuitBreaker {

	//允许大量读
	mu.RLock()
	if _, ok := BreakerList[name]; ok {
		mu.RUnlock()
		return BreakerList[name]
	}
	mu.RUnlock()

	//写锁定
	mu.Lock()
	defer mu.Unlock()

	//到这里锁住， 但之前的读锁有可能有有多个并发，需再次检查
	if breaker, ok := BreakerList[name]; ok {
		return breaker
	}

	//初始化
	breaker := &CircuitBreaker{}
	breaker.setConfig(triggerOpen, maxRequest, beHalOpenInterval, clearInterval)

	//添加到map中
	BreakerList[name] = breaker

	return BreakerList[name]
}
