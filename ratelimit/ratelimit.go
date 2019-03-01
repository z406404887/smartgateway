package ratelimit

import (
	"sync"
	"time"
)

//TokenBucket 令牌桶对象
type TokenBucket struct {
	capacity        int64      //当前令牌桶的总容量
	putTimeInterval int64      //令牌放入桶里的时间间隔(PutTimeInterval，PutNumInterval 即表示每隔几秒放入多少令牌)
	putNumInterval  int64      //令牌放入桶里时的数量
	remainTokens    int64      //剩余的令牌数
	lastTakeTime    time.Time  //上一次拿令牌的时间
	mu              sync.Mutex //锁
}

//BucketConfig 令牌桶配置对象
type BucketConfig struct {
	Capacity        int64 //当前令牌桶的总容量
	PutTimeInterval int64 //令牌放入桶里的时间间隔(PutTimeInterval，PutNumInterval 即表示每隔几秒放入多少令牌)
	PutNumInterval  int64 //令牌放入桶里时的数量
	remainTokens    int64 //剩余的令牌数
}

//bucketsManage 内存存储所有使用令牌桶算法的
type bucketsManage struct {
	mu        sync.Mutex //锁
	AllBukets map[string]*TokenBucket
}

var (
	manageBucket      *bucketsManage
	dfCapacity        int64 = 1000000
	dfRemainTokens    int64 = 1000000
	dfPutTimeInterval int64 = 10
	dfPutNumInterval  int64 = 1000000
)

func init() {
	manageBucket = &bucketsManage{
		mu:        sync.Mutex{},
		AllBukets: make(map[string]*TokenBucket),
	}
}

//NewTokenBucket 返回令牌桶对象
func NewTokenBucket(name string, bucketConfig *BucketConfig) *TokenBucket {
	manageBucket.mu.Lock()
	defer manageBucket.mu.Unlock()

	//获取配置对象
	var config *BucketConfig
	if bucketConfig != nil {
		config = bucketConfig
	} else {
		config = &BucketConfig{
			Capacity:        dfCapacity,
			PutTimeInterval: dfPutTimeInterval,
			PutNumInterval:  dfPutNumInterval,
			remainTokens:    dfRemainTokens,
		}
	}

	//获取对应name 的Bucket对象
	bucket, ok := manageBucket.AllBukets[name]

	//还没有该name 的Bucket 对象
	if !ok {
		bucket = &TokenBucket{
			capacity:        config.Capacity,
			putTimeInterval: config.PutTimeInterval,
			putNumInterval:  config.PutNumInterval,
			remainTokens:    config.remainTokens,
			lastTakeTime:    time.Now(),
			mu:              sync.Mutex{},
		}

		manageBucket.AllBukets[name] = bucket
	}

	return bucket
}

//Take 从令牌桶中获取令牌
func (bucket *TokenBucket) Take(count int64) bool {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	//上一次的时间到现在新增的令牌
	now := time.Now()
	addTokens := (int64(now.Sub(bucket.lastTakeTime)) / (bucket.putTimeInterval * 1e9)) * bucket.putNumInterval

	bucket.remainTokens = bucket.remainTokens + addTokens

	//新加的令牌 已大于总量， 重置剩余量为总量
	if bucket.remainTokens >= bucket.capacity {
		bucket.remainTokens = bucket.capacity
	}

	//令牌桶里的令牌不足
	if bucket.remainTokens < count {
		return false
	}

	//可以获取令牌桶里的令牌， 更新相关信息
	bucket.lastTakeTime = now
	bucket.remainTokens = bucket.remainTokens - count

	return true
}
