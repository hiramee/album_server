package main

import (
	"album-server/application/usecase"
	"album-server/util"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tags := request.QueryStringParameters["tag"]
	tagSlice := strings.Split(tags, ",") // a=b,c or a=b&a=c
	userName, _ := util.GetUsernameFromHeader(request)

	taggedImageUsecase := usecase.NewTaggedImageUsecase(ctx)

	response, err := taggedImageUsecase.ListTaggedImageByTagNames(*userName, tagSlice)
	if err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}

	return util.CreateOKResponse(response)
}

func main() {
	lambda.Start(handler)
}
