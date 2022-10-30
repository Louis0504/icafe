package utils

import (
	"sync"
	"time"
)

type WarmingUpRateLimiter struct {
	maxToken        int64 // 令牌上限
	warmUpPeriod    int64 // 爬坡时长, 秒
	step            int64 // 爬坡步长
	currentMaxToken int64 // 当前令牌上限

	lastQPS    int64 // 前一秒的 QPS
	currentQPS int64 // 当前秒的 QPS

	lastAccess time.Time     // 上一次访问时间
	interval   time.Duration // 两次发令牌时间间隔

	lock sync.Mutex
}

func NewWarmingUpRateLimiter(maxToken, warmUpPeriod int64) *WarmingUpRateLimiter {
	step := MaxInt64(maxToken/warmUpPeriod, 1)
	w := &WarmingUpRateLimiter{
		maxToken:        maxToken,
		step:            step,
		warmUpPeriod:    warmUpPeriod,
		currentMaxToken: step, // 冷启动，直接爬坡
		interval:        time.Duration(int64(time.Second) / (2 * step)),
		lastAccess:      time.Now(),
	}

	go func() {
		for range time.Tick(1 * time.Second) {
			w.updateToken()
		}
	}()

	return w
}

// Take 持续等待.
func (w *WarmingUpRateLimiter) Take() {
	for {
		w.lock.Lock()
		now := time.Now()
		if w.currentQPS < w.currentMaxToken {
			if desire := w.lastAccess.Add(w.interval); now.Before(desire) {
				w.lock.Unlock()
				time.Sleep(desire.Sub(now))
				continue
			}
			w.lastAccess = now
			w.currentQPS++
			w.lock.Unlock()
			return
		}
		w.lock.Unlock()
	}
}

// updateToken 更新令牌桶.
func (w *WarmingUpRateLimiter) updateToken() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.lastQPS = w.currentQPS
	w.currentQPS = int64(0)

	w.currentMaxToken = MinInt64(w.lastQPS+w.step, w.maxToken)
	w.interval = time.Duration(int64(time.Second) / (2 * w.currentMaxToken))
}

// SetMaxToken 设置令牌上限.
func (w *WarmingUpRateLimiter) SetLimit(limit int64) {
	w.SetLimitAndWarmingPeriod(limit, w.warmUpPeriod)
}

// SetWarmUpPeriod 设置爬坡时长.
func (w *WarmingUpRateLimiter) SetWarmUpPeriod(period int64) {
	w.SetLimitAndWarmingPeriod(w.maxToken, period)
}

// SetMaxTokenAndWarmingPeriod 设置令牌上限与爬坡时长.
func (w *WarmingUpRateLimiter) SetLimitAndWarmingPeriod(limit, period int64) {
	if limit <= 0 {
		return
	}
	w.lock.Lock()
	w.maxToken = limit
	w.warmUpPeriod = period
	w.step = MaxInt64(limit/period, 1)
	w.lock.Unlock()
}

func (w *WarmingUpRateLimiter) GetCurrentStatus() (lastQPS, maxLimit, currentMax, warmUpPeriod int64) {
	w.lock.Lock()
	lastQPS, maxLimit, currentMax, warmUpPeriod = w.lastQPS, w.maxToken, w.currentMaxToken, w.warmUpPeriod
	w.lock.Unlock()
	return
}
