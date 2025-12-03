package rule

import (
	"encoding/json"
	"fmt"
	"mail-cleaner/internal/rules"
	"os"
)

type RuleFactory func(data map[string]any) (rules.Rule, error)

var factories = make(map[string]RuleFactory)

func RegisterRuleFactory(ruleType string, factory RuleFactory) {
	factories[ruleType] = factory
}

func CreateFromFile(rule_set_file string) ([]rules.Rule, error) {
	// load data from json file and create rules
	file_data, err := os.ReadFile(rule_set_file)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var raw_rules []map[string]any
	if err := json.Unmarshal(file_data, &raw_rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules file: %w", err)
	}

	var rulesList []rules.Rule
	for _, raw_rule := range raw_rules {
		ruleType, ok := raw_rule["type"].(string)
		if !ok {
			fmt.Println("Invalid rule type in rules file")
			continue
		}

		factory, exists := factories[ruleType]
		if !exists {
			fmt.Printf("No factory registered for rule type: %s\n", ruleType)
			continue
		}

		rule, err := factory(raw_rule)
		if err != nil {
			fmt.Printf("Error creating rule of type %s: %v\n", ruleType, err)
			continue
		}

		rulesList = append(rulesList, rule)
	}

	if len(rulesList) == 0 {
		fmt.Println("No valid rules found in the rules file.")
	}

	return rulesList, nil
}
