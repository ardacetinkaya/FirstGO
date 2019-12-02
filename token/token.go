package token

//forked from https://github.com/harlow/authtoken
import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func GetRequestToken(req *http.Request) (string, error) {

	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("Authorization header required")
	}

	if !strings.HasPrefix(authorizationHeader, "Basic") && !strings.HasPrefix(authorizationHeader, "Bearer") {
		return "", errors.New("Authorization requires Basic/Bearer scheme")
	}

	if strings.HasPrefix(authorizationHeader, "Basic") {
		basicAuthorization, err := base64.StdEncoding.DecodeString(authorizationHeader[len("Basic"):])
		if err != nil {
			return "", errors.New("Base64 encoding issue")
		}
		credentials := strings.Split(string(basicAuthorization), ":")
		return credentials[0], nil
	} else {
		return authorizationHeader[len("Bearer"):], nil
	}
}
