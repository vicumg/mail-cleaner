package rule

import (
	"fmt"
	"os"

	"mail-cleaner/internal/ai/ollama"
	"mail-cleaner/internal/rules"

	"github.com/emersion/go-imap"
)

type Classifier interface {
	IsSpam(emailAddress string, subject string, prompt string) (bool, error)
}

type AIRule struct {
	Enabled            bool   `json:"enabled"`
	Action             string `json:"action"` // "log" или "delete"
	prompt             string
	classifier         Classifier
	excluded_domains   []string
	excluded_addresses []string
	logFile            *os.File
}

func init() {
	RegisterRuleFactory("ai_local_rule", func(config map[string]interface{}) (rules.Rule, error) {
		enabled, ok := config["enabled"].(bool)
		if !ok {
			enabled = false
		}

		action, ok := config["action"].(string)
		if !ok {
			action = "log" // по умолчанию только логирование
		}
		if action != "log" && action != "delete" {
			return nil, fmt.Errorf("action must be 'log' or 'delete', got: %s", action)
		}

		prompt, ok := config["prompt"].(string)
		if !ok {
			prompt = "Is this email spam? answer only one word:(spam or ham)."
		}

		baseURL, _ := config["host_url"].(string)
		model, _ := config["model"].(string)

		//todo: make client configurable
		client := ollama.NewClient(baseURL, model)

		excludedDomains, ok := config["exclude_domains"].([]string)
		if !ok && enabled {
			panic("Must contains at least one excluded domain")
		}

		excludedAddresses, ok := config["exclude_addresses"].([]string)
		if !ok && enabled {
			panic("Must contains at least one excluded address")
		}

		return NewAIRule(enabled, action, prompt, client, excludedDomains, excludedAddresses)
	})
}

func NewAIRule(enabled bool, action string, prompt string, classifier Classifier, excludedDomains []string, excludedAddresses []string) (*AIRule, error) {
	var logFile *os.File
	var err error

	if enabled {
		logFile, err = os.OpenFile("spam_classification.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Warning: failed to open log file: %v\n", err)
		}
	}

	return &AIRule{
		Enabled:            enabled,
		Action:             action,
		prompt:             prompt,
		classifier:         classifier,
		excluded_domains:   excludedDomains,
		excluded_addresses: excludedAddresses,
		logFile:            logFile,
	}, nil
}

func (ar *AIRule) ShouldDelete(msg *imap.Message) bool {
	if !ar.Enabled {
		return false
	}

	if msg.Envelope == nil || len(msg.Envelope.From) == 0 {
		return false
	}

	for _, addr := range msg.Envelope.From {
		emailAddress := addr.MailboxName + "@" + addr.HostName
		subject := msg.Envelope.Subject

		// Check excluded addresses
		excluded := false
		for _, excludedAddress := range ar.excluded_addresses {
			if emailAddress == excludedAddress {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		// Check excluded domains
		for _, excludedDomain := range ar.excluded_domains {
			if addr.HostName == excludedDomain {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if ar.apply(emailAddress, subject) {
			return true
		}
	}

	return false
}

func (ar *AIRule) apply(emailAddress string, subject string) bool {
	if ar.classifier == nil {
		return false
	}

	isSpam, err := ar.classifier.IsSpam(emailAddress, subject, ar.prompt)
	if err != nil {
		fmt.Printf("Error classifying email: %v\n", err)
		return false
	}

	if isSpam {
		message := fmt.Sprintf("Classified as spam: %s - %s\n", emailAddress, subject)
		if ar.logFile != nil {
			ar.logFile.WriteString(message)
		} else {
			fmt.Print(message)
		}

		// the main logic to decide whether to delete or not
		if ar.Action == "delete" {
			return true
		}
	}

	return false
}

func (ar *AIRule) Close() error {
	if ar.logFile != nil {
		return ar.logFile.Close()
	}
	return nil
}

func (ar *AIRule) String() string {
	return fmt.Sprintf("AIRule{Enabled: %v, Action: %s}", ar.Enabled, ar.Action)
}
