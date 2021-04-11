package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TagUsecase := usecase.NewTagUsecase()
	userName := request.Headers["Auth"]
	results, err := TagUsecase.ListAll(userName)
	headers := map[string]string{"Access-Control-Allow-Origin": "*"}

	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}
	response := new(openapi.GetTagsResponse)
	var tags []string
	for _, e := range results {
		tags = append(tags, e.TagName)
	}
	response.Tags = &tags
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
