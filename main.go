package main

import (
	"fmt"
	"mail-cleaner/internal/config"
	"mail-cleaner/internal/imap"
	"mail-cleaner/internal/rules"
	"mail-cleaner/internal/rules/rule"
	"os"
)

func main() {
	//get service name from input arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: mail-cleaner <service_name> <rule_set_file>")
		os.Exit(1)
	}

	service_name := os.Args[1]
	fmt.Printf("Loading config for service: %s\n", service_name)
	cfg := config.LoadConfig(service_name)
	fmt.Println(cfg)

	rule_set_file := os.Args[2]

	rules_list, err := rule.CreateFromFile(rule_set_file)
	if err != nil {
		fmt.Printf("Failed to create rules from file: %v\n", err)
		os.Exit(1)
	}

	imapClient := imap.NewClient(cfg)
	if err := imapClient.Connect(); err != nil {
		fmt.Printf("Error connecting to IMAP server: %v\n", err)
		return
	}
	defer imapClient.Disconnect()

	imapClient.CleanEmails(rules.NewRules(rules_list))

}
