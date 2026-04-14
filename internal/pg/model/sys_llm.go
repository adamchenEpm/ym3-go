package pgmodel

import (
	"time"
)

type SysLlm struct {
	Id         int       `pg:"id" json:"id,omitempty"`
	Name       string    `pg:"name" json:"name,omitempty"`
	Type       string    `pg:"type" json:"type,omitempty"`
	BaseUrl    string    `pg:"base_url" json:"baseUrl,omitempty"`
	APIKey     string    `pg:"api_key" json:"apiKey,omitempty"`
	ModelId    string    `pg:"model_id" json:"modelId,omitempty"`
	ModelName  string    `pg:"model_name" json:"modelName,omitempty"`
	Remark     string    `pg:"remark" json:"remark,omitempty"`
	CreateTime time.Time `pg:"create_time" json:"createTime,omitempty"`
	UpdateTime time.Time `pg:"update_time" json:"updateTime,omitempty"`
	Status     int       `pg:"status" json:"status,omitempty"`
}

var (
	SysLlmSelect = `SELECT id, name, type, base_url, api_key, model_id, model_name,
	remark, create_time, update_time, status FROM sys_llm `
)
