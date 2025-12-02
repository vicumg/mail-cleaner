package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/emersion/go-imap"
)

type Rule interface {
	ShouldDelete(msg *imap.Message) bool
}

type Rules struct {
	rules []Rule
}

type FromAddressRule struct {
	Address string
}

type FromDomainRule struct {
	Domain string
}

func (r *FromDomainRule) ShouldDelete(msg *imap.Message) bool {
	if msg.Envelope == nil {
		return false
	}
	for _, addr := range msg.Envelope.From {
		// if domain is a part of the email address check for substring match for mail.citrus.com.ua emasiles.citrus.com.ua
		if strings.Contains(addr.HostName, r.Domain) {
			fmt.Printf("Deleting email from domain: %s\n", r.Domain)
			return true
		}
	}
	return false
}

func (r *FromAddressRule) ShouldDelete(msg *imap.Message) bool {
	if msg.Envelope == nil {
		return false
	}
	for _, addr := range msg.Envelope.From {
		if addr.MailboxName+"@"+addr.HostName == r.Address {
			fmt.Printf("Deleting email from: %s\n", r.Address)
			return true
		}
	}
	return false
}

func NewRules(rule_set_file string) *Rules {
	var rules []Rule

	// load data from json file and create rules
	file_data, err := os.ReadFile(rule_set_file)
	if err != nil {
		panic("Failed to read rules file: " + err.Error())
	}

	var raw_rules []map[string]interface{}
	if err := json.Unmarshal(file_data, &raw_rules); err != nil {
		panic("Failed to parse rules file: " + err.Error())
	}

	for _, raw_rule := range raw_rules {
		if raw_rule["type"] == "FromAddressRule" {
			address, ok := raw_rule["address"].(string)
			if !ok {
				panic("Invalid address in FromAddressRule")
			}
			rules = append(rules, &FromAddressRule{Address: address})
		}
		// Domain rule
		if raw_rule["type"] == "FromDomainRule" {
			domain, ok := raw_rule["address"].(string)
			if !ok {
				panic("Invalid address in FromDomainRule")
			}
			rules = append(rules, &FromDomainRule{Domain: domain})
		}

		// Add more rule types here as needed
	}

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
