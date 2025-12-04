package rule

import (
	"errors"
	"testing"

	"github.com/emersion/go-imap"
)

type mockClassifier struct {
	shouldReturnSpam bool
	shouldReturnErr  bool
}

func (m *mockClassifier) IsSpam(emailAddress string, subject string, prompt string) (bool, error) {
	if m.shouldReturnErr {
		return false, errors.New("mock error")
	}
	return m.shouldReturnSpam, nil
}

func TestAIRule_apply(t *testing.T) {
	tests := []struct {
		name         string
		action       string
		classifier   Classifier
		emailAddress string
		subject      string
		want         bool
	}{
		{
			name:         "classifier returns spam - log mode",
			action:       "log",
			classifier:   &mockClassifier{shouldReturnSpam: true},
			emailAddress: "spam@example.com",
			subject:      "Get rich quick!",
			want:         false, // log mode - не удаляем
		},
		{
			name:         "classifier returns spam - delete mode",
			action:       "delete",
			classifier:   &mockClassifier{shouldReturnSpam: true},
			emailAddress: "spam@example.com",
			subject:      "Get rich quick!",
			want:         true, // delete mode - удаляем
		},
		{
			name:         "classifier returns not spam",
			action:       "delete",
			classifier:   &mockClassifier{shouldReturnSpam: false},
			emailAddress: "test@example.com",
			subject:      "Test Subject",
			want:         false,
		},
		{
			name:         "classifier returns error",
			action:       "delete",
			classifier:   &mockClassifier{shouldReturnErr: true},
			emailAddress: "test@example.com",
			subject:      "Test Subject",
			want:         false,
		},
		{
			name:         "nil classifier",
			action:       "log",
			classifier:   nil,
			emailAddress: "test@example.com",
			subject:      "Test Subject",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &AIRule{
				Enabled:    true,
				Action:     tt.action,
				prompt:     "test prompt",
				classifier: tt.classifier,
			}
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
			name: "enabled rule with spam classifier - log mode",
			rule: &AIRule{
				Enabled:    true,
				Action:     "log",
				prompt:     "test prompt",
				classifier: &mockClassifier{shouldReturnSpam: true},
			},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
					Subject: "Get rich quick!",
				},
			},
			want: false, // log mode - do not delete
		},
		{
			name: "enabled rule with spam classifier - delete mode",
			rule: &AIRule{
				Enabled:    true,
				Action:     "delete",
				prompt:     "test prompt",
				classifier: &mockClassifier{shouldReturnSpam: true},
			},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "spam", HostName: "example.com"},
					},
					Subject: "Get rich quick!",
				},
			},
			want: true, // delete mode - delete
		},
		{
			name: "enabled rule with not spam classifier",
			rule: &AIRule{
				Enabled:    true,
				prompt:     "test prompt",
				classifier: &mockClassifier{shouldReturnSpam: false},
			},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From: []*imap.Address{
						{MailboxName: "test", HostName: "example.com"},
					},
					Subject: "Normal email",
				},
			},
			want: false,
		},
		{
			name: "disabled rule",
			rule: &AIRule{
				Enabled:    false,
				classifier: &mockClassifier{shouldReturnSpam: true},
			},
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
			rule: &AIRule{
				Enabled:    true,
				classifier: &mockClassifier{shouldReturnSpam: true},
			},
			msg: &imap.Message{
				Envelope: nil,
			},
			want: false,
		},
		{
			name: "empty from addresses",
			rule: &AIRule{
				Enabled:    true,
				classifier: &mockClassifier{shouldReturnSpam: true},
			},
			msg: &imap.Message{
				Envelope: &imap.Envelope{
					From:    []*imap.Address{},
					Subject: "Test",
				},
			},
			want: false,
		},
		{
			name: "nil classifier",
			rule: &AIRule{
				Enabled:    true,
				classifier: nil,
			},
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
		name       string
		enabled    bool
		action     string
		prompt     string
		classifier *mockClassifier
		wantErr    bool
	}{
		{
			name:       "enabled with classifier - log mode",
			enabled:    true,
			action:     "log",
			prompt:     "test prompt",
			classifier: &mockClassifier{},
			wantErr:    false,
		},
		{
			name:       "enabled with classifier - delete mode",
			enabled:    true,
			action:     "delete",
			prompt:     "test prompt",
			classifier: &mockClassifier{},
			wantErr:    false,
		},
		{
			name:       "disabled with classifier",
			enabled:    false,
			action:     "log",
			prompt:     "test prompt",
			classifier: &mockClassifier{},
			wantErr:    false,
		},
		{
			name:       "enabled without classifier",
			enabled:    true,
			action:     "log",
			prompt:     "test prompt",
			classifier: nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAIRule(tt.enabled, tt.action, tt.prompt, tt.classifier, []string{}, []string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAIRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Enabled != tt.enabled {
					t.Errorf("NewAIRule() enabled = %v, want %v", got.Enabled, tt.enabled)
				}
				if got.Action != tt.action {
					t.Errorf("NewAIRule() action = %v, want %v", got.Action, tt.action)
				}
			}
		})
	}
}
