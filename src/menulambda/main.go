package main

import (
	"encoding/json"
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

	// Get entire spreadsheet
	spreadsheetService := sheets.NewSpreadsheetsService(sheetsService)
	sheetResp, sheetErr := spreadsheetService.Get(spreadsheetID).Do()

	if sheetErr != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Was unable to get the spreadsheet from SPREADSHEET_ID")
	}

	var sectionSheet *sheets.Sheet
	for i, sheetElement := range sheetResp.Sheets {
		fmt.Printf("%d: %s\n", i, sheetElement.Properties.Title)
		if sheetElement.Properties.Title == "Sheet1" {
			sectionSheet = sheetElement
		}
	}

	fmt.Println(sectionSheet.Data)

	for i, gridDataObj := range sectionSheet.Data {
		fmt.Printf("%d: %s\n", i, gridDataObj.ColumnMetadata)
	}

	// Using Sheet Range
	sheetResp2, sheetErr2 := spreadsheetService.Values.Get(spreadsheetID, "Sheet1!A1:G200").Do()
	if sheetErr2 != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Was unable to get the spreadsheet from SPREADSHEET_ID")
	}

	marshalledJson, err := sheetResp2.MarshalJSON()
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Was unable to get the spreadsheet from SPREADSHEET_ID")
	}
	var jsonResp map[string]interface{}
	json.Unmarshal(marshalledJson, &jsonResp)
	fmt.Println(jsonResp)

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
