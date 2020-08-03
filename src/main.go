package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	// Get configuration from environment vars
	conf := newConfig()
	log.Printf("Config Summary:\n- Output file: %s\n- Interval: %s\n", conf.OutputFile, conf.Interval)

	// Loop and log Kraftwerk client counter to CSV
	nextTime := time.Now().Truncate(conf.IntervalUnit)
	for {
		// Get results from Kraftwerk client counter
		counter, max, err := getClientCounter()
		if err != nil {
			log.Println("Unable to get results: ", err)
		}

		// Write results to output file
		if conf.OutputFile != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05 -0700 MST")
			csvData := fmt.Sprintf("%s;%d;%d\n", timestamp, counter, max)
			if err := writeToDataFile(conf.OutputFile, csvData); err != nil {
				log.Println("Unable to write to data file: ", err)
			}
		}

		// Print results
		log.Printf("Received results successfully. Counter: %d / Max: %d\n", counter, max)

		// Sleep until next execution
		nextTime = nextTime.Add(conf.Interval)
		time.Sleep(time.Until(nextTime))
	}
}

func getClientCounter() (int, int, error) {
	// Request the gym-clientcounter page.
	response, err := http.Get("https://www.boulderado.de/boulderadoweb/gym-clientcounter/index.php?mode=get&token=eyJhbGciOiJIUzI1NiIsICJ0eXAiOiJKV1QifQ.eyJjdXN0b21lciI6IktyYWZ0d2VyayJ9.L6rb7R_lt_auWPtJFtTIWziXOcKhUMNPPhBiwtWGnrY")
	if err != nil {
		return -1, -1, err
	}
	defer response.Body.Close()

	// Read response data in
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, -1, err
	}

	// Create a regular expression to find counter
	counterRegex := regexp.MustCompile("<div data-value=\"(\\d+)\" class=\"actcounter zoom\">")
	counterResult := counterRegex.FindStringSubmatch(string(body))

	// Create a regular expression to find max
	maxRegex := regexp.MustCompile("<div data-value=\"(\\d+)\" class=\"freecounter zoom\">")
	maxResult := maxRegex.FindStringSubmatch(string(body))

	// Parse results
	if (counterResult == nil || maxResult == nil) && (len(counterResult) < 1 || len(maxResult) < 1) {
		return -1, -1, fmt.Errorf("Unable to find result in body: %s", string(body))
	}

	counter, err := strconv.Atoi(counterResult[1])
	if err != nil {
		return -1, -1, err
	}
	max, err := strconv.Atoi(maxResult[1])
	if err != nil {
		return -1, -1, err
	}

	return counter, max, nil
}

func writeToDataFile(outputFile string, data string) error {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}
