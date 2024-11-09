package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 必须导入 MySQL 驱动
	"log"
)

func storeAlertInDB(labels map[string]string, endTime sql.NullString) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	query := "INSERT INTO alerts (instance, job, alertname, status, summary, start_time, end_time, severity) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, labels["instance"], labels["job"], labels["alertname"], labels["status"], labels["summary"], labels["startTime"], endTime, labels["severity"])
	if err != nil {
		log.Printf("插入告警数据失败: %v", err)
	}
}
