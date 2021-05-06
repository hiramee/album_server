package util

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func CreateOKResponse(response interface{}) (events.APIGatewayProxyResponse, error) {
	if response == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "create response failed",
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func CreateErrorResponse(response interface{}, errorType ErrorType, err error) (events.APIGatewayProxyResponse, error) {
	statusCode := GetStatusCodeByErrorType(errorType)
	if response == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
		}, nil
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "create response failed",
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: statusCode,
	}, nil
}
