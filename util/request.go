package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func GetUsernameFromHeader(request events.APIGatewayProxyRequest) (*string, error) {
	idToken := request.Headers["x-authorization"]
	sections := strings.Split(idToken, ".")
	payload := sections[1]
	decoded, error := base64.RawStdEncoding.DecodeString(payload)
	if error != nil {
		return nil, error
	}
	claim := new(Claim)
	json.Unmarshal(decoded, claim)
	return &claim.UserName, nil
}

type Claim struct {
	UserName string `json:"cognito:username"`
}
