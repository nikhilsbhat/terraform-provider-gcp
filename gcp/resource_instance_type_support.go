package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(filename string) ([]byte, error) {
	if _, dirneuerr := os.Stat(filename); os.IsNotExist(dirneuerr) {
		return nil, dirneuerr
	}

	content, conterr := ioutil.ReadFile(filename)
	if conterr != nil {
		return nil, conterr
	}
	return content, nil
	//return nil, nil
}

func jsonDecode(data []byte, i interface{}) error {
	err := json.Unmarshal(data, i)
	if err != nil {

		switch err.(type) {
		case *json.UnmarshalTypeError:
			return unknownTypeError(data, err)
		case *json.SyntaxError:
			return syntaxError(data, err)
		}
	}

	return nil
}

func syntaxError(data []byte, err error) error {
	syntaxErr, ok := err.(*json.SyntaxError)
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	start := bytes.LastIndex(data[:syntaxErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %s\n%s",
		line, syntaxErr, data[start:end])
	return err
}

func unknownTypeError(data []byte, err error) error {
	unknownTypeErr, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	start := bytes.LastIndex(data[:unknownTypeErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %s\n%s\nThe data type you entered for the value is wrong",
		line, unknownTypeErr, data[start:end])
	return err
}

func getStringOfMessage(g interface{}) string {
	switch g.(type) {
	case string:
		return g.(string)
	case error:
		return g.(error).Error()
	default:
		return "unknown messagetype"
	}
}
