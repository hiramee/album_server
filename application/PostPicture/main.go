package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"encoding/base64"
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

	tagUsecase := usecase.NewTagUsecase()
	if err := tagUsecase.CreateIfAbsent(*userName, req.Tags); err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}
	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	decoded, err := base64.StdEncoding.DecodeString(req.Picture)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}
	taggedImageUsecase.Create(*userName, req.Tags, req.Ext, decoded)

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
