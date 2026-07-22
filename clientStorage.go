package rui

import (
	"encoding/base64"
	"strconv"
	"strings"
)

// ClientStorage is an interface for accessing client-side key-value storage.
type ClientStorage interface {

	// Request performs an asynchronous request to obtain key-value pairs from client-side storage.
	//
	// The first argument specifies the function that is called for each key-value pair
	// (if the pair is not in storage, then the function is called with an empty string as its value).
	//
	// The second argument specifies the set of keys to request.
	// If no key is specified, then a query is performed for all possible key-value pairs.
	Request(result func(key, value string), key ...string)

	// Get returns a value by key from the client-side storage.
	//
	// If the key-value pair is not in the client-side storage, then the function returns an empty string.
	Get(key string) string

	// Set stores a key-value pair in the client-side storage.
	//
	// An empty string as a value removes the key-value pair.
	Set(key, value string)

	// RemoveAll removes all key-value pair from the client-side storage
	RemoveAll()

	handleEvent(command string, data DataObject)
}

type clientStorageData struct {
	bridge      bridge
	getResult   map[int]func(string, string)
	lastRequest int
}

func (session *sessionData) ClientStorage() ClientStorage {
	if session.clientStorage == nil {
		storage := new(clientStorageData)
		storage.bridge = session.bridge
		storage.getResult = map[int]func(string, string){}
		session.clientStorage = storage
	}
	return session.clientStorage
}

func encodeClientStorageText(text string) string {
	return base64.URLEncoding.EncodeToString([]byte(text))
}

func decodeClientStorageText(text string) string {
	result, _ := base64.URLEncoding.DecodeString(text)
	return string(result)
}

func (storage *clientStorageData) Get(key string) string {
	result := make(chan string)
	defer close(result)

	storage.Request(func(key, value string) {
		result <- value
	}, key)

	return <-result
}

func (storage *clientStorageData) Request(result func(key, value string), key ...string) {
	if result == nil {
		return
	}

	storage.lastRequest++
	request := storage.lastRequest
	storage.getResult[request] = result

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	for _, k := range key {
		if k != "" {
			if buffer.Len() > 0 {
				buffer.WriteRune(',')
			}
			buffer.WriteString(encodeClientStorageText(k))
		}
	}

	if buffer.Len() == 0 {
		storage.bridge.localStorageRequest("localStorageGetAll", request)
	} else {
		storage.bridge.localStorageRequest("localStorageGet", request, buffer.String())
	}
}

func (storage *clientStorageData) Set(key, value string) {
	key = encodeClientStorageText(key)
	if value != "" {
		storage.bridge.callFunc("localStorageSet", key, encodeClientStorageText(value))
	} else {
		storage.bridge.callFunc("localStorageRemove", key)
	}
}

func (storage *clientStorageData) RemoveAll() {
	storage.bridge.callFunc("localStorageClear")
}

func (storage *clientStorageData) handleEvent(command string, data DataObject) {
	request := 0
	if text, ok := data.PropertyValue("request"); ok {
		var err error
		if request, err = strconv.Atoi(text); err != nil {
			ErrorLog(err.Error())
			return
		}
	} else {
		ErrorLog("'request' property not found (command: " + command + ")")
		return
	}

	switch command {
	case "storageError":
		if text, ok := data.PropertyValue("error"); ok {
			ErrorLog(text)
		}

	//case "storageSuccess":

	case "storageValues":
		fn, ok := storage.getResult[request]
		if !ok {
			return
		}
		delete(storage.getResult, request)

		text, ok := data.PropertyValue("error")
		if !ok {
			ErrorLog("'values' property not found (command: storageValues)")
			return
		}

		for pair := range strings.SplitSeq(text, ";") {
			if data := strings.Split(pair, ":"); len(data) == 2 {
				key := decodeClientStorageText(data[0])
				value := decodeClientStorageText(data[1])
				fn(key, value)
			}
		}
	}
}
