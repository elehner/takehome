package users

type UserOutput struct {
	UserId         int    `json:"user_id"`
	Name           string `json:"name"`
	WeekdayOfBirth string `json:"weekday_of_birth"`
	CreatedOn      string `json:"created_on"`
}
