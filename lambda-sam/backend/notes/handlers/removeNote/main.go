package main

import (
	"log"
	"notes/notes"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	repo = notes.NewDynamoDBNotesRepositoryFromEnv()
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	noteID := request.PathParameters["NoteID"]
	note := &notes.Note{
		ID: uuid.MustParse(noteID),
	}

	if err := repo.RemoveNote(note); err != nil {
		log.Printf("cannot remove note: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
