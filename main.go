package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	logbookDir   = flag.String("dir", os.Getenv("HOME")+"/logbook", "The directory logbook files are stored")
	summaryMatch = regexp.MustCompile(`- SUMMARY\:(?:(.+)\:)? (.+)`)
)

func main() {
	flag.Parse()
	week, err := time.Parse("2006-01-02", flag.Arg(0))
	if err != nil {
		exit(err)
	}
	if week.Weekday() != time.Monday {
		exit("date must be a Monday")
	}

	logfiles := getLogfiles(week)

	summary := make(map[string][]string)
	// TODO: parallelise?
	for _, logfile := range logfiles {
		daySummary, err := getSummaries(logfile)
		if err != nil {
			exit(err)
		}
		for k, v := range daySummary {
			summary[k] = append(summary[k], v...)
		}
	}

	printSummary(week, summary)
}

func getLogfiles(start time.Time) (logfiles []string) {
	for i := 0; i < 7; i++ {
		day := start.Add(time.Hour * 24 * time.Duration(i))
		logfile := fmt.Sprintf("%s.md", day.Format("2006-01-02"))
		logfiles = append(logfiles, logfile)
	}
	return
}

func getSummaries(logfile string) (map[string][]string, error) {
	summaries := make(map[string][]string)
	file, err := filepath.Abs(filepath.Join(*logbookDir, logfile))
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	matches := summaryMatch.FindAllSubmatch(content, -1)

	for _, match := range matches {
		topic := string(match[1])
		topic = strings.TrimSpace(topic)
		if topic == "" {
			topic = "Misc"
		}
		summary := string(match[2])
		summary = strings.TrimSpace(summary)
		summaries[topic] = append(summaries[topic], summary)
	}

	return summaries, nil
}

func printSummary(week time.Time, summary map[string][]string) {
	fmt.Printf("James Routley Weekly Update w/c %s\n\n", week.Format("2006-01-02"))
	var topics []string
	for topic := range summary {
		// We wish summaries without a topic to be printed last
		if topic == "Misc" {
			continue
		}
		topics = append(topics, topic)
	}

	sort.Strings(topics)
	topics = append(topics, "Misc")
	for _, topic := range topics {
		fmt.Println(topic)
		for _, s := range summary[topic] {
			fmt.Printf("- %s\n", s)
		}
		fmt.Println("")
	}
}

func exit(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
