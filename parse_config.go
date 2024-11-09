package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// 定义配置文件对应的结构体
type Config struct {
	Server       ServerConfig       `yaml:"server" json:"server,omitempty"`
	MySQL        MySQLConfig        `yaml:"mysql" json:"my_sql"`
	Webhooks     []Webhook          `yaml:"webhooks" json:"webhooks,omitempty"`
	NoticePerson NoticePersonConfig `yaml:"notice_person" json:"notice_person"`
	SelfHeal     []SelfHealRule     `yaml:"self_heal" json:"self_heal,omitempty"`
}

type ServerConfig struct {
	Port      int    `yaml:"port"`
	URLSuffix string `yaml:"url_suffix"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Webhook struct {
	Job          string        `yaml:"job"`
	IMURL        string        `yaml:"im_url"`
	VoiceURL     string        `yaml:"voice_url,omitempty"`
	SilenceRules []SilenceRule `yaml:"voice_silence_rules,omitempty"`
}

type MatchConditions struct {
	Conditions map[string]string `yaml:",inline"`              // 正向或取反匹配条件，标签键值对
	TimeRange  *TimeRange        `yaml:"time_range,omitempty"` // 时间范围，包含 start、end 和 days 子字段
}

type TimeRange struct {
	Start string   `yaml:"start,omitempty"` // 静默期开始时间，格式为 "HH:MM" 或 "all-day"
	End   string   `yaml:"end,omitempty"`   // 静默期结束时间，格式为 "HH:MM" 或 "all-day"
	Days  []string `yaml:"days,omitempty"`  // 静默的特定日期（如 ["Saturday", "Sunday"]）
}

type SilenceRule struct {
	Match    *MatchConditions `yaml:"match,omitempty"`     // 正向匹配条件
	NotMatch *MatchConditions `yaml:"not_match,omitempty"` // 取反匹配条件
}

type NoticePersonConfig struct {
	Default map[string][]string  `yaml:"default"`
	Custom  []CustomNoticePerson `yaml:"custom"`
}

type CustomNoticePerson struct {
	Match        map[string]string `yaml:"match"`
	NoticePerson []string          `yaml:"notice_person"`
}

// SelfHealRule 定义自愈规则
type SelfHealRule struct {
	Match  map[string]string `yaml:"match"`
	Action string            `yaml:"action"`
	Delay  string            `yaml:"delay,omitempty"` // 延迟时间，带单位的字符串，可选
}

var config Config

// init 函数负责加载配置文件
func init() {
	configPath := filepath.Join("config", "config.yaml")
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}
}
