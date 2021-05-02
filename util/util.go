package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func GetUsernameFromHeader(request events.APIGatewayProxyRequest) (*string, error) {
	idToken := request.Headers["x-authorization"]
	sections := strings.Split(idToken, ".")
	payload := sections[1]
	decoded, error := base64.RawStdEncoding.DecodeString(payload)
	if error != nil {
		return nil, error
	}
	claim := new(Claim)
	json.Unmarshal(decoded, claim)
	return &claim.UserName, nil
}

type Claim struct {
	UserName string `json:"cognito:username"`
}

func CreateOKResponse(response interface{}) (events.APIGatewayProxyResponse, error) {
	if response == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Create Response Failed",
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
			Body:       "Create Response Failed",
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: statusCode,
	}, nil
}
