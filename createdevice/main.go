package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-playground/validator"
)

// Define the standard request body schema
type deviceInfo struct {
	ID          string `json:"id" validate:"required"`
	DeviceModel string `json:"deviceModel" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Note        string `json:"note" validate:"required"`
	Serial      string `json:"serial" validate:"required"`
}

var validate *validator.Validate

func main() {

	//Init the AWS request handler
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	validate = validator.New()

	var req deviceInfo
	err := json.Unmarshal([]byte(event.Body), &req)

	//dyanmodb configs
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// validate input json
	missingStr := ""
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			missingStr += err.Field() + ", "
		}
		return events.APIGatewayProxyResponse{Body: string("Some values are missing, " + missingStr), StatusCode: 400}, nil
	}

	if strings.Contains(req.ID, "/devices/") {
		req.ID = strings.Split(req.ID, "/")[2]
	}

	av, err := dynamodbattribute.MarshalMap(req)

	// Create item in table Movies
	tableName := "devices"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Internal Server Error"), StatusCode: 500}, nil
	}

	// Send back the response
	return events.APIGatewayProxyResponse{Body: string("Created"), StatusCode: 201}, nil
}