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
	thumbNail := request.QueryStringParameters["thumbnail"]

	taggedImageUsecase := usecase.NewTaggedImageUsecase()

	if thumbNail == "true" {
		response, err := taggedImageUsecase.GetThumbNailImageById(id, *userName)
		if err != nil {
			return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
		}
		return util.CreateOKResponse(response)
	} else {
		response, err := taggedImageUsecase.GetImageById(id, *userName)
		if err != nil {
			return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
		}
		return util.CreateOKResponse(response)
	}

}

func main() {
	lambda.Start(handler)
}
