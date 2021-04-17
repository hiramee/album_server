package main

import (
	"album-server/application/usecase"
	"album-server/util"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tags := request.QueryStringParameters["tag"]
	tagSlice := strings.Split(tags, ",") // a=b,c or a=b&a=c
	userName, _ := util.GetUsernameFromHeader(request)

	taggedImageUsecase := usecase.NewTaggedImageUsecase()

	headers := map[string]string{"Access-Control-Allow-Origin": "*"}
	response, err := taggedImageUsecase.ListByTagNames(*userName, tagSlice)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}

	jsonBytes, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
