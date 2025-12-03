package rule

import (
	"errors"
	"fmt"
	"mail-cleaner/internal/rules"
	"strings"

	"github.com/emersion/go-imap"
)

type AddressRule struct {
	Address string
}

func init() {
	RegisterRuleFactory("address_rule", func(data map[string]any) (rules.Rule, error) {
		address, ok := data["address"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid or missing 'address' field")
		}
		return NewAddressRule(address)
	})
}

func NewAddressRule(address string) (*AddressRule, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}
	return &AddressRule{
		Address: address,
	}, nil
}

func (r *AddressRule) ShouldDelete(msg *imap.Message) bool {
	if msg.Envelope == nil {
		return false
	}
	for _, addr := range msg.Envelope.From {
		if r.apply(addr.MailboxName+"@"+addr.HostName, r.Address) {
			fmt.Printf("Deleting email from: %s\n", r.Address)
			return true
		}
	}
	return false
}

func (r *AddressRule) apply(emailAddress string, ruleAddress string) bool {
	return strings.EqualFold(emailAddress, ruleAddress)
}
