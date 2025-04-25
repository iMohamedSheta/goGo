package requests

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=4,max=12,unique_db=users_username"`
	Email     string `json:"email" validate:"required,email,unique_db=users_email"`
	FirstName string `json:"first_name" validate:"required,min=4,max=12"`
	LastName  string `json:"last_name" validate:"required,min=4,max=12"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
}

func (r *RegisterRequest) Messages() map[string]string {
	return map[string]string{
		"Username.required":  "Username is required",
		"Username.unique_db": "Username must be unique",
		"Username.min":       "Username must be at least 4 characters",
		"Username.max":       "Username must be at most 12 characters",
		"Email.required":     "Email is required",
		"Email.unique_db":    "Email must be unique",
		"Email.email":        "Email must be a valid email address",
		"FirstName.required": "first name is required",
		"FirstName.min":      "first name must be at least 4 characters",
		"FirstName.max":      "first name must be at most 12 characters",
		"LastName.required":  "LastName is required",
		"LastName.min":       "LastName must be at least 4 characters",
		"LastName.max":       "LastName must be at most 12 characters",
		"Password.required":  "Password is required",
		"Password.min":       "Password must be at least 8 characters",
		"Password.max":       "Password must be at most 30 characters",
	}
}
