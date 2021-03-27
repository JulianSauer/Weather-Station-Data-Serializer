package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "os"
)

const TABLE = "WeatherStation"

type WeatherData struct {
    Source    string `json:"source,omitempty"`
    Timestamp string `json:"timestamp"`
    Temperature   float64 `json:"temperature"`
    Humidity      uint8   `json:"humidity"`
    WindSpeed     float64 `json:"windSpeed"`
    GustSpeed     float64 `json:"gustSpeed"`
    Rain          float64 `json:"rain"`
    WindDirection float32 `json:"windDirection"`
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
    fmt.Println("Connecting to DynamoDB")
    session := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    db := dynamodb.New(session)

    fmt.Printf("Received %d message(s):", len(snsEvent.Records))
    for _, record := range snsEvent.Records {
        fmt.Println(record.SNS.Message)
        d := WeatherData{}
        e := json.Unmarshal([]byte(record.SNS.Message), &d)
        handle(e)

        d.Source = "WeatherStation"
        put(&d, db)
    }
}

func put(data *WeatherData, db *dynamodb.DynamoDB) {
    attribute, e := dynamodbattribute.MarshalMap(data)
    handle(e)

    input := &dynamodb.PutItemInput{
        Item:      attribute,
        TableName: aws.String(TABLE),
    }

    _, e = db.PutItem(input)
    handle(e)
}

func handle(e error) {
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
}

func main() {
    lambda.Start(handler)
}
