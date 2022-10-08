package users

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type UserInput struct {
	UserId      *int    `json:"user_id"`
	Name        *string `json:"name"`
	DateOfBirth *string `json:"date_of_birth"`
	CreatedOn   *int64  `json:"created_on"`
}

// Validates whether or not a given UserInput is valid (all fields are defined)
func (ui UserInput) validate() error {
	if ui.UserId == nil || ui.Name == nil || ui.DateOfBirth == nil || ui.CreatedOn == nil {
		return errors.New("the UserInput entity is missing required fields")
	}

	return nil
}

// generateUserOutput uses a UserInput to generate the expected UserOutput.
// On error, the object will be returned up to the point it was processed
// with the associated error.
func (ui UserInput) generateUserOutput() (userOutput UserOutput, err error) {
	userOutput = UserOutput{
		UserId: *ui.UserId,
		Name:   *ui.Name,
	}

	// attempt to extract the day of the week from the date of birth
	dateOfBirth, err := time.Parse("2006-01-02", *ui.DateOfBirth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while parsing the user's DOB: %s", err.Error())
		return userOutput, err
	}
	userOutput.WeekdayOfBirth = dateOfBirth.Weekday().String()

	// attempt to extract the time in the appropriate timezone and format
	location, err := time.LoadLocation("EST")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while finding the time zone: %s", err.Error())
		return userOutput, err
	}
	userOutput.CreatedOn = time.Unix(*ui.CreatedOn, 0).In(location).Format(time.RFC3339)

	return userOutput, nil
}
