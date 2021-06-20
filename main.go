package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Problem struct {
	q string
	a string
}

func parseProblems(p [][]string) []Problem {
	problems := make([]Problem, len(p))
	for i, question := range p {
		problems[i] = Problem{
			q: question[0],
			a: question[1],
		}
	}
	return problems
}

func getProblems() []Problem {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Sometihng whent wrong Open File!!!")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Sometihng whent wrong Read File!!!")
	}

	problems := parseProblems(records)
	return problems
}

func main() {
	timeLimit := flag.Int("timer", 30, "the limit of quiz")
	flag.Parse()
	var counter = 0
	problems := getProblems()

	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	r.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, value := range problems {
		fmt.Printf("problem #%d: %s =", i+1, value.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYour time run out your score:", counter)
			return
		case answer := <-answerCh:
			if answer == value.a {
				counter++
			}
			timer.Reset(time.Duration(*timeLimit) * time.Second)
		}
	}

	fmt.Println("Your score is", counter)
}
