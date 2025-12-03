package rules

import (
	"github.com/emersion/go-imap"
)

type Rule interface {
	ShouldDelete(msg *imap.Message) bool
}

type Rules struct {
	rules []Rule
}

func NewRules(rules []Rule) *Rules {
	return &Rules{rules: rules}
}

func (r *Rules) ShouldDelete(msg *imap.Message) bool {
	for _, rule := range r.rules {
		if rule.ShouldDelete(msg) {
			return true
		}
	}
	return false
}
