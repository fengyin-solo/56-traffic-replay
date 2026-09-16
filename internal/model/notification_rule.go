package model

import (
	"strings"
	"time"
)

const (
	NotifyChannelEmail  = "email"
	NotifyChannelWebhook = "webhook"
	NotifyChannelSMS    = "sms"
)

const (
	NotifyStatusEnabled  = "enabled"
	NotifyStatusDisabled = "disabled"
)

type NotificationRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Channel   string    `json:"channel"`
	Target    string    `json:"target"`
	Events    string    `json:"events"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *NotificationRule) Validate() error {
	n.Name = strings.TrimSpace(n.Name)
	n.Target = strings.TrimSpace(n.Target)
	if n.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if n.Channel == "" {
		return NewValidationError("channel", "通知渠道不能为空")
	}
	if n.Channel != NotifyChannelEmail && n.Channel != NotifyChannelWebhook && n.Channel != NotifyChannelSMS {
		return NewValidationError("channel", "通知渠道不合法")
	}
	if n.Target == "" {
		return NewValidationError("target", "通知目标不能为空")
	}
	if n.Status == "" {
		n.Status = NotifyStatusEnabled
	}
	if n.Status != NotifyStatusEnabled && n.Status != NotifyStatusDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type NotificationRuleFilter struct {
	Name    string
	Channel string
	Status  string
	Keyword string
}

func (f NotificationRuleFilter) Match(n *NotificationRule) bool {
	if f.Name != "" && n.Name != f.Name {
		return false
	}
	if f.Channel != "" && n.Channel != f.Channel {
		return false
	}
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(n.Name), k) &&
			!strings.Contains(strings.ToLower(n.Target), k) {
			return false
		}
	}
	return true
}
