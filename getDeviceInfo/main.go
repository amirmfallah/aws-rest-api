package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Define the standard request body schema
type deviceInfo struct {
	ID          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

func main() {
	//Init the AWS request handler
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// parse id parameter from url
	id := event.PathParameters["id"]

	//dyanmodb configs
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// run query on dynamodb with the requested device id
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("devices"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, nil
	}

	// if
	if result.Item == nil {
		return events.APIGatewayProxyResponse{Body: "Device not found", StatusCode: 404}, nil
	}

	device := deviceInfo{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &device)

	resp, _ := json.Marshal(device)

	// Send back the response
	return events.APIGatewayProxyResponse{Body: string(resp), StatusCode: 200}, nil
}
