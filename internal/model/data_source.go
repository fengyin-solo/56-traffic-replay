package model

import (
	"strings"
	"time"
)

const (
	DataSourceTypeMySQL    = "mysql"
	DataSourceTypePostgreSQL = "postgresql"
	DataSourceTypeHTTP     = "http"
)

type DataSource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *DataSource) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	d.Config = strings.TrimSpace(d.Config)
	if d.Name == "" {
		return NewValidationError("name", "数据源名称不能为空")
	}
	if d.Type == "" {
		return NewValidationError("type", "数据源类型不能为空")
	}
	if d.Type != DataSourceTypeMySQL && d.Type != DataSourceTypePostgreSQL && d.Type != DataSourceTypeHTTP {
		return NewValidationError("type", "数据源类型不合法")
	}
	if d.Config == "" {
		return NewValidationError("config", "配置不能为空")
	}
	if d.Status == "" {
		d.Status = "active"
	}
	return nil
}

type DataSourceFilter struct {
	Name    string
	Type    string
	Status  string
	Keyword string
}

func (f DataSourceFilter) Match(d *DataSource) bool {
	if f.Name != "" && d.Name != f.Name {
		return false
	}
	if f.Type != "" && d.Type != f.Type {
		return false
	}
	if f.Status != "" && d.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Name), k) {
			return false
		}
	}
	return true
}
