package main

import (
	"album-server/application/usecase"
	"album-server/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	id := request.PathParameters["id"]

	headers := map[string]string{"Access-Control-Allow-Origin": "*"}

	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	taggedImageUsecase.Delete(*userName, id)

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
