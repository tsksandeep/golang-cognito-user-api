package main

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Username  string `json:"username"`
	UserEmail string `json:"useremail"`
}

type Users struct {
	Users          []map[string]interface{} `json:"users"`
	TotalItemCount *int                     `json:"totalItemCount"`
}

func EnableConfirmUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		log.Error("Enable confirm user: missing username in body")
		return writeGatewayProxyResponse("missing username in body", ErrBadRequest)
	}

	log.Info("Enable confirm user: ", user.Username)

	enableUserInput := &cognito.AdminEnableUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(user.Username),
	}

	_, err = cognitoClient.AdminEnableUser(enableUserInput)
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error("Enable confirm user: ", err)
		switch awsErr.Code() {
		case cognito.ErrCodeNotAuthorizedException:
			return writeGatewayProxyResponse("Invalid access token", ErrUnauthorized)
		case cognito.ErrCodeUserNotFoundException:
			return writeGatewayProxyResponse("User not found", ErrNotFound)
		default:
			return writeGatewayProxyResponse("", ErrInternalServerError)
		}
	}

	return writeGatewayProxyResponse("successfully enabled confirmed user", nil)
}

func DisableConfirmUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		log.Error("Disable confirm user: missing username in body")
		return writeGatewayProxyResponse("missing username in body", ErrBadRequest)
	}

	log.Info("Disable confirm user: ", user.Username)

	disableUserInput := &cognito.AdminDisableUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(user.Username),
	}

	_, err = cognitoClient.AdminDisableUser(disableUserInput)
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error("Disable confirm user: ", err)
		switch awsErr.Code() {
		case cognito.ErrCodeNotAuthorizedException:
			return writeGatewayProxyResponse("Invalid access token", ErrUnauthorized)
		case cognito.ErrCodeUserNotFoundException:
			return writeGatewayProxyResponse("User not found", ErrNotFound)
		default:
			return writeGatewayProxyResponse("", ErrInternalServerError)
		}
	}

	return writeGatewayProxyResponse("successfully disabled confirmed user", nil)
}

func DeleteUnconfirmUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		log.Error("Delete unconfirm user: missing username in body")
		return writeGatewayProxyResponse("missing username in body", ErrBadRequest)
	}

	log.Info("Delete unconfirm user: ", user.Username)

	deleteUserInput := &cognito.AdminDeleteUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(user.Username),
	}

	_, err = cognitoClient.AdminDeleteUser(deleteUserInput)
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error("Delete unconfirm user: ", err)
		switch awsErr.Code() {
		case cognito.ErrCodeNotAuthorizedException:
			return writeGatewayProxyResponse("Invalid access token", ErrUnauthorized)
		case cognito.ErrCodeUserNotFoundException:
			return writeGatewayProxyResponse("User not found", ErrNotFound)
		default:
			return writeGatewayProxyResponse("", ErrInternalServerError)
		}
	}

	// err = sendMail(user.UserEmail, SubjectRejected, HtmlBodyRejected, TextBodyRejected)
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	log.Info("email sent to address: " + user.UserEmail)
	// }

	return writeGatewayProxyResponse("successfully deleted unconfirmed user", nil)
}

func ConfirmUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var user User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		log.Error("Confirm user: missing username in body")
		return writeGatewayProxyResponse("missing username in body", ErrBadRequest)
	}

	log.Info("Confirm user: ", user.Username)

	confirmSignupInput := &cognito.AdminConfirmSignUpInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(user.Username),
	}

	_, err = cognitoClient.AdminConfirmSignUp(confirmSignupInput)
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error("Confirm user: ", err)
		switch awsErr.Code() {
		case cognito.ErrCodeNotAuthorizedException:
			return writeGatewayProxyResponse("Invalid access token", ErrUnauthorized)
		case cognito.ErrCodeUserNotFoundException:
			return writeGatewayProxyResponse("User not found", ErrNotFound)
		default:
			return writeGatewayProxyResponse("Internal server error", ErrInternalServerError)
		}
	}

	_, err = cognitoClient.AdminUpdateUserAttributes(&cognito.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(user.Username),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
			{
				Name:  aws.String("phone_number_verified"),
				Value: aws.String("true"),
			},
		},
	})
	if err != nil {
		log.Error("Confirm user: ", err)
		return writeGatewayProxyResponse("", ErrInternalServerError)
	}

	// err = sendMail(user.UserEmail, SubjectApproved, HtmlBodyApproved, TextBodyApproved)
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	log.Info("email sent to address: " + user.UserEmail)

	// }

	return writeGatewayProxyResponse("successfully confirmed user", nil)
}

func getUserInfoFromAttributes(userEnabled *bool, attributes []*cognito.AttributeType) map[string]interface{} {
	userInfo := map[string]interface{}{}

	for _, attr := range attributes {
		name := strcase.ToLowerCamel(*attr.Name)

		switch name {
		case "sub":
			userInfo["userId"] = *attr.Value
		case "emailVerified", "phoneNumberVerified":
			parseBool, err := strconv.ParseBool(*attr.Value)
			if err != nil {
				break
			}
			userInfo[name] = parseBool

		default:
			userInfo[name] = *attr.Value
		}
	}

	if userEnabled != nil {
		userInfo["isEnabled"] = *userEnabled
	} else {
		userInfo["isEnabled"] = false
	}

	return userInfo
}

func GetUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	token, ok := request.QueryStringParameters["accessToken"]
	if !ok || token == "" {
		log.Error("Get user: access token not found")
		return writeGatewayProxyResponse("", ErrBadRequest)
	}

	userInput := &cognito.GetUserInput{
		AccessToken: aws.String(token),
	}

	userOutput, err := cognitoClient.GetUser(userInput)
	if awsErr, ok := err.(awserr.Error); ok {
		log.Error("Get user: ", err)
		switch awsErr.Code() {
		case cognito.ErrCodeNotAuthorizedException:
			return writeGatewayProxyResponse("Invalid access token", ErrUnauthorized)
		case cognito.ErrCodeUserNotFoundException:
			return writeGatewayProxyResponse("User not found", ErrNotFound)
		default:
			return writeGatewayProxyResponse("", ErrInternalServerError)
		}
	}

	userInfo := getUserInfoFromAttributes(nil, userOutput.UserAttributes)

	resBytes, err := json.Marshal(userInfo)
	if err != nil {
		log.Error("Get user: ", err)
		return writeGatewayProxyResponse("", ErrInternalServerError)
	}

	return writeGatewayProxyResponse(string(resBytes), nil)
}

func GetConfirmedUserList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	nextPageToken := ""
	userInfoList := []map[string]interface{}{}
	listUsersInput := &cognito.ListUsersInput{UserPoolId: aws.String(userPoolID)}

	for {
		if nextPageToken != "" {
			listUsersInput.PaginationToken = aws.String(nextPageToken)
		}

		listUsersOutput, err := cognitoClient.ListUsers(listUsersInput)
		if err != nil {
			log.Error("Get confirmed user list: ", err)
			break
		}

		for _, user := range listUsersOutput.Users {
			if user.UserStatus != nil && *user.UserStatus == "CONFIRMED" {
				userInfoList = append(userInfoList, getUserInfoFromAttributes(user.Enabled, user.Attributes))
			}
		}

		if listUsersOutput.PaginationToken == nil || *listUsersOutput.PaginationToken == "" {
			break
		}

		nextPageToken = *listUsersOutput.PaginationToken
	}

	totalItemCount := len(userInfoList)

	confirmedUsers := Users{
		Users:          userInfoList,
		TotalItemCount: &totalItemCount,
	}

	outputBytes, err := json.Marshal(confirmedUsers)
	if err != nil {
		log.Error("Get confirmed user list: ", err)
		return writeGatewayProxyResponse("", ErrInternalServerError)
	}

	return writeGatewayProxyResponse(string(outputBytes), nil)
}

func GetUnconfirmedUserList(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	nextPageToken := ""
	userInfoList := []map[string]interface{}{}
	listUsersInput := &cognito.ListUsersInput{UserPoolId: aws.String(userPoolID)}

	for {
		if nextPageToken != "" {
			listUsersInput.PaginationToken = aws.String(nextPageToken)
		}

		listUsersOutput, err := cognitoClient.ListUsers(listUsersInput)
		if err != nil {
			log.Error("Get unconfirmed user list: ", err)
			break
		}

		for _, user := range listUsersOutput.Users {
			if user.UserStatus != nil && *user.UserStatus == "UNCONFIRMED" {
				userInfoList = append(userInfoList, getUserInfoFromAttributes(user.Enabled, user.Attributes))
			}
		}

		if listUsersOutput.PaginationToken == nil || *listUsersOutput.PaginationToken == "" {
			break
		}

		nextPageToken = *listUsersOutput.PaginationToken
	}

	totalItemCount := len(userInfoList)

	unconfirmedUsers := Users{
		Users:          userInfoList,
		TotalItemCount: &totalItemCount,
	}

	outputBytes, err := json.Marshal(unconfirmedUsers)
	if err != nil {
		log.Error("Get unconfirmed user list: ", err)
		return writeGatewayProxyResponse("", ErrInternalServerError)
	}

	return writeGatewayProxyResponse(string(outputBytes), nil)
}
