package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	admin      string = "ADMIN"
	superAdmin string = "SUPER_ADMIN"
)

func isSuperAdmin(profile string) bool {
	return profile == superAdmin
}

func getProfile(headers map[string]string) (string, error) {
	token, ok := headers["Authorization"]
	if !ok || token == "" {
		return "", errors.New("authorization token not present")
	}

	tokenSplit := strings.Split(token, ".")
	if len(tokenSplit) < 2 {
		return "", errors.New("invalid auth token")
	}

	userDataBytes, err := base64.RawURLEncoding.DecodeString(tokenSplit[1])
	if err != nil {
		log.Error(err)
		return "", errors.New("unable to decode token body")
	}

	var userData map[string]interface{}
	err = json.Unmarshal(userDataBytes, &userData)
	if err != nil {
		log.Error(err)
		return "", errors.New("unable to unmarshal token body")
	}

	profile, ok := userData["profile"]
	if !ok || profile == "" {
		return "", errors.New("profile attribute not present")
	}

	switch p := profile.(type) {
	case string:
		return p, nil
	default:
		return "", errors.New("profile attribute type not string")
	}
}

func profileWithSuperAdmin(headers map[string]string) error {
	profile, err := getProfile(headers)
	if err != nil {
		log.Error(err)
		return ErrBadRequest
	}

	if !isSuperAdmin(profile) {
		return ErrUnauthorized
	}

	return nil
}
