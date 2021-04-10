package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)


func main() {
	log.Printf("executing fair billing...")

	var configFile = flag.String("input", "", "configuration file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	result, err := ReadFile(configFile)
	if err != nil {
		log.Fatalf("error: reading file - %s", err.Error())
	}

	output := ComputeUserSessions(result)
	for _, val := range output {
		log.Printf("%s %d %f", val.Name, val.ActiveSessions, val.TimeSpent)
	}

	log.Printf("done...")
}

type SessionOutput struct {
	Name           string
	ActiveSessions int
	TimeSpent      float64
}

//validates the log records
//Ideal record would be - 14:02:03 ALICE99 Start
func ValidateRecords(rec string) bool {
	//1. validate it has three split time name and session-info
	recArr := strings.Split(rec, " ")
	if len(recArr) != 3 {
		return false
	}

	//2. validate the time it has hr:mm:ss
	if ok, err := regexp.MatchString("^(?:(?:([01]?\\d|2[0-3]):)?([0-5]?\\d):)?([0-5]?\\d)$", recArr[0]); !ok || err != nil {
		log.Printf("error: regex - %v", err)
		return false
	}

	//3. Name should not be empty
	if len(strings.TrimSpace(recArr[1])) < 1 {
		log.Printf("error: name - %s", recArr[1])
		return false
	}

	//4. session should be available with start / end values
	if !strings.EqualFold(recArr[2], "Start") && !strings.EqualFold(recArr[2], "End") {
		log.Printf("error: session - %s", recArr[2])
		return false
	}

	return true
}

func ReadFile(filePath *string) ([]string, error) {
	file, err := os.Open(*filePath)
	if err != nil {
		log.Printf("error: ReadFile - %s", err.Error())
		return nil, err
	}

	defer file.Close()

	var lines []string
	scIns := bufio.NewScanner(file)

	for scIns.Scan() {
		lines = append(lines, scIns.Text())
	}

	return lines, nil
}

type SessionInfo struct {
	StartTime *string
	EndTime   *string
	TimeSpent float64
	IsComplete bool
}

func (s *SessionInfo) ComputeTimeSpent() {

	sTime, err := time.Parse("15:04:05", *s.StartTime)
	if err != nil {
		log.Printf("error: parse time error (start time) - %s", err.Error())
		return
	}

	eTime, err := time.Parse("15:04:05", *s.EndTime)
	if err != nil {
		log.Printf("error: parse time error (start time) - %s", err.Error())
		return
	}

	s.TimeSpent = eTime.Sub(sTime).Seconds()
}


func ComputeUserSessions(logs []string) map[string]SessionOutput {
	var (
		earliestTime  *string
		latestTime    *string
		userSessionCache  map[string][]*SessionInfo
	)

	earliestTime = GetEarliestTime(logs)
	latestTime = GetLatestTime(logs)

	userSessionCache = make(map[string][]*SessionInfo)

	for l := range logs {
		//validate the records are perfect without any missing fields
		if !ValidateRecords(logs[l]) {
			continue
		}
		logInsArr := strings.Split(logs[l], " ")

		sessionInfoIns := new(SessionInfo)
		if val, ok := userSessionCache[logInsArr[1]]; ok {
			newComputedSession := MapExistingSessionsForUser(val, logInsArr, earliestTime)
			userSessionCache[logInsArr[1]] = newComputedSession
			continue
		}

		//if record for the user is not found and its session start
		if logInsArr[2] == "Start" {
			sessionInfoIns.StartTime = &logInsArr[0]
			sessionInfoIns.IsComplete = false
			userSessionCache[logInsArr[1]] = []*SessionInfo{sessionInfoIns}
			continue
		}

		//if the record for the user is not found and the session is end
		sessionInfoIns.StartTime = earliestTime
		sessionInfoIns.EndTime = &logInsArr[0]
		sessionInfoIns.ComputeTimeSpent()
		sessionInfoIns.IsComplete = true
		userSessionCache[logInsArr[1]] = []*SessionInfo{sessionInfoIns}
	}

	output := make(map[string]SessionOutput)
	for key, val := range userSessionCache {
		var totalTimeSpent float64
		for i := range val {
			if val[i].EndTime == nil {
				val[i].EndTime = latestTime
				val[i].ComputeTimeSpent()
			}

			totalTimeSpent += val[i].TimeSpent
		}

		output[key] = SessionOutput{
			Name: key,
			ActiveSessions: len(val),
			TimeSpent: totalTimeSpent,
		}
	}

	return output
}

//earliest time is the first record of the file
func GetEarliestTime(logs []string) *string {
	if len(logs) < 1 {
		return nil
	}

	logInsArr := strings.Split(logs[0], " ")
	return &logInsArr[0]
}

//latest time would be the last record of the log file
func GetLatestTime(logs []string) *string {
	if len(logs) < 1 {
		return nil
	}

	logInsArr := strings.Split(logs[len(logs) -1], " ")
	return &logInsArr[0]
}

func MapExistingSessionsForUser(sessions []*SessionInfo, log []string, earliestTime *string) []*SessionInfo {
	for i := range sessions {
		if sessions[i].IsComplete {
			continue
		}

		if !sessions[i].IsComplete {
			if log[2] == "End" {
				sessions[i].IsComplete = true
				sessions[i].EndTime = &log[0]
				sessions[i].ComputeTimeSpent()
				return sessions
			}

			sessions = append(sessions, &SessionInfo{
				StartTime: &log[0],
				IsComplete: false,
			})
			return sessions
		}
	}

	if log[2] == "End" {
		sessionIns := &SessionInfo{
			StartTime:  earliestTime,
			IsComplete: true,
			EndTime:    &log[0],
		}
		sessionIns.ComputeTimeSpent()
		sessions = append(sessions, sessionIns)
		return sessions
	}

	sessionIns := &SessionInfo{
		StartTime:  &log[0],
		IsComplete: false,
	}
	sessions = append(sessions, sessionIns)
	return sessions
}
