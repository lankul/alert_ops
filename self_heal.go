package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

func executeSelfHeal(labels map[string]string) {
	for _, rule := range config.SelfHeal {
		matches := true
		for key, value := range rule.Match {
			labelValue, exists := labels[key]
			if !exists || !strings.Contains(labelValue, value) {
				matches = false
				break
			}

		}
		if matches {
			// 构建标签信息，传递给自愈程序
			labelArgs := make([]string, 0, len(labels))
			for key, value := range labels {
				labelArgs = append(labelArgs, key+"="+value)
			}

			// 默认情况下不延迟执行
			delayDuration := time.Duration(0)

			// 如果 delay 字段存在且不为 0 或者空字符串，则解析并延迟执行
			var err error
			if rule.Delay != "" && rule.Delay != "0" {
				delayDuration, err = time.ParseDuration(rule.Delay)
				if err != nil {
					log.Fatalf("Invalid delay format: %v", err)
				}
			}

			// 如果需要延迟执行
			if delayDuration > 0 {
				log.Printf("Waiting for %s before executing the action...", delayDuration)
				time.Sleep(delayDuration)
			}

			// 通过命令行参数传递标签
			cmd := exec.Command("/bin/sh", "-c", "config/"+rule.Action+" "+strings.Join(labelArgs, " "))
			log.Printf("Executing command: %s", cmd.String())

			// 通过环境变量传递标签（可选）
			//cmd.Env = append(os.Environ(), labelArgs...)

			// 执行自愈脚本
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("自愈操作失败: %v, 输出: %s", err, string(output))
			} else {
				log.Printf("自愈操作成功, 输出: %s", string(output))
			}
		}
	}
}
