package main

import "strings"

// getNoticePerson 根据告警标签获取通知人列表
func getNoticePerson(labels map[string]string) []string {
	for _, custom := range config.NoticePerson.Custom {
		matches := true
		for key, value := range custom.Match {
			if labelValue, exists := labels[key]; !exists || !strings.Contains(labelValue, value) {
				matches = false
				break
			}
		}
		if matches {
			return custom.NoticePerson
		}
	}
	if defaultNoticePerson, exists := config.NoticePerson.Default[labels["job"]]; exists {
		return defaultNoticePerson
	}
	return []string{}
}
