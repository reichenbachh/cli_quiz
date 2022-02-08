package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const quiz_api_url string ="https://api.trivia.willfry.co.uk/questions?limit=10"

type Questions struct {
	Category string `json:"category"`
	CorrectAnswer string `json:"correctAnswer"`
	Id int `json:"id"`
	IncorrectAnswers  []string `json:"incorrectAnswers"`
	Question string `json:"question"`
	Type string `json:"type"`
}

func main() {
	// "category": "Society and Culture", 
    // "correctAnswer": "Japan",
    // "id": 18776,
    // "incorrectAnswers": [
    //   "Mongolia",
    //   "Nigeria",
    //   "Angola"
    // ],
    // "question": "In what country do people speak the Language they call Nihongo?",
    // "type": "Multiple Choice"
	resp, err := http.Get(quiz_api_url)
	if err != nil{
		exit("Error fectching trivia questions")
	}

	resBody,_ := io.ReadAll(resp.Body)

	var questions []Questions

	err = json.Unmarshal(resBody,&questions)

	if err != nil{
		println(err)
		exit("error parsing json response")
	}

	reader:= bufio.NewReader(os.Stdin)

	var correct int;

	for i,q := range questions{
		possibleAnswers:=[]string{q.CorrectAnswer}
		possibleAnswers = append(possibleAnswers, q.IncorrectAnswers...)
		shuffledAnswers := shuffleAnswers(possibleAnswers)
		fmt.Printf("Question %d. %s \n",i+1,q.Question)
			for _,a:= range shuffledAnswers{
			fmt.Printf("%v \n",a)
		}
		input, _,_:= reader.ReadLine()
		if  string(input) == q.CorrectAnswer {
			correct++
		}

	}
	fmt.Printf("You got %v  answers correct",correct)
}


func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

func shuffleAnswers(answerList []string)(shuffledAnsewers []string){
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(answerList), func(i, j int) {
        answerList[i], answerList[j] = answerList[j], answerList[i]
    })
	return answerList
}