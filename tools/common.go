package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Stage string

const (
	StageComplete Stage = "complete"
	StageBuilding Stage = "building"
	StageError    Stage = "error"
	StageInitial  Stage = "initial"
)

type irexRunInfo struct {
	RunID     int    `json:"run_id"`
	RunName   string `json:"run_name"`
	Timestamp string `json:"timestamp"`
	Stage     Stage  `json:"stage,omitempty"`
	LogFile   string `json:"log_file"`
}

var (
	runInfoPath = filepath.Join(os.TempDir(), "irex_run_info.json")
)

func saveRunFile(runInfo irexRunInfo) error {
	jsonData, err := json.MarshalIndent(runInfo, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(runInfoPath, jsonData, 0644)
	return err
}

func addRunInfoFile() {
	// remove older run info files
	os.RemoveAll(runInfoPath)
	// create initial run info file
	cwd, _ := os.Getwd()
	runInfo := irexRunInfo{
		RunID:     0,
		RunName:   "Initial Run",
		Timestamp: time.Now().Format(time.RFC3339),
		Stage:     StageInitial,
		// use current working directory + irextools-log.log as log file
		LogFile: filepath.Join(cwd, "irextools-log.log"),
	}
	saveRunFile(runInfo)
}

func getRunInfoPath() (irexRunInfo, error) {
	// read run info from file
	data, err := os.ReadFile(runInfoPath)
	if err != nil {
		return irexRunInfo{}, err
	}
	var info irexRunInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return irexRunInfo{}, err
	}
	return info, nil
}

func updateBuildingStage() {
	runInfo, err := getRunInfoPath()
	if err != nil {
		return
	}
	runInfo.Stage = StageBuilding
	runInfo.Timestamp = time.Now().Format(time.RFC3339)
	saveRunFile(runInfo)
}

func updateCompleteStage() {
	runInfo, err := getRunInfoPath()
	if err != nil {
		return
	}
	runInfo.RunID++
	runInfo.Stage = StageComplete
	runInfo.Timestamp = time.Now().Format(time.RFC3339)
	saveRunFile(runInfo)
}
