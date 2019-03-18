package breaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type doFunc func(context.Context) error

type failFunc func(context.Context) error

//状态
const (
	Closed   = iota //关闭状态， 健康允许访问
	HalfOpen        //半开状态， 允许一部分 访问
	Open            //开启状态， 将阻挡访问
)

var (
	defaultMaxRequest = 15
)

//CircuitBreaker 断路器对象
type CircuitBreaker struct {
	Name              string
	mu                sync.Mutex
	TriggerOpen       func(Counts) bool //根据func 的返回条件 触发 Open状态
	MaxRequest        int               //当是 半开状态时， 允许继续请求的数目
	BeHalOpenInterval time.Duration     //当进入Open 状态时， 定时变为 HalfOpen的时间
	ClearInterval     time.Duration     //当是Close状态时， 定时清除的时间
	Count             Counts            //计时器
	Status            int               //当前状态

	viewTime time.Time // 最后访问时间
}

//Counts 计数器
type Counts struct {
	SuccessCounts    int //成功次数
	FailCounts       int //失败次数
	ContinuesSuccess int //连续失败次数
	ContinuesFail    int //连续成功次数
	Totals           int //总请求数
}

func (breaker *CircuitBreaker) SetConfig(triggerOpen func(Counts) bool, maxRequest int, beHalOpenInterval time.Duration, clearInterval time.Duration) {
	if triggerOpen == nil {
		breaker.TriggerOpen = defaultTriggerOpen
	} else {
		breaker.TriggerOpen = triggerOpen
	}

	if maxRequest <= 0 {
		breaker.MaxRequest = defaultMaxRequest
	} else {
		breaker.MaxRequest = maxRequest
	}

	breaker.BeHalOpenInterval = beHalOpenInterval

	breaker.ClearInterval = clearInterval
}

//Handle 执行请求， 记录相关行为
func (breaker *CircuitBreaker) Handle(context context.Context, doFunc, failback failFunc) error {

	_, err := breaker.beforeHandle()

	if err != nil {
		return failback(context)
	}

	//成功， 继续执行
	isSuccess := true
	err = doFunc(context)

	//执行失败
	if err != nil {
		isSuccess = false
	}

	//处理执行后
	breaker.afterHandle(isSuccess)

	if !isSuccess {
		return failback(context)
	}

	return nil
}

//ReStartCount 重新计数
func (breaker *CircuitBreaker) ReStartCount() {
	breaker.Count.SuccessCounts = 0
	breaker.Count.FailCounts = 0
	breaker.Count.ContinuesSuccess = 0
	breaker.Count.ContinuesFail = 0
	breaker.Count.Totals = 0
}

//AddSuccess 请求通过
func (breaker *CircuitBreaker) AddSuccess() {
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	breaker.Count.ContinuesSuccess++
	breaker.Count.ContinuesFail = 0
	breaker.Count.SuccessCounts++
	breaker.Count.Totals++
}

//AddSuccess 请求失败
func (breaker *CircuitBreaker) AddFail() {
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	breaker.Count.ContinuesFail++
	breaker.Count.ContinuesSuccess = 0
	breaker.Count.FailCounts++
	breaker.Count.Totals++
}

//执行请求前
func (breaker *CircuitBreaker) beforeHandle() (Status int, err error) {
	//开始增加次数
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	Status = breaker.Status

	if breaker.Status == Closed { //关闭状态
		//需清零错误
		if time.Now().After(breaker.viewTime.Add(breaker.ClearInterval)) {
			breaker.ReStartCount()
		}

		//触发 Open
		if breaker.TriggerOpen(breaker.Count) {
			breaker.Status = Open
		}
	}

	if breaker.Status == Open { //打开状态
		//需重置为 半打开状态
		if time.Now().After(breaker.viewTime.Add(breaker.BeHalOpenInterval)) {
			breaker.Status = HalfOpen
		}
	}

	if breaker.Status == HalfOpen { //半打开状态
		//最大成功数已到，重置为 关闭状态
		if breaker.Count.ContinuesSuccess >= breaker.MaxRequest {
			breaker.Status = Closed
		}

		//超过允许的请求数
		if breaker.Count.Totals >= breaker.MaxRequest {
			return breaker.Status, errors.New("Too Many Handle")
		}
	} else if breaker.Status == Open {
		return breaker.Status, errors.New("Breaker Open")
	}

	//设置 访问时间
	breaker.viewTime = time.Now()

	return breaker.Status, nil
}

//执行请求后
func (breaker *CircuitBreaker) afterHandle(isSuccess bool) {
	//开始增加次数
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	if isSuccess {
		breaker.AddSuccess()
	} else {
		breaker.AddFail()
	}
}

//默认 30 次错误 触发 断路器开启
func defaultTriggerOpen(count Counts) bool {
	if count.ContinuesFail > 30 {
		return true
	}

	return false
}
