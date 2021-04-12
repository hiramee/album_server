package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	req := new(openapi.PostPicturesRequest)

	headers := map[string]string{"Access-Control-Allow-Origin": "*"}
	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}

	TagUsecase := usecase.NewTagUsecase()
	results, err := TagUsecase.ListAll(*userName)

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
