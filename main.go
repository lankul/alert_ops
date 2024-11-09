package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	port := validatePort(config.Server.Port)
	urlSuffix := validateURLSuffix(config.Server.URLSuffix)

	router.POST(fmt.Sprintf("%s", urlSuffix), handleAlerts)
	router.Run(fmt.Sprintf(":%d", port))
}

// validatePort 验证端口号是否在合法范围内，如果不合法则返回默认值 8080
func validatePort(port int) int {
	if port >= 1000 && port <= 65535 {
		return port
	}
	return 8080
}

// validateURLSuffix 验证URL后缀是否为空，如果为空则返回默认值 "alerts"
func validateURLSuffix(suffix string) string {
	if suffix != "" && suffix[0] == '/' {
		return suffix
	}
	return "alerts"
}
