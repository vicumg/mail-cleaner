package rule

import (
	"fmt"

	"mail-cleaner/internal/rules"

	"github.com/emersion/go-imap"
)

type AIRule struct {
	Enabled bool `json:"enabled"`
}

func init() {
	RegisterRuleFactory("ai_rule", func(config map[string]interface{}) (rules.Rule, error) {
		enabled, ok := config["enabled"].(bool)
		if !ok {
			enabled = false
		}
		return NewAIRule(enabled)
	})
}

func NewAIRule(enabled bool) (*AIRule, error) {
	return &AIRule{
		Enabled: enabled,
	}, nil
}

func (r *AIRule) ShouldDelete(msg *imap.Message) bool {
	if !r.Enabled {
		return false
	}

	if msg.Envelope == nil || len(msg.Envelope.From) == 0 {
		return false
	}

	for _, addr := range msg.Envelope.From {
		emailAddress := addr.MailboxName + "@" + addr.HostName
		subject := msg.Envelope.Subject
		if r.apply(emailAddress, subject) {
			return true
		}
	}

	return false
}

func (r *AIRule) apply(emailAddress string, subject string) bool {
	// TODO: implement AI-based logic
	return false
}

func (r *AIRule) String() string {
	return fmt.Sprintf("AIRule{Enabled: %v}", r.Enabled)
}
