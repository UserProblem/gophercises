package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var problemsFilename string
var timeLimit int

func init() {
	const (
		defaultProblemsFilename = "problems.csv"
		usageProblemsFilename   = "The set of problems in CSV format."
		defaultTimeLimit        = 30
		usageTimeLimit          = "The time limit for completing the problem set."
	)
	flag.StringVar(&problemsFilename, "problems", defaultProblemsFilename, usageProblemsFilename)
	flag.StringVar(&problemsFilename, "p", defaultProblemsFilename, usageProblemsFilename+" (shorthand)")
	flag.IntVar(&timeLimit, "time_limit", defaultTimeLimit, usageTimeLimit)
	flag.IntVar(&timeLimit, "t", defaultTimeLimit, usageTimeLimit+" (shorthand)")
}

func main() {
	flag.Parse()

	problems, err := loadProblems(problemsFilename)
	if err != nil {
		fmt.Println("Error loading problem set: ", err)
		return
	}

	totalProblems := len(problems)
	correctAnswers := runQuiz(problems, timeLimit)

	fmt.Printf("\nYou got %d out of %d correct answers.\n", correctAnswers, totalProblems)
}

func loadProblems(filename string) ([][]string, error) {
	fd, err := os.Open(problemsFilename)
	defer fd.Close()

	var problems [][]string

	if err != nil {
		return nil, err
	}

	problemSet := csv.NewReader(fd)
	for {
		record, err := problemSet.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		record[0] = strings.TrimSpace(record[0])
		record[1] = strings.ToLower(strings.TrimSpace(record[1]))
		problems = append(problems, record)
	}

	return problems, nil
}

func runQuiz(problems [][]string, timeLimit int) int {
	fmt.Println("Press ENTER to start.")
	userInputChan := userInputReader()
	_ = <-userInputChan

	timeout := time.After(time.Duration(timeLimit) * time.Second)
	correct := 0

Qloop:
	for idx := range problems {
		question, answer := problems[idx][0], problems[idx][1]
		fmt.Printf("\n%s ", question)

		select {
		case userInput := <-userInputChan:
			if answer == userInput {
				correct++
			}

		case <-timeout:
			fmt.Println("\n\nTime's up!")
			break Qloop
		}
	}

	return correct
}

func userInputReader() <-chan string {
	outc := make(chan string)
	go func() {
		for {
			userInput := ""
			fmt.Scanln(&userInput)
			outc <- strings.ToLower(strings.TrimSpace(userInput))
		}
	}()
	return outc
}
