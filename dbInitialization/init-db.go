package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"os"
)

type Book struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

func readBooks(fileName string) ([]Book, error) {
	books := make([]Book, 0)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return books, err
	}
	err = json.Unmarshal(data, &books)
	if err != nil {
		return books, err
	}
	return books, nil

}

func insertBook(cfg aws.Config, book Book) error {
	item, err := attributevalue.MarshalMap(book)
	fmt.Println(item)
	if err != nil {
		return err
	}
	dyDb := dynamodb.NewFromConfig(cfg)
	_, err = dyDb.PutItem(context.TODO(), &dynamodb.PutItemInput{TableName: aws.String("books"), Item: item})
	if err != nil {
		return err
	}
	return nil
}
func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	books, err := readBooks("books.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		fmt.Println("Inserting", book.Name)
		err = insertBook(cfg, book)
		if err != nil {
			log.Fatal(err)
		}
	}

}
