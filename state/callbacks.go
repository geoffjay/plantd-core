package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/geoffjay/plantd/core/service"

	log "github.com/sirupsen/logrus"
)

type createScopeCallback struct {
	name  string
	store *Store
}

type deleteScopeCallback struct {
	name  string
	store *Store
}

type deleteCallback struct {
	name  string
	store *Store
}

type getCallback struct {
	name  string
	store *Store
}

type setCallback struct {
	name  string
	store *Store
}

// Execute callback function to handle `create-scope` requests.
func (cb *createScopeCallback) Execute(msgBody string) ([]byte, error) {
	var (
		scope   string
		found   bool
		request service.RawRequest
	)

	log.Tracef("name: %s", cb.name)
	log.Tracef("body: %s", msgBody)

	if err := json.Unmarshal([]byte(msgBody), &request); err != nil {
		return []byte(`{"error": "` + err.Error() + `"}`), err
	}

	if scope, found = request["service"].(string); !found {
		return []byte(`{"error": "service required for create-scope request"}`), errors.New("`service` missing")
	}

	if cb.store.HasScope(scope) {
		// this shouldn't fail, just report to the caller
		msg := fmt.Sprintf("{\"error\":\"the scope %s already exists\"}", scope)
		return []byte(msg), nil
	}

	err := cb.store.CreateScope(scope)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	return []byte("{}"), nil
}

// Execute callback function to handle `delete-scope` requests.
func (cb *deleteScopeCallback) Execute(msgBody string) ([]byte, error) {
	var (
		scope   string
		found   bool
		request service.RawRequest
	)

	log.Tracef("name: %s", cb.name)
	log.Tracef("body: %s", msgBody)

	if err := json.Unmarshal([]byte(msgBody), &request); err != nil {
		return []byte(`{"error": "` + err.Error() + `"}`), err
	}

	if scope, found = request["service"].(string); !found {
		return []byte(`{"error": "service required for delete-scope request"}`), errors.New("`service` missing")
	}

	if !cb.store.HasScope(scope) {
		// this shouldn't fail, just report to the caller
		msg := fmt.Sprintf("{\"error\":\"the scope %s doesn't exist\"}", scope)
		return []byte(msg), nil
	}

	err := cb.store.DeleteScope(scope)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	return []byte("{}"), nil
}

// Execute callback function to handle `delete` requests.
func (cb *deleteCallback) Execute(msgBody string) ([]byte, error) {
	var (
		scope   string
		key     string
		found   bool
		request service.RawRequest
	)

	log.Tracef("name: %s", cb.name)
	log.Tracef("body: %s", msgBody)

	if err := json.Unmarshal([]byte(msgBody), &request); err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	if scope, found = request["service"].(string); !found {
		return []byte(`{"error": "service required for delete request"}`), errors.New("`service` missing")
	}

	if key, found = request["key"].(string); !found {
		return []byte(`{"error": "key required for delete request"}`), errors.New("`key` missing")
	}

	err := cb.store.Delete(scope, key)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	return []byte("{}"), nil
}

// Execute callback function to handle `get` requests.
func (cb *getCallback) Execute(msgBody string) ([]byte, error) {
	var (
		scope   string
		key     string
		found   bool
		request service.RawRequest
	)

	log.Tracef("name: %s", cb.name)
	log.Tracef("body: %s", msgBody)

	if err := json.Unmarshal([]byte(msgBody), &request); err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	if scope, found = request["service"].(string); !found {
		return []byte(`{"error": "service required for get request"}`), errors.New("`service` missing")
	}

	if key, found = request["key"].(string); !found {
		return []byte(`{"error": "key required for get request"}`), errors.New("`key` missing")
	}

	value, err := cb.store.Get(scope, key)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	log.Tracef("value: %s", value)
	return []byte(`{"key": "` + key + `", "value": "` + value + `"}`), nil
}

// Execute callback function to handle `set` requests.
func (cb *setCallback) Execute(msgBody string) ([]byte, error) {
	var (
		scope   string
		key     string
		value   string
		found   bool
		request service.RawRequest
	)

	log.Tracef("name: %s", cb.name)
	log.Tracef("body: %s", msgBody)

	if err := json.Unmarshal([]byte(msgBody), &request); err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	if scope, found = request["service"].(string); !found {
		return []byte(`{"error": "service required for get request"}`), errors.New("`service` missing")
	}

	if key, found = request["key"].(string); !found {
		return []byte(`{"error": "key required for get request"}`), errors.New("`key` missing")
	}

	if value, found = request["value"].(string); !found {
		return []byte(`{"error": "value required for get request"}`), errors.New("`value` missing")
	}

	err := cb.store.Set(scope, key, value)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
		return []byte(msg), err
	}

	log.Tracef("value: %s", value)
	return []byte(`{"key": "` + key + `", "value": "` + value + `"}`), nil
}
