CREATE TABLE alerts (
                        id INT AUTO_INCREMENT PRIMARY KEY,        -- 自动递增的唯一主键
                        instance VARCHAR(255) NOT NULL,           -- 告警实例（通常是主机或服务的地址）
                        job VARCHAR(255) NOT NULL,                -- 告警的作业名称（job 标签）
                        alertname VARCHAR(255) NOT NULL,          -- 告警名称
                        status VARCHAR(50) NOT NULL,              -- 告警状态（如 firing 或 resolved）
                        summary TEXT,                             -- 告警摘要信息
                        start_time DATETIME NOT NULL,             -- 告警开始时间
                        end_time DATETIME,                        -- 告警结束时间（如果已恢复）
                        severity VARCHAR(50) NOT NULL,            -- 告警级别（从 Alertmanager 获取）
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 记录创建时间
);
