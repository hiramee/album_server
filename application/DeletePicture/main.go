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

	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	taggedImageUsecase.DeleteTaggedImage(*userName, id)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}
