package summary

import (
	"bytes"
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
	logbookDir   = os.Getenv("HOME") + "/logbook"
	summaryMatch = regexp.MustCompile(`- SUMMARY\:(?:(.+)\:)? (.+)`)
)

type Summary map[string][]string

func Summarise(start, end time.Time) *Summary {
	logfiles := getLogfiles(start, end)
	summary := make(Summary)
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
	return &summary
}

func (s Summary) String() string {
	var b bytes.Buffer
	topics := s.sortedKeys()
	for _, topic := range topics {
		fmt.Fprintln(&b, topic)
		for _, summary := range s[topic] {
			fmt.Fprintf(&b, "- %s\n", summary)
		}
		fmt.Fprintln(&b, "")
	}
	return strings.TrimRight(b.String(), "\n")
}

func (s *Summary) sortedKeys() (keys []string) {
	for k := range *s {
		if k == "Misc" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keys = append(keys, "Misc")
	return
}

func sortedKeys(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func getLogfiles(start, end time.Time) (logfiles []string) {
	days := daysInRange(start, end)
	for _, day := range days {
		logfile := fmt.Sprintf("%s.md", day.Format("2006-01-02"))
		logfiles = append(logfiles, logfile)
	}
	return
}

func getSummaries(logfile string) (map[string][]string, error) {
	summaries := make(map[string][]string)
	file, err := filepath.Abs(filepath.Join(logbookDir, logfile))
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

func daysInRange(start, end time.Time) (days []time.Time) {
	t := start.Truncate(time.Hour * 24)
	end = end.Truncate(time.Hour * 24).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	for t.Before(end) {
		days = append(days, t)
		t = t.AddDate(0, 0, 1)
	}
	return
}
