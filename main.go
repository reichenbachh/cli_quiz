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

const quiz_api_url string ="https://api.trivia.willfry.co.uk/questions?limit=3"
// number of questions can be changed by changing the limit query paramter
var timeLimit int = 30

type questions struct {
	Category string `json:"category"`
	CorrectAnswer string `json:"correctAnswer"`
	Id int `json:"id"`
	IncorrectAnswers  []string `json:"incorrectAnswers"`
	Question string `json:"question"`
	Type string `json:"type"`
}

func main() {
	//example response
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

	questions:= FetchQuiz()

	var correct int = 0
	reader:= bufio.NewReader(os.Stdin)

	timer := time.NewTimer(time.Duration(timeLimit)*time.Second)

	for i,q := range questions{
		possibleAnswers:=[]string{q.CorrectAnswer}
		possibleAnswers = append(possibleAnswers, q.IncorrectAnswers...)
		shuffledAnswers := shuffleAnswers(possibleAnswers)
		fmt.Printf("Question %d. %s \n",i+1,q.Question)
			for _,a:= range shuffledAnswers{
			fmt.Printf("%v \n",a)
		}
		input := make(chan string)
		go answerChannel(reader,input)

		select{
			case <- timer.C:
			fmt.Printf("Time is up!, you scored %d out of %d \n",correct, len(questions))
			os.Exit(1)
			case answer := <- input:
				if answer == q.CorrectAnswer{
					correct++
				}
		}
	}
	fmt.Printf("you scored %d out of %d \n",correct, len(questions))
}


func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

func FetchQuiz()([]questions){
	resp, err := http.Get(quiz_api_url)
	if err != nil{
		exit("Error fectching trivia questions")
	}

	resBody,_ := io.ReadAll(resp.Body)

	var questions []questions
	
	err = json.Unmarshal(resBody,&questions)

	if err != nil{
		println(err)
		exit("error parsing json response")
	}
	return questions
	
}

func answerChannel(reader *bufio.Reader,answer chan string){
	input, _,_:= reader.ReadLine()
	answer <- string(input)
}


func shuffleAnswers(answerList []string)(shuffledAnsewers []string){
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(answerList), func(i, j int) {
		answerList[i], answerList[j] = answerList[j], answerList[i]
    })
	return answerList
}