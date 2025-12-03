package rule

import (
	"errors"
	"fmt"
	"mail-cleaner/internal/rules"
	"strings"

	"github.com/emersion/go-imap"
)

type DomainRule struct {
	Domain string
}

func init() {
	RegisterRuleFactory("domain_rule", func(data map[string]any) (rules.Rule, error) {
		domain, ok := data["domain"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid or missing 'domain' field")
		}
		return NewDomainRule(domain)
	})
}

func NewDomainRule(domain string) (*DomainRule, error) {
	if domain == "" {
		return nil, errors.New("domain cannot be empty")
	}
	return &DomainRule{
		Domain: domain,
	}, nil
}

func (r *DomainRule) ShouldDelete(msg *imap.Message) bool {
	if msg.Envelope == nil {
		return false
	}
	for _, addr := range msg.Envelope.From {
		if strings.Contains(addr.HostName, r.Domain) {
			fmt.Printf("Deleting email from domain: %s\n", r.Domain)
			return true
		}
	}
	return false
}
