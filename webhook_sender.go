package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func processAlert(labels map[string]string) {
	for _, webhook := range config.Webhooks {
		if webhook.Job == labels["job"] {
			// 发送到 IM webhook
			if labels["severity"] == "1" || labels["severity"] == "1.5" {
				sendToIMWebhook(webhook.IMURL, labels)
			}

			// 如果 severity 为 1，并且不在静默期内，才发送到语音 webhook
			if labels["severity"] == "1" && labels["status"] == "告警通知" && webhook.VoiceURL != "" && !isInSilencePeriod(labels, webhook.SilenceRules) {
				noticePerson := getNoticePerson(labels)
				sendToVoiceWebhook(webhook.VoiceURL, labels, noticePerson)
			}
		}
	}

	if labels["status"] == "告警通知" && strings.Contains(labels["summary"], "故障自愈") {
		executeSelfHeal(labels)
	}
}

// sendToIMWebhook 发送告警信息到 IM webhook
func sendToIMWebhook(url string, labels map[string]string) {
	message := fmt.Sprintf("类型：%s\n告警名称：%s\n实例：%s\n信息：%s\n发生时间：%s\n恢复时间：%s",
		labels["status"], labels["alertname"], labels["instance"], labels["summary"], labels["startTime"], labels["endTime"])
	payload := map[string]string{"content": message}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("序列化消息出错: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("发送 IM webhook 出错: %v", err)
	} else {
		defer resp.Body.Close()
		log.Printf("IM告警通知发送到 %s, 告警名称: %s, %v", url, labels["alertname"], resp.Status)
	}
}

// sendToVoiceWebhook 发送告警信息到语音 webhook
func sendToVoiceWebhook(url string, labels map[string]string, noticePerson []string) {
	payload := map[string]interface{}{
		"business":      labels["alertname"],
		"host_ip":       labels["instance"],
		"status":        labels["status"],
		"startAt":       labels["startTime"],
		"monitor":       "prometheus",
		"notice_person": noticePerson, // 语音通知人员
		"group":         labels["job"],
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("序列化消息出错: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("发送语音 webhook 出错: %v", err)
	} else {
		defer resp.Body.Close()
		log.Printf("语音通知发送到 %s, 告警名称: %s, %v", url, labels["alertname"], resp.Status)
	}
}
