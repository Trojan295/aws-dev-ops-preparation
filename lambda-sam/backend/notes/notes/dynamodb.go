package notes

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type NoteItem struct {
	NoteID      string
	Title       string
	Description string
}

type DynamoDBNotesRepository struct {
	TableName string
	session   *session.Session
}

func NewDynamoDBNotesRepositoryFromEnv() NotesRepository {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		panic("missing DYNAMODB_TABLE environment variable")
	}

	return &DynamoDBNotesRepository{
		TableName: tableName,
	}
}

func (repo *DynamoDBNotesRepository) dynamoDBClient() *dynamodb.DynamoDB {
	if repo.session == nil || repo.session.Config.Credentials.IsExpired() {
		repo.session = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	return dynamodb.New(repo.session)
}

func (repo *DynamoDBNotesRepository) GetNotes() ([]Note, error) {
	svc := repo.dynamoDBClient()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(repo.TableName),
	})
	if err != nil {
		return []Note{}, err
	}

	notes := make([]Note, 0, *result.Count)

	for _, i := range result.Items {
		item := NoteItem{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return notes, err
		}

		note := Note{
			ID:          uuid.MustParse(item.NoteID),
			Title:       item.Title,
			Description: item.Description,
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (repo *DynamoDBNotesRepository) AddNote(note *Note) error {
	note.ID = uuid.New()

	item := NoteItem{
		NoteID:      note.ID.String(),
		Title:       note.Title,
		Description: note.Description,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	svc := repo.dynamoDBClient()
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.TableName),
	})
	return err
}

func (repo *DynamoDBNotesRepository) RemoveNote(note *Note) error {
	svc := repo.dynamoDBClient()
	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(repo.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(note.ID.String()),
			},
		},
	})
	return err
}
