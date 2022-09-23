package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const apiPrefix = "/api/v1"

func getMethodPathString(method, path string) string {
	return method + " " + apiPrefix + path
}

var (
	getUser             = getMethodPathString("GET", "/user/info")
	getConfirmedUsers   = getMethodPathString("GET", "/user/confirmed")
	getUnconfirmedUsers = getMethodPathString("GET", "/user/unconfirmed")
	enableUser          = getMethodPathString("POST", "/user/enable")
	disableUser         = getMethodPathString("POST", "/user/disable")
	confirmUser         = getMethodPathString("POST", "/user/confirm")
	deleteUnconfirmUser = getMethodPathString("DELETE", "/user/unconfirm")
)

func router(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	methodAndPathString := request.HTTPMethod + " " + request.Path

	switch methodAndPathString {

	// Profile: Allow all
	case getUser:
		return GetUser(request)

	// Profile: Only Super Admin
	case getConfirmedUsers:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return GetConfirmedUserList(request)

	// Profile: Only Super Admin
	case getUnconfirmedUsers:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return GetUnconfirmedUserList(request)

	// Profile: Only Super Admin
	case confirmUser:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return ConfirmUser(request)

	// Profile: Only Super Admin
	case enableUser:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return EnableConfirmUser(request)

	// Profile: Only Super Admin
	case disableUser:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return DisableConfirmUser(request)

	// Profile: Only Super Admin
	case deleteUnconfirmUser:
		err := profileWithSuperAdmin(request.Headers)
		if err != nil {
			return writeGatewayProxyResponse(err.Error(), err)
		}
		return DeleteUnconfirmUser(request)

	default:
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "'GET,PUT,DELETE,OPTIONS'",
			},
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid method and path: " + methodAndPathString,
		}
	}
}
