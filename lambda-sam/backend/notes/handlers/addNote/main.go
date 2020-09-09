package main

import (
	"encoding/json"
	"log"
	"notes/notes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	repo = notes.NewDynamoDBNotesRepositoryFromEnv()
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	note := &notes.Note{}
	if err := json.Unmarshal([]byte(request.Body), note); err != nil {
		log.Printf("cannot unmarshal payload to note: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	if err := repo.AddNote(note); err != nil {
		log.Printf("cannot add note: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	noteData, err := json.Marshal(note)
	if err != nil {
		log.Printf("cannot marshal note: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(noteData),
	}, nil
}

func main() {
	lambda.Start(handler)
}
