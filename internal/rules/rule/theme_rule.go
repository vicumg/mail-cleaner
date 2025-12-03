package rule

import (
	"errors"
	"fmt"
	"mail-cleaner/internal/rules"
	"strings"

	"github.com/emersion/go-imap"
)

type ThemeRule struct {
	Text string
}

func init() {
	RegisterRuleFactory("theme_rule", func(data map[string]any) (rules.Rule, error) {
		text, ok := data["text"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid or missing 'text' field")
		}
		return NewThemeRule(text)
	})
}

func NewThemeRule(text string) (*ThemeRule, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}
	return &ThemeRule{
		Text: text,
	}, nil
}

func (r *ThemeRule) ShouldDelete(msg *imap.Message) bool {
	if msg.Envelope == nil {
		return false
	}

	if msg.Envelope.Subject != "" && containsIgnoreCase(msg.Envelope.Subject, r.Text) {
		fmt.Printf("Deleting email with subject containing: %s\n", r.Text)
		return true
	}
	return false
}

func containsIgnoreCase(s1, s2 string) bool {
	s1Lower := strings.ToLower(s1)
	s2Lower := strings.ToLower(s2)
	return strings.Contains(s1Lower, s2Lower)
}
