package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type APIHandlerFunc func(http.ResponseWriter, *http.Request, string, int)

var validActions = []string{"on", "off", "toggle"}

type API struct {
	config *Config
}

func NewAPI(config *Config) *API {
	return &API{
		config: config,
	}
}

func (a *API) ValidateRequest(f APIHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		action, err := parseAction(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		groupId, err := parseIntField(r, "group_id")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		f(w, r, action, groupId)
	})
}

func (a *API) HandleOutletGroupRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	w.Write([]byte(r.RequestURI + "\n"))
	w.Write([]byte(action + "\n"))
	w.Write([]byte(fmt.Sprintf("group_id: %d\n", groupId)))

	log.Println(r.RequestURI)
}

func (a *API) HandleOutletRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletId, err := parseIntField(r, "outlet_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(r.RequestURI + "\n"))
	w.Write([]byte(action + "\n"))
	w.Write([]byte(fmt.Sprintf("group_id: %d\n", groupId)))
	w.Write([]byte(fmt.Sprintf("outlet_id: %d\n", outletId)))

	log.Println(r.RequestURI)
}

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
