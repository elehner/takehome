package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	ErrorMethodNotSupported = "Only POST is supported"
	ErrorParsingInput       = "Error parsing user input"
	ErrorProcessingInput    = "Error processing the users input"
	ErrorEncodingInput      = "Error encoding the processed data"
)

// HandleUserRequest directs the request to the appropriate call based
// on the request method.
func HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleUserInputs(w, r)
	default:
		http.Error(w, ErrorMethodNotSupported, http.StatusMethodNotAllowed)
	}
}

// handleUserInputs processes the user defined objects received
// from the client, and responds to the request with either the
// transformed inputs or an error.
func handleUserInputs(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	// Utilize a json decoder since we're dealing with a stream
	userInputs, err := processUserInputs(&body)
	if err != nil {
		http.Error(w, ErrorParsingInput, http.StatusBadRequest)
		return
	}
	// Don't bother parsing an empty request
	if userInputs == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userOutputs, err := transformUserInputs(userInputs)
	if err != nil {
		http.Error(w, ErrorProcessingInput, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(userOutputs)
	if err != nil {
		http.Error(w, ErrorEncodingInput, http.StatusInternalServerError)
		return
	}
}

// processUserInputs transforms the body of an http request into a slice of UserInputs.
// On Error, it returns nil and the associated error.
func processUserInputs(body *io.ReadCloser) (userInputs []UserInput, err error) {
	// Utilize a json decoder since we're dealing with a stream
	userInputsDecoder := json.NewDecoder(*body)
	for {
		// Loop over elements to ensure the entire message is parsed correctly
		if err = userInputsDecoder.Decode(&userInputs); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			// These (and other Fprintfs) should be moved to logs to track data over time & appropriate error levels
			fmt.Fprintf(os.Stderr, "Error occurred while parsing the user's input: %s", err.Error())
			return nil, err
		}
	}

	// Validate the parsed input objects
	for _, userInput := range userInputs {
		err = userInput.validate()
		if err != nil {
			return nil, err
		}
	}

	return userInputs, err
}

// transformUserInputs generates a slice of UserOuputs from the given slice of UserInputs.
// On Error, it will return nil and the associated error.
func transformUserInputs(userInputs []UserInput) (userOutputs []UserOutput, err error) {
	// Generate the slice of user outputs from the slice of user inputs
	userOutputs = make([]UserOutput, len(userInputs))
	for index, userInput := range userInputs {
		userOutputs[index], err = userInput.generateUserOutput()
		if err != nil {
			return nil, err
		}
	}

	return userOutputs, nil
}
