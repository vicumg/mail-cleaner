package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SpamStats struct {
	TotalEmails int
	ByDomain    map[string]int
	ByAddress   map[string]int
	Subjects    []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: analyze <log-file>")
		fmt.Println("Example: analyze spam_classification.log")
		os.Exit(1)
	}

	logFile := os.Args[1]
	stats, err := parseLogFile(logFile)
	if err != nil {
		fmt.Printf("Error parsing log file: %v\n", err)
		os.Exit(1)
	}

	printStats(stats)
}

func parseLogFile(filename string) (*SpamStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &SpamStats{
		ByDomain:  make(map[string]int),
		ByAddress: make(map[string]int),
		Subjects:  make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "Classified as spam:") {
			continue
		}

		stats.TotalEmails++

		parts := strings.SplitN(line, "Classified as spam:", 2)
		if len(parts) < 2 {
			continue
		}

		content := strings.TrimSpace(parts[1])
		emailSubject := strings.SplitN(content, " - ", 2)
		if len(emailSubject) < 2 {
			continue
		}

		email := strings.TrimSpace(emailSubject[0])
		subject := strings.TrimSpace(emailSubject[1])

		if atPos := strings.Index(email, "@"); atPos != -1 {
			domain := email[atPos+1:]
			stats.ByDomain[domain]++
		}

		stats.ByAddress[email]++
		stats.Subjects = append(stats.Subjects, subject)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func printStats(stats *SpamStats) {
	fmt.Println("=== Spam Classification Statistics ===")
	fmt.Printf("\nTotal spam emails detected: %d\n", stats.TotalEmails)

	if stats.TotalEmails == 0 {
		fmt.Println("\nNo spam emails found in the log file.")
		return
	}

	fmt.Println("\n=== Top Spam Domains ===")
	topDomains := getTopN(stats.ByDomain, 10)
	for i, item := range topDomains {
		fmt.Printf("%2d. %s: %d emails (%.1f%%)\n",
			i+1, item.Key, item.Count,
			float64(item.Count)*100/float64(stats.TotalEmails))
	}

	fmt.Println("\n=== Top Spam Addresses ===")
	topAddresses := getTopN(stats.ByAddress, 10)
	for i, item := range topAddresses {
		fmt.Printf("%2d. %s: %d emails (%.1f%%)\n",
			i+1, item.Key, item.Count,
			float64(item.Count)*100/float64(stats.TotalEmails))
	}

	fmt.Println("\n=== Sample Spam Subjects (first 10) ===")
	limit := 10
	if len(stats.Subjects) < limit {
		limit = len(stats.Subjects)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("%2d. %s\n", i+1, stats.Subjects[i])
	}
}

type KeyValue struct {
	Key   string
	Count int
}

func getTopN(m map[string]int, n int) []KeyValue {
	items := make([]KeyValue, 0, len(m))
	for k, v := range m {
		items = append(items, KeyValue{Key: k, Count: v})
	}

	for i := 0; i < len(items) && i < n; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Count > items[i].Count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	if len(items) > n {
		items = items[:n]
	}

	return items
}
