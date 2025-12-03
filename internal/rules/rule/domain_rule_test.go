package rule

import "testing"

func Test_apply(t *testing.T) {
	tests := []struct {
		name        string
		emailDomain string
		ruleDomain  string
		want        bool
	}{
		{
			name:        "exact match",
			emailDomain: "example.com",
			ruleDomain:  "example.com",
			want:        true,
		},
		{
			name:        "subdomain match",
			emailDomain: "promo.epicentrk.ua",
			ruleDomain:  "epicentrk.ua",
			want:        true,
		},
		{
			name:        "partial domain match",
			emailDomain: "newsletter.com",
			ruleDomain:  "newsletter",
			want:        true,
		},
		{
			name:        "case insensitive match",
			emailDomain: "Example.COM",
			ruleDomain:  "example.com",
			want:        true,
		},
		{
			name:        "no match different domain",
			emailDomain: "example.com",
			ruleDomain:  "test.com",
			want:        false,
		},
		{
			name:        "no match partial overlap",
			emailDomain: "example.com",
			ruleDomain:  "ample.co",
			want:        true,
		},
		{
			name:        "empty email domain",
			emailDomain: "",
			ruleDomain:  "example.com",
			want:        false,
		},
		{
			name:        "empty rule domain",
			emailDomain: "example.com",
			ruleDomain:  "",
			want:        true,
		},
		{
			name:        "both empty",
			emailDomain: "",
			ruleDomain:  "",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&DomainRule{}).apply(tt.emailDomain, tt.ruleDomain)
			if got != tt.want {
				t.Errorf("apply(%q, %q) = %v, want %v",
					tt.emailDomain, tt.ruleDomain, got, tt.want)
			}
		})
	}
}
