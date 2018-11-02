package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func parseIntField(r *http.Request, fieldName string) (int, error) {
	rawValue := r.FormValue(fieldName)
	if rawValue == "" {
		return 0, fmt.Errorf("%s field missing", fieldName)
	}

	id, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func parseAction(r *http.Request) (string, error) {
	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) < 3 {
		return "", errors.New("invalid url path")
	}

	action := urlParts[3]
	if !isValidAction(action) {
		return "", fmt.Errorf("%s is not a valid action", action)
	}

	return action, nil
}

func isValidAction(action string) bool {
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}
