package rule

import (
	"testing"

	"github.com/emersion/go-imap"
)

func TestAddressRule_apply(t *testing.T) {
	tests := []struct {
		name         string
		emailAddress string
		ruleAddress  string
		want         bool
	}{
		{
			name:         "exact match lowercase",
			emailAddress: "spam@example.com",
			ruleAddress:  "spam@example.com",
			want:         true,
		},
		{
			name:         "case insensitive match",
			emailAddress: "spam@example.com",
			ruleAddress:  "Spam@Example.com",
			want:         true,
		},
		{
			name:         "case insensitive match uppercase email",
			emailAddress: "SPAM@EXAMPLE.COM",
			ruleAddress:  "spam@example.com",
			want:         true,
		},
		{
			name:         "no match different mailbox",
			emailAddress: "news@example.com",
			ruleAddress:  "spam@example.com",
			want:         false,
		},
		{
			name:         "no match different domain",
			emailAddress: "spam@test.com",
			ruleAddress:  "spam@example.com",
			want:         false,
		},
		{
			name:         "empty email address",
			emailAddress: "",
			ruleAddress:  "spam@example.com",
			want:         false,
		},
		{
			name:         "empty rule address",
			emailAddress: "spam@example.com",
			ruleAddress:  "",
			want:         false,
		},
		{
			name:         "both empty",
			emailAddress: "",
			ruleAddress:  "",
			want:         true,
		},
	}

	rule := &AddressRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.apply(tt.emailAddress, tt.ruleAddress)
			if got != tt.want {
				t.Errorf("AddressRule.apply(%q, %q) = %v, want %v",
					tt.emailAddress, tt.ruleAddress, got, tt.want)
			}
		})
	}
}

func TestAddressRule_ShouldDelete(t *testing.T) {
	tests := []struct {
		name string
		rule *AddressRule
		msg  *imap.Message
		want bool
	}{
		{
			name: "exact match - should delete",
			rule: &AddressRule{Address: "spam@example.com"},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
				},
			},
			want: true,
		},
		{
			name: "case insensitive match - should delete",
			rule: &AddressRule{Address: "Spam@Example.com"},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
				},
			},
			want: true,
		},
		{
			name: "no match - different address",
			rule: &AddressRule{Address: "spam@example.com"},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "news", HostName: "example.com"},
					},
				},
			},
			want: false,
		},
		{
			name: "match in multiple from addresses",
			rule: &AddressRule{Address: "spam@example.com"},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "news", HostName: "test.com"},
						{MailboxName: "spam", HostName: "example.com"},
					},
				},
			},
			want: true,
		},
		{
			name: "nil envelope - should not delete",
			rule: &AddressRule{Address: "spam@example.com"},
			msg: &imap.Message{
				Envelope: nil,
			},
			want: false,
		},
		{
			name: "empty from addresses - should not delete",
			rule: &AddressRule{Address: "spam@example.com"},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.ShouldDelete(tt.msg)
			if got != tt.want {
				t.Errorf("AddressRule.ShouldDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddressRule(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "valid address",
			address: "test@example.com",
			wantErr: false,
		},
		{
			name:    "empty address - should error",
			address: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAddressRule(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAddressRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Address != tt.address {
				t.Errorf("NewAddressRule() address = %v, want %v", got.Address, tt.address)
			}
		})
	}
}
