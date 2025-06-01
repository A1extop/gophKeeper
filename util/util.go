package util

const (
	Admin    = "admin"
	Attendee = "attendee"
)

func IsValidUserTypeForRegistration(userType string) bool {
	validUserType := map[string]bool{
		"attendee": true,
	}

	_, ok := validUserType[userType]
	return ok
}

func IsValidUserTypeForAdmin(userType string) bool {
	validUserType := map[string]bool{
		"admin":    true,
		"attendee": true,
	}

	_, ok := validUserType[userType]
	return ok
}
