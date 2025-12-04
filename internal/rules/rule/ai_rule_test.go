package rule

import (
	"testing"

	"github.com/emersion/go-imap"
)

func TestAIRule_apply(t *testing.T) {
	tests := []struct {
		name         string
		emailAddress string
		subject      string
		want         bool
	}{
		{
			name:         "any email and subject",
			emailAddress: "test@example.com",
			subject:      "Test Subject",
			want:         false,
		},
		{
			name:         "empty email",
			emailAddress: "",
			subject:      "Test Subject",
			want:         false,
		},
		{
			name:         "empty subject",
			emailAddress: "test@example.com",
			subject:      "",
			want:         false,
		},
		{
			name:         "both empty",
			emailAddress: "",
			subject:      "",
			want:         false,
		},
	}

	rule := &AIRule{Enabled: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rule.apply(tt.emailAddress, tt.subject)
			if got != tt.want {
				t.Errorf("AIRule.apply(%q, %q) = %v, want %v",
					tt.emailAddress, tt.subject, got, tt.want)
			}
		})
	}
}

func TestAIRule_ShouldDelete(t *testing.T) {
	tests := []struct {
		name string
		rule *AIRule
		msg  *imap.Message
		want bool
	}{
		{
			name: "enabled rule - returns false for now",
			rule: &AIRule{Enabled: true},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
					Subject: "Get rich quick!",
				},
			},
			want: false,
		},
		{
			name: "disabled rule",
			rule: &AIRule{Enabled: false},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
					Subject: "Get rich quick!",
				},
			},
			want: false,
		},
		{
			name: "nil envelope",
			rule: &AIRule{Enabled: true},
			msg: &imap.Message{
				Envelope: nil,
			},
			want: false,
		},
		{
			name: "empty from addresses",
			rule: &AIRule{Enabled: true},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From:    []*imap.Address{},
					Subject: "Test",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.ShouldDelete(tt.msg)
			if got != tt.want {
				t.Errorf("AIRule.ShouldDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAIRule(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		wantErr bool
	}{
		{
			name:    "enabled",
			enabled: true,
			wantErr: false,
		},
		{
			name:    "disabled",
			enabled: false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAIRule(tt.enabled)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAIRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Enabled != tt.enabled {
				t.Errorf("NewAIRule() enabled = %v, want %v", got.Enabled, tt.enabled)
			}
		})
	}
}
