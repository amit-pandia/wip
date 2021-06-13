package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	line := "14:04:05 ALICE99 End"
	lineContents, err := parseLine(line)
	if err != nil {
		t.Error(err)
	}
	if !(lineContents.user == "ALICE99" && lineContents.sessionStatus == "End" && lineContents.time.Format("15:04:05") == "14:04:05") {
		t.Error("mismatched line contents")
	}
}

func TestFetchLatestTime(t *testing.T) {
	lines := []string{
		"14:04:05 ALICE99 End",
		"14:04:23 ALICE99 End",
		"14:04:41 CHARLIE Start"}
	fetchedTime, err := fetchLatestTime(lines)
	if err != nil {
		t.Error(err)
	}
	if fetchedTime.Format("15:04:05") != "14:04:41" {
		t.Error("mismatched time")
	}
}

func TestFetchEarliestTime(t *testing.T) {
	lines := []string{
		"14:04:05 ALICE99 End",
		"14:04:23 ALICE99 End",
		"14:04:41 CHARLIE Start"}
	fetchedTime, err := fetchEarliestTime(lines)
	if err != nil {
		t.Error(err)
	}
	if fetchedTime.Format("15:04:05") != "14:04:05" {
		t.Error("mismatched time")
	}
}

func TestGetLines(t *testing.T) {
	fileData := []byte(`14:04:05 ALICE99 End
14:04:23 ALICE99 End
14:04:41 CHARLIE Start
`)
	lines := getLines(fileData)
	if !(lines[0] == "14:04:05 ALICE99 End" && lines[1] == "14:04:23 ALICE99 End" && lines[2] == "14:04:41 CHARLIE Start") {
		t.Error("mismatched lines")
	}
}

func TestGenerateReport(t *testing.T) {
	// test 1:
	// 14:02:03 ALICE99 Start
	// 14:02:05 CHARLIE End
	// 14:02:34 ALICE99 End
	// 14:02:58 ALICE99 Start
	// 14:03:02 CHARLIE Start
	// 14:03:33 ALICE99 Start
	// 14:03:35 ALICE99 End
	// 14:03:37 CHARLIE End
	// 14:04:05 ALICE99 End
	// 14:04:23 ALICE99 End
	// 14:04:41 CHARLIE Start
	report, err := generateReport("./samplelogfiles/log1.txt")
	if err != nil {
		t.Error(err)
	}
	if !(report[0] == "ALICE99 4 240" && report[1] == "CHARLIE 3 37" || report[0] == "CHARLIE 3 37" && report[1] == "ALICE99 4 240") {
		t.Error("incorrect report for test 1")
	}

	// test 2:
	// 14:02:03 ALICE99 Start
	// 14:02:05 CHARLIE End
	// 14:02:34 ALICE99 End
	report, err = generateReport("./samplelogfiles/log2.txt")
	if err != nil {
		t.Error(err)
	}
	if !(report[0] == "ALICE99 1 31" && report[1] == "CHARLIE 1 2" || report[0] == "CHARLIE 1 2" && report[1] == "ALICE99 1 31") {
		t.Error("incorrect report for test 2")
	}

	// test 3:
	// 14:04:05 ALICE99 End
	// 14:04:23 ALICE99 End
	// 14:04:41 CHARLIE Start
	report, err = generateReport("./samplelogfiles/log3.txt")
	if err != nil {
		t.Error(err)
	}
	if !(report[0] == "ALICE99 2 18" && report[1] == "CHARLIE 1 0" || report[0] == "CHARLIE 1 0" && report[1] == "ALICE99 2 18") {
		t.Error("incorrect report for test 3")
	}
}
