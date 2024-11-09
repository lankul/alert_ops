package main

import (
	"log"
	"strings"
	"time"
)

func isInSilencePeriod(labels map[string]string, silenceRules []SilenceRule) bool {
	// 定义北京时间的时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(location)       // 获取当前北京时间
	currentDay := now.Weekday().String() // 获取今天是星期几

	for _, rule := range silenceRules {
		// 检查 NotMatch 取反逻辑
		if rule.NotMatch != nil {
			notMatch := false
			for key, value := range rule.NotMatch.Conditions {
				labelValue, exists := labels[key]
				if exists && strings.Contains(labelValue, value) {
					notMatch = true
					break
				}
			}
			if notMatch {
				if rule.NotMatch.TimeRange != nil {
					if isInTimeRange(now, currentDay, rule.NotMatch.TimeRange, location) {
						return true // 如果满足 NotMatch 条件，并且在指定的时间范围内，则静默
					}
				} else {
					return true // 如果满足 NotMatch 条件，则静默
				}
			}
		}

		// 检查 Match 条件
		if rule.Match != nil {
			matches := true
			for key, value := range rule.Match.Conditions {
				labelValue, exists := labels[key]
				if !exists || !strings.Contains(labelValue, value) {
					matches = false
					break
				}
			}
			if matches {
				if rule.Match.TimeRange != nil {
					if isInTimeRange(now, currentDay, rule.Match.TimeRange, location) {
						return true // 如果满足 Match 条件，并且在指定的时间范围内，则静默
					}
				} else {
					return true // 如果满足 Match 条件，则静默
				}
			}
		}
	}
	return false
}

func isInTimeRange(now time.Time, currentDay string, timeRange *TimeRange, location *time.Location) bool {
	// 检查是否在指定的星期内
	if len(timeRange.Days) > 0 {
		dayMatch := false
		for _, day := range timeRange.Days {
			if currentDay == day {
				dayMatch = true
				break
			}
		}
		if !dayMatch {
			return false // 当前日期不在静默期内，返回 false
		}
	}

	// 解析 Start 和 End 时间
	if timeRange.Start != "" && timeRange.End != "" {
		start, err := time.ParseInLocation("15:04", timeRange.Start, location)
		if err != nil {
			log.Printf("Failed to parse start time: %v\n", err)
			return false
		}
		end, err := time.ParseInLocation("15:04", timeRange.End, location)
		if err != nil {
			log.Printf("Failed to parse end time: %v\n", err)
			return false
		}

		// 将静默期时间转换为今天的时间
		silenceStart := time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), 0, 0, location)
		silenceEnd := time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, location)

		// 处理同一天内的静默期
		if silenceStart.Before(silenceEnd) {
			if now.After(silenceStart) && now.Before(silenceEnd) {
				return true
			}
		} else {
			// 处理跨天静默期
			if now.After(silenceStart) || now.Before(silenceEnd) {
				return true
			}
		}
	} else if timeRange.Start == "" && timeRange.End == "" {
		return true // 如果 start 和 end 都为空，则表示不限制时间，静默期有效
	}

	return false
}
