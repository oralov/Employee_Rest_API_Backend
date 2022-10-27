package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TABLE_NAME = "employees"

var dynamo *dynamodb.DynamoDB

// Connecting to DynamoDB with hardcoded credentials (not recomended)
func ConnectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKIAQ6ANJCLF74EVWAJK", "T+crpqMpduWRCA/w6NVjycA25bbQkLO+FE0u/Mym", ""), // for testing purposes only
	})))
}

// This function creates a table if one doesn't exist yet. If a table already exsits,
// then the message will be printed to the console
func createTable() {
	_, err := dynamo.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

// This function creates a new employee in the employees table
func PutItem(employee Employee) {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(employee.Id))},
			"Name": {
				S: aws.String(employee.Name)},
			"Age": {
				N: aws.String(strconv.Itoa(employee.Age))},
			"Salary": {
				N: aws.String(strconv.Itoa(employee.Salary))},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

// this function deletes an employee by the specified id
func DeleteItem(id int) {
	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

// This function is for getting a single employee from the database by id
func GetItem(id int) (employee Employee) {
	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &employee)
	if err != nil {
		panic(err)
	}
	return employee
}

// This function scans through the entrie table and returns a JSON list of all employees in the table
func GetItems() (data string) {
	var records []Employee
	err := dynamo.ScanPages(&dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	}, func(page *dynamodb.ScanOutput, last bool) bool {
		recs := []Employee{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &recs)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}

		records = append(records, recs...)

		return true
	})
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(records)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	return string(j)

}
