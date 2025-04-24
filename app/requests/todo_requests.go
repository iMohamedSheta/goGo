package requests

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=5"`
}

// Return validation failure messages
func (r *CreateTodoRequest) Messages() map[string]string {
	return map[string]string{
		"Title.required":       "The title is required.",
		"Title.min":            "The title must be at least 3 characters.",
		"Title.max":            "The title must not exceed 100 characters.",
		"Description.required": "The description is required.",
		"Description.min":      "The description must be at least 5 characters.",
	}
}
