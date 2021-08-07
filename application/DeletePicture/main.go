package main

import (
	"album-server/application/usecase"
	"album-server/util"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	id := request.PathParameters["id"]

	taggedImageUsecase := usecase.NewTaggedImageUsecase(ctx)
	taggedImageUsecase.DeleteTaggedImage(*userName, id)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}
