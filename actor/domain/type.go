package domain

import (
	"errors"
	"strings"
)

type (
	Gender    string
	OauthName string
)

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"

	OauthNameAndroid OauthName = "Android"
	OauthNameIOS     OauthName = "iOS"
	OauthNameWeb     OauthName = "Web"
	OauthNameDesktop OauthName = "Desktop"
)

var (
	OauthNames = []OauthName{OauthNameAndroid, OauthNameWeb, OauthNameIOS, OauthNameDesktop}
)

func GenderValidation(gender Gender) (Gender, error) {
	switch strings.ToLower(string(gender)) {
	case strings.ToLower(string(GenderMale)):
		return GenderMale, nil
	case strings.ToLower(string(GenderFemale)):
		return GenderFemale, nil
	default:
		return "", errors.New("unknown gender")
	}
}

func OauthNameValidation(name OauthName) (OauthName, error) {
	switch strings.ToLower(string(name)) {
	case strings.ToLower(string(OauthNameAndroid)):
		return OauthNameAndroid, nil
	case strings.ToLower(string(OauthNameIOS)):
		return OauthNameIOS, nil
	case strings.ToLower(string(OauthNameWeb)):
		return OauthNameWeb, nil
	case strings.ToLower(string(OauthNameDesktop)):
		return OauthNameDesktop, nil
	default:
		return "", errors.New("unknown oauth name")
	}
}
