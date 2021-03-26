package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	sheetResp, sheetErr := spreadsheetService.Get(spreadsheetID).Do()

	if sheetErr != nil {
		fmt.Println("Was unable to get the spreadsheet from SPREADSHEET_ID")
	} else {
		fmt.Println(sheetResp)
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, World"),
		StatusCode: 200,
	}, nil
}

func main() {
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
