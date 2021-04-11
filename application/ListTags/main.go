package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Claim struct {
	UserName string `json:"cognito:username"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TagUsecase := usecase.NewTagUsecase()
	idToken := request.Headers["Auth"]
	sections := strings.Split(idToken, ".")
	payload := sections[1]
	decoded, _ := base64.RawStdEncoding.DecodeString(payload)
	claim := new(Claim)
	json.Unmarshal(decoded, claim)
	userName := claim.UserName

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
