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

type problem struct {
	question string
	answer   string
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

func loadProblems(filename string) ([]problem, error) {
	fd, err := os.Open(problemsFilename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var problems []problem

	problemSet := csv.NewReader(fd)
	for {
		rec, err := problemSet.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		problems = append(problems, problem{
			question: strings.TrimSpace(rec[0]),
			answer:   strings.ToLower(strings.TrimSpace(rec[1])),
		})
	}

	return problems, nil
}

func runQuiz(problems []problem, timeLimit int) int {
	fmt.Println("Press ENTER to start.")
	inputChan := userInputReader()
	_ = <-inputChan

	timeout := time.After(time.Duration(timeLimit) * time.Second)
	correct := 0

Qloop:
	for _, p := range problems {
		fmt.Printf("\n%s ", p.question)

		select {
		case userInput := <-inputChan:
			if p.answer == userInput {
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
