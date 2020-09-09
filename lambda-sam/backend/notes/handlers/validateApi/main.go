package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/codedeploy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	CloudformationStackName = os.Getenv("CLOUDFORMATION_STACK_NAME")
	APIURL                  = ""
)

func createNote() (map[string]interface{}, error) {
	note := map[string]interface{}{
		"title":       "Test note",
		"description": "This is a test note",
	}
	noteData, _ := json.Marshal(&note)

	resp, err := http.Post(fmt.Sprintf("%s/notes", APIURL), "application/json", bytes.NewBuffer(noteData))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, errors.New("failed to create note")
	}

	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return nil, err
	}

	return note, nil
}

func getNotes() error {
	resp, err := http.Get(fmt.Sprintf("%s/notes", APIURL))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("failed to get notes")
	}
	return nil
}

func removeNote(noteID string) error {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/notes/%s", APIURL, noteID), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New("failed to remove note")
	}
	return nil
}

type response struct{}

func handler(req codedeploy.PutLifecycleEventHookExecutionStatusInput) (response, error) {
	reqData, _ := json.Marshal(req)
	log.Printf("event data: %v", string(reqData))

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := cloudformation.New(session)
	stacks, _ := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(CloudformationStackName),
	})

	stack := stacks.Stacks[0]
	for _, output := range stack.Outputs {
		if *output.OutputKey == "NotesAPI" {
			APIURL = *output.OutputValue
		}
	}

	codedeploySvc := codedeploy.New(session)

	params := &codedeploy.PutLifecycleEventHookExecutionStatusInput{
		DeploymentId:                  req.DeploymentId,
		LifecycleEventHookExecutionId: req.LifecycleEventHookExecutionId,
		Status:                        aws.String(codedeploy.DeploymentStatusSucceeded),
	}

	note, err := createNote()
	if err != nil {
		log.Printf("failed to create note: %v", err.Error())
		params.Status = aws.String(codedeploy.DeploymentStatusFailed)
		codedeploySvc.PutLifecycleEventHookExecutionStatus(params)
		return response{}, nil
	}

	if err := getNotes(); err != nil {
		log.Printf("failed to get notes: %v", err.Error())
		params.Status = aws.String(codedeploy.DeploymentStatusFailed)
		codedeploySvc.PutLifecycleEventHookExecutionStatus(params)
		return response{}, nil
	}

	if err := removeNote(note["id"].(string)); err != nil {
		log.Printf("failed to remove note: %v", err.Error())
		params.Status = aws.String(codedeploy.DeploymentStatusFailed)
		codedeploySvc.PutLifecycleEventHookExecutionStatus(params)
		return response{}, nil
	}

	if _, err := codedeploySvc.PutLifecycleEventHookExecutionStatus(params); err != nil {
		log.Printf("failed to update deployment status: %v", err.Error())
		return response{}, err
	}

	log.Printf("deployment succeeded")

	return response{}, nil
}

func main() {
	lambda.Start(handler)
}
