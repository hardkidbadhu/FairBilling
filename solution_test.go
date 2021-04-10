package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFairBilling(t *testing.T) {
	//1. Tests for Read file
	t.Run("TestReadFileShouldPassForFileExists", TestReadFileShouldPassForFileExists)
	t.Run("TestReadFileShouldFailForFileDoesNotExists", TestReadFileShouldFailForFileDoesNotExists)

	//2. Tests for Validate Records
	t.Run("TestValidateRecordsShouldReturnTrue", TestValidateRecordsShouldReturnTrue)
	t.Run("TestValidateRecordsShouldReturnFalse", TestValidateRecordsShouldReturnFalse)
	t.Run("TestValidateRecordsShouldReturnFalseForInvalidName", TestValidateRecordsShouldReturnFalseForInvalidName)
	t.Run("TestValidateRecordsShouldReturnFalseForInvalidSession", TestValidateRecordsShouldReturnFalseForInvalidSession)
	t.Run("TestValidateRecordsShouldReturnFalseForInvalidSessionInput", TestValidateRecordsShouldReturnFalseForInvalidSessionInput)

	//3. Tests for getEarliestTime
	t.Run("TestGetEarliestTime", TestGetEarliestTime)
	t.Run("TestGetEarliestTimeShouldReturnNil", TestGetEarliestTimeShouldReturnNil)

	//4. Tests for getLatestTime
	t.Run("TestGetLatestTime", TestGetLatestTime)
	t.Run("TestGetLatestTimeShouldReturnNil", TestGetLatestTimeShouldReturnNil)

	//5. Tests for ComputeTimeSpent
	t.Run("TestComputeTimeSpentForValidTime", TestComputeTimeSpentForValidTime)
	t.Run("TestComputeTimeSpentForInValidTime", TestComputeTimeSpentForInValidTime)
	t.Run("TestComputeTimeSpentForInValidEndTime", TestComputeTimeSpentForInValidEndTime)

	//6. Tests fot ComputeUserSessions
	t.Run("TestMapExistingSessionsForUserCase1", TestMapExistingSessionsForUserCase1)

	//master tests
	t.Run("TestComputeUserSessions", TestComputeUserSessions)
	t.Run("TestComputeUserSessionsCase2", TestComputeUserSessionsCase2)
}

var expectedOut = []string{
	"14:02:03 ALICE99 Start",
	"14:02:05 CHARLIE End",
	"14:02:34 ALICE99 End",
	"14:02:58 ALICE99 Start",
	"14:03:02 CHARLIE Start",
	"14:03:33 ALICE99 Start",
	"14:03:35 ALICE99 End",
	"14:03:37 CHARLIE End",
	"14:04:05 ALICE99 End",
	"14:04:23 ALICE99 End",
	"14:04:41 CHARLIE Start",
}

func TestReadFileShouldPassForFileExists(tf *testing.T) {
	filePath := "input.txt"
	output, err := ReadFile(&filePath)
	if err != nil {
		assert.Fail(tf, err.Error())
		return
	}

	expectedOut := []string{
		"14:02:03 ALICE99 Start",
		"14:02:05 CHARLIE End",
		"14:02:34 ALICE99 End",
		"14:02:58 ALICE99 Start",
		"14:03:02 CHARLIE Start",
		"14:03:33 ALICE99 Start",
		"14:03:35 ALICE99 End",
		"14:03:37 CHARLIE End",
		"14:04:05 ALICE99 End",
		"14:04:23 ALICE99 End",
		"14:04:41 CHARLIE Start",
	}

	assert.Equal(tf, expectedOut, output)
}

func TestReadFileShouldFailForFileDoesNotExists(tf *testing.T) {
	filePath := "inut.txt"
	output, err := ReadFile(&filePath)
	assert.NotNil(tf, err, err.Error())
	assert.Nil(tf, output)
}

func TestValidateRecordsShouldReturnTrue(tf *testing.T) {
	ok := ValidateRecords("14:02:03 ALICE99 Start")
	assert.Equal(tf, true, ok, "valid record")
}

func TestValidateRecordsShouldReturnFalse(tf *testing.T) {
	ok := ValidateRecords("10::03 ALICE99 Start")
	assert.Equal(tf, false, ok, "invalid record")
}

func TestValidateRecordsShouldReturnFalseForInvalidSessionInput(tf *testing.T) {
	ok := ValidateRecords("10::03 ALICE99 Sta")
	assert.Equal(tf, false, ok, "invalid record")
}

func TestValidateRecordsShouldReturnFalseForInvalidName(tf *testing.T) {
	ok := ValidateRecords("10::03 Start")
	assert.Equal(tf, false, ok, "invalid record")
}

func TestValidateRecordsShouldReturnFalseForInvalidSession(tf *testing.T) {
	ok := ValidateRecords("10::03 ALice")
	assert.Equal(tf, false, ok, "invalid record")
}

func TestGetEarliestTime(te *testing.T) {
	out := GetEarliestTime(expectedOut)
	assert.Equal(te, *out, "14:02:03")
}

func TestGetEarliestTimeShouldReturnNil(te *testing.T) {
	out := GetEarliestTime([]string{})
	assert.Nil(te, out)
}

func TestGetLatestTime(te *testing.T) {
	out := GetLatestTime(expectedOut)
	assert.Equal(te, *out, "14:04:41")
}

func TestGetLatestTimeShouldReturnNil(te *testing.T) {
	out := GetLatestTime([]string{})
	assert.Nil(te, out)
}

func TestComputeTimeSpentForValidTime(ts *testing.T) {
	sessionInfo := SessionInfo{
		StartTime: &[]string{"14:02:03"}[0],
		EndTime:   &[]string{"14:02:05"}[0],
	}

	sessionInfo.ComputeTimeSpent()
	assert.Equal(ts, sessionInfo.TimeSpent, 2.0)
}

func TestComputeTimeSpentForInValidTime(ts *testing.T) {
	sessionInfo := SessionInfo{
		StartTime: &[]string{"14::03"}[0],
		EndTime:   &[]string{"14:02:05"}[0],
	}

	sessionInfo.ComputeTimeSpent()
	assert.Equal(ts, sessionInfo.TimeSpent, 0.0)
}

func TestComputeTimeSpentForInValidEndTime(ts *testing.T) {
	sessionInfo := SessionInfo{
		StartTime: &[]string{"14:02:03"}[0],
		EndTime:   &[]string{"14::05"}[0],
	}

	sessionInfo.ComputeTimeSpent()
	assert.Equal(ts, sessionInfo.TimeSpent, 0.0)
}

func TestMapExistingSessionsForUserCase1(ms *testing.T) {
	expectedSessionInfo := []*SessionInfo{
		{StartTime: &[]string{"14:02:03"}[0], EndTime: &[]string{"14:02:05"}[0], TimeSpent: 2, IsComplete: true},
	}

	sessionInfo := []*SessionInfo{
		{StartTime: &[]string{"14:02:03"}[0], EndTime: nil, IsComplete: false},
	}

	logs := []string{"14:02:05", "CHARLIE", "End"}
	mappedInfo := MapExistingSessionsForUser(sessionInfo, logs, &[]string{"14:02:03"}[0])
	assert.Equal(ms, expectedSessionInfo, mappedInfo)
}

func TestMapExistingSessionsForUserCase2(ms *testing.T) {
	expectedSessionInfo := []*SessionInfo{
		{StartTime: &[]string{"14:02:03"}[0], EndTime: &[]string{"14:02:05"}[0], TimeSpent: 2, IsComplete: true},
		{StartTime: &[]string{"14:04:03"}[0], EndTime: nil, IsComplete: false},
	}

	sessionInfo := []*SessionInfo{
		{StartTime: &[]string{"14:02:03"}[0], EndTime: &[]string{"14:02:05"}[0], TimeSpent: 2, IsComplete: true},
	}

	logs := []string{"14:04:03", "CHARLIE", "Start"}
	mappedInfo := MapExistingSessionsForUser(sessionInfo, logs, &[]string{"14:02:03"}[0])
	assert.Equal(ms, expectedSessionInfo, mappedInfo)
}

//master test
func TestComputeUserSessions(ct *testing.T) {
	expectedMap := map[string]SessionOutput{
		"ALICE99": {Name: "ALICE99", TimeSpent: 240, ActiveSessions: 4},
		"CHARLIE": {Name: "CHARLIE", TimeSpent: 37, ActiveSessions: 3},
	}
	out := ComputeUserSessions(expectedOut)
	assert.Equal(ct, out, expectedMap)
}

//master test case 2
func TestComputeUserSessionsCase2(ct *testing.T) {

	inputCase2 := []string{
		"14:02:05 CHARLIE End",
		"14:02:34 ALICE99 End",
		"14:02:58 ALICE99 Start",
		"14:03:02 CHARLIE Start",
		"14:03:33 ALICE99 Start",
		"14:03:35 ALICE99 End",
		"14:03:37 CHARLIE End",
		"14:04:05 ALICE99 End",
		"14:04:23 ALICE99 End",
		"14:04:41 CHARLIE Start",
	}

	expectedMap := map[string]SessionOutput{
		"ALICE99": {Name: "ALICE99", TimeSpent: 236, ActiveSessions: 4},
		"CHARLIE": {Name: "CHARLIE", TimeSpent: 35, ActiveSessions: 3},
	}
	out := ComputeUserSessions(inputCase2)
	assert.Equal(ct, out, expectedMap)
}
