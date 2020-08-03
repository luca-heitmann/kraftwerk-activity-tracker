package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	outputFileEnvVar = "OUTPUT_FILE"
	defaultOutputFile = "data.csv"

	intervalEnvVar = "INTERVAL_SECONDS"
	defaultInterval = 5

	defaultIntervalUnit = time.Second
)

type config struct {
	OutputFile string
	Interval time.Duration
	IntervalUnit time.Duration
}

func newConfig() *config {
	outputFile := os.Getenv(outputFileEnvVar)
	intervalVar, intervalEnvExists := os.LookupEnv(intervalEnvVar)
	var interval time.Duration
	if !intervalEnvExists {
		interval = defaultInterval
	} else {
		intervalParsed, err := strconv.Atoi(intervalVar)
		if err != nil {
			log.Printf("Invalid format of variable %s with value %s. Using default value.\n", intervalEnvVar, intervalVar)
		}
		interval = time.Duration(intervalParsed)
	}
	
	return &config{
		OutputFile: outputFile,
		Interval: interval * defaultIntervalUnit,
		IntervalUnit: defaultIntervalUnit,
	}
}