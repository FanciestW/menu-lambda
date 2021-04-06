package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	sheetsAPIKey, sheetsAPIKeyFound   = os.LookupEnv("GOOGLE_SHEETS_API_KEY")
	spreadsheetID, spreadsheetIDFound = os.LookupEnv("SPREADSHEET_ID")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !sheetsAPIKeyFound {
		return events.APIGatewayProxyResponse{}, errors.New("GOOGLE_SHEETS_API_KEY is not set in environment variables")
	} else if !spreadsheetIDFound {
		return events.APIGatewayProxyResponse{}, errors.New("SPREADSHEET_ID is not set in environment variables")
	}

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithAPIKey(sheetsAPIKey))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	spreadsheetService := sheets.NewSpreadsheetsService(sheetsService)
	sheetResp, sheetErr := spreadsheetService.Values.Get(spreadsheetID, "Sheet1!A1:G200").Do()
	if sheetErr != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Was unable to get the spreadsheet from SPREADSHEET_ID")
	}

	marshalledJson, err := sheetResp.MarshalJSON()
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Was unable to get the spreadsheet from SPREADSHEET_ID")
	}
	var jsonResp map[string]interface{}
	json.Unmarshal(marshalledJson, &jsonResp)
	// fmt.Println(jsonResp)

	var labels []string = make([]string, len(sheetResp.Values[0]))
	for i, label := range sheetResp.Values[0] {
		labels[i] = fmt.Sprint(label)
	}
	fmt.Printf("%s\n", labels)

	var dataList []map[string]interface{}
	for _, row := range sheetResp.Values[1:] {
		var tmpObj map[string]interface{} = make(map[string]interface{})
		for j, val := range row {
			if val == nil || strings.TrimSpace(fmt.Sprint(val)) == "" {
				continue
			}
			tmpObj[labels[j]] = val
		}
		dataList = append(dataList, tmpObj)
	}

	// Encode to JSON
	dataMap := make(map[string][]map[string]interface{})
	dataMap["data"] = dataList
	jsonData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Unable to encode data into JSON format")
	}
	fmt.Println(string(jsonData))

	// AWS S3 PutObject
	awsCreds := credentials.NewSharedCredentials("~/.aws/credentials", "default")
	awsCredsValues, err := awsCreds.Get()
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Unable to get AWS credentials")
	}
	fmt.Println(awsCredsValues.AccessKeyID)
	fmt.Println(awsCredsValues.SecretAccessKey)
	fmt.Println(awsCredsValues.SessionToken)

	// s3Session, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-east-1"),
	// 	Credentials: credentials.NewEnvCredentials(),
	// })

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, World"),
		StatusCode: 200,
	}, nil
}

func main() {
	// AWS S3 PutObject
	awsCreds := credentials.NewSharedCredentials("/Users/william/.aws/credentials", "default")
	awsCredsValues, err := awsCreds.Get()
	if err != nil {
		fmt.Println("Unable to get AWS credentials")
		return
	}
	fmt.Println(awsCredsValues.AccessKeyID)
	fmt.Println(awsCredsValues.SecretAccessKey)
	fmt.Println(awsCredsValues.SessionToken)

	// s3Session, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-east-1"),
	// 	Credentials: credentials.NewEnvCredentials(),
	// })
	return
	_, isLambda := os.LookupEnv("LAMBDA_TASK_ROOT")
	if isLambda {
		lambda.Start(handler)
	} else {
		var emptyRequest events.APIGatewayProxyRequest
		res, err := handler(emptyRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}

}
