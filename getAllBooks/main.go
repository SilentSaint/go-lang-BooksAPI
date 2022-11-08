package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
	"os"
)

//type Book struct {
//	ID   string `json:"id"`
//	Name string `json:"name"`
//}xx

func findAll() (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving aws credentials",
		}, err
	}
	dyDb := dynamodb.NewFromConfig(cfg)
	res, err := dyDb.Scan(context.TODO(), &dynamodb.ScanInput{TableName: aws.String(os.Getenv("TABLE_NAME"))})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning dynamodb",
		}, err
	}
	response, err := json.Marshal(res.Items)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding string value",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}
func main() {
	lambda.Start(findAll)
}
