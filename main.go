package main

import (
	"fmt"
	"mail-cleaner/internal/config"
	"mail-cleaner/internal/imap"
	"mail-cleaner/internal/rules"
	"os"
)

func main() {
	//get service name from input arguments
	if len(os.Args) < 3 {
		panic("Service name or rule set file argument is missing")
	}

	service_name := os.Args[1]
	fmt.Printf("Loading config for service: %s\n", service_name)
	cfg := config.LoadConfig(service_name)
	fmt.Println(cfg)

	rule_set_file := os.Args[2]
	rules := rules.NewRules(rule_set_file)
	fmt.Printf("Loading rules from file: %s\n", rule_set_file)

	// Here you can proceed to create an IMAP client and connect using the loaded config
	imapClient := imap.NewClient(cfg)
	if err := imapClient.Connect(); err != nil {
		fmt.Printf("Error connecting to IMAP server: %v\n", err)
		return
	}
	defer imapClient.Disconnect()

	imapClient.CleanEmails(rules)

}
