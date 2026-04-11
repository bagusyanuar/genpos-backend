package response

type APIResponse[T any] struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       T           `json:"data,omitempty"`
	Errors     any         `json:"errors,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	Limit       int   `json:"limit"`
}

// Success returns a standard success response
func Success[T any](data T, message string) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// SuccessWithPagination returns a success response with pagination metadata
func SuccessWithPagination[T any](data T, pagination Pagination, message string) APIResponse[T] {
	return APIResponse[T]{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
	}
}

// Error returns a standard error response
func Error(message string) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Message: message,
	}
}

// ErrorWithDetails returns an error response with additional error details
func ErrorWithDetails(message string, errors any) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

// ValidationError returns a standard 422 Unprocessable Entity error response
func ValidationError(errors any) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	}
}
