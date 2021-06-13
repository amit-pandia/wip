package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// lineContents holds the values parsed from each line
type lineContents struct {
	time          time.Time
	user          string
	sessionStatus string
}

// sessionDetails holds the details of a session
type sessionDetails struct {
	username      string
	startTimes    []time.Time
	endTimes      []time.Time
	totalTime     time.Duration
	totalSessions int
}

func main() {
	// parse arguments
	if err := validateArgs(os.Args); err != nil {
		fmt.Println(err)
		return
	}

	// generate report
	report, err := generateReport(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// print report
	for _, r := range report {
		fmt.Println(r)
	}
}

// validateArgs validates the command line arguments
func validateArgs(args []string) (err error) {
	argCount := len(os.Args)
	if argCount > 2 {
		err = errors.New("too many arguments")
		return
	}
	if argCount < 2 {
		err = errors.New("file path missing")
		return
	}
	return
}

// generateReport generates the report
func generateReport(filePath string) (report []string, err error) {
	var earliestTime, latestTime time.Time
	var sessions map[string]sessionDetails
	sessions = make(map[string]sessionDetails, 0)

	//ReadFile
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// getLines from file
	lines := getLines(fileData)

	// fetchEarliestTime
	earliestTime, err = fetchEarliestTime(lines)
	if err != nil {
		return nil, err
	}

	// fetchLatestTime
	latestTime, err = fetchLatestTime(lines)
	if err != nil {
		return
	}

	// session wise pairing of start times and end times
	for _, l := range lines {
		lineContents, err := parseLine(l)
		if err != nil {
			continue
		}
		if _, ok := sessions[lineContents.user]; ok {
			// existing user
			if lineContents.sessionStatus == "Start" {
				sessionDetails := sessions[lineContents.user]
				sessionDetails.startTimes = append(sessionDetails.startTimes, lineContents.time)
				sessions[lineContents.user] = sessionDetails
			}
			if lineContents.sessionStatus == "End" {
				sessionDetails := sessions[lineContents.user]
				sessionDetails.endTimes = append(sessionDetails.endTimes, lineContents.time)
				sessions[lineContents.user] = sessionDetails
			}
		} else {
			// new user detected
			if lineContents.sessionStatus == "Start" {
				var sd sessionDetails
				sd.username = lineContents.user
				sd.startTimes = append(sd.startTimes, lineContents.time)
				sessions[lineContents.user] = sd
			}
			if lineContents.sessionStatus == "End" {
				var sd sessionDetails
				sd.username = lineContents.user
				sd.startTimes = append(sd.startTimes, earliestTime)
				sd.endTimes = append(sd.endTimes, lineContents.time)
				sessions[lineContents.user] = sd
			}
		}
	}

	// calculate total time, calculate total sessions and prepare report
	for _, s := range sessions {
		strtCount := len(s.startTimes)
		endCount := len(s.endTimes)
		if strtCount > endCount {
			for i := 0; i < strtCount; i++ {
				if i < endCount {
					s.totalTime += s.endTimes[i].Sub(s.startTimes[i])
				} else {
					s.totalTime += latestTime.Sub(s.startTimes[i])
				}
			}
			s.totalSessions = strtCount
		} else {
			for i := 0; i < endCount; i++ {
				if i < strtCount {
					s.totalTime += s.endTimes[i].Sub(s.startTimes[i])
				} else {
					s.totalTime += s.endTimes[i].Sub(earliestTime)
				}
			}
			s.totalSessions = endCount
		}
		sessions[s.username] = s
		report = append(report, fmt.Sprintf("%v %v %v", s.username, s.totalSessions, s.totalTime.Seconds()))
	}
	return
}

// parseLine extracts time, username and session status from every line
func parseLine(line string) (lineContents lineContents, err error) {
	words := strings.Fields(line)
	if len(words) < 3 {
		err = errors.New("line parsing failed")
		return
	}
	timeStamp, err := time.Parse("15:04:05", words[0])
	if err != nil {
		return
	}
	lineContents.time = timeStamp
	lineContents.user = words[1]
	lineContents.sessionStatus = words[2]
	return
}

// fetchEarliestTime fetches the earliest time
func fetchEarliestTime(lines []string) (t time.Time, err error) {
	for i := 0; i <= len(lines)-1; i++ {
		lineContents, err := parseLine(lines[i])
		if err != nil {
			if i < (len(lines) - 1) {
				continue
			} else {
				return t, err
			}
		}
		t = lineContents.time
		return t, err
	}
	return
}

// fetchLatestTime fetches the latest time
func fetchLatestTime(lines []string) (t time.Time, err error) {
	for i := len(lines) - 1; i >= 0; i-- {
		lineContents, err := parseLine(lines[i])
		if err != nil {
			if i > 0 {
				continue
			} else {
				return t, err
			}
		}
		t = lineContents.time
		return t, err
	}
	return
}

// getLines returns lines of file in form of string array
func getLines(fileData []byte) []string {
	return strings.Split(strings.TrimSpace(string(fileData)), "\n")
}
