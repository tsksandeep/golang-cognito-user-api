package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type DefaultRespBody struct {
	Message string `json:"message,omitempty"`
}

func writeGatewayProxyResponse(body string, err error) events.APIGatewayProxyResponse {
	res := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "'GET,PUT,DELETE,OPTIONS'",
		},
	}
	if err != nil {
		res.StatusCode = GetCodeFromError(err)

		if body != "" {
			respByte, _ := json.Marshal(&DefaultRespBody{Message: body})
			res.Body = string(respByte)
		}

		return res
	}

	res.StatusCode = http.StatusOK
	res.Body = body

	return res
}
