package utils

import (
	"context"
	"fmt"
	"time"
)

const (
	Seconds         int64 = 1
	MinuteInSeconds       = 60 * Seconds
	HourInSeconds         = 60 * MinuteInSeconds
	DayInSeconds          = 24 * HourInSeconds
	MonthInSeconds        = 30 * DayInSeconds
	YearInSeconds         = 12 * MonthInSeconds
)

// GetRelativeTime 最近一段时间文本生成方法
// 逻辑与客户端渲染逻辑一致。其中有几个小细节： (1) 1 小时 50 分钟前，会渲染为 1 小时前  (2) 30 天计算为一个月.
func GetRelativeTime(ctx context.Context, timeInSeconds int64) string {
	now := time.Now().Unix()
	duration := now - timeInSeconds
	switch {
	case duration < MinuteInSeconds:
		return "刚刚"
	case duration < HourInSeconds:
		return fmt.Sprintf("%d 分钟前", duration/MinuteInSeconds)
	case duration < DayInSeconds:
		return fmt.Sprintf("%d 小时前", duration/HourInSeconds)
	case duration < MonthInSeconds:
		return fmt.Sprintf("%d 天前", duration/DayInSeconds)
	case duration < YearInSeconds:
		return fmt.Sprintf("%d 月前", duration/MonthInSeconds)
	default:
		return fmt.Sprintf("%d 年前", duration/YearInSeconds)
	}
}
