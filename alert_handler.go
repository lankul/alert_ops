package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handleAlerts(c *gin.Context) {
	var alertData map[string]interface{}
	if err := c.BindJSON(&alertData); err != nil {
		log.Printf("解析告警数据出错: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求错误"})
		return
	}

	alerts := alertData["alerts"].([]interface{})
	for _, alert := range alerts {
		alertMap := alert.(map[string]interface{})
		labels := make(map[string]string)

		// 动态获取所有标签信息
		for key, value := range alertMap["labels"].(map[string]interface{}) {
			if strVal, ok := value.(string); ok {
				labels[key] = strVal
			}
		}

		status := alertMap["status"].(string)
		startTimeUTC := alertMap["startsAt"].(string)
		endTimeUTC := alertMap["endsAt"].(string)
		summary, ok := alertMap["annotations"].(map[string]interface{})["summary"].(string)
		if !ok {
			log.Printf("alertname: %s, summary is missing", labels["alertname"])
			summary = ""
		}
		// 转换 status 字段
		status = strings.Replace(status, "resolved", "恢复通知", -1)
		status = strings.Replace(status, "firing", "告警通知", -1)

		// 转换 startTime 和 endTime 到北京时间
		startTime, err := time.Parse(time.RFC3339, startTimeUTC)
		if err != nil {
			log.Printf("解析开始时间出错: %v", err)
			continue
		}
		startTime = startTime.In(time.FixedZone("CST", 8*3600))
		startTimeStr := startTime.Format(time.RFC3339)

		endTime, err := time.Parse(time.RFC3339, endTimeUTC)
		if err != nil {
			log.Printf("解析结束时间出错: %v", err)
			continue
		}
		endTime = endTime.In(time.FixedZone("CST", 8*3600))
		endTimeStrTmp := endTime.Format(time.RFC3339)
		var endTimeStr string
		if endTimeStrTmp == "0001-01-01T08:00:00+08:00" {
			endTimeStr = ""
		} else {
			endTimeStr = endTimeStrTmp
		}

		labels["summary"] = summary
		labels["status"] = status
		labels["startTime"] = startTimeStr
		labels["endTime"] = endTimeStr

		// 存储告警信息到数据库
		var endTimeSql sql.NullString
		if endTimeStr == "" {
			endTimeSql = sql.NullString{
				String: "",
				Valid:  false,
			}
		} else {
			endTimeSql = sql.NullString{
				String: endTimeStr,
				Valid:  true,
			}
		}
		storeAlertInDB(labels, endTimeSql)

		// 处理告警并发送到 Webhook
		processAlert(labels)
	}

	c.JSON(http.StatusOK, gin.H{"message": "告警已处理"})
}
