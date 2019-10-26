// Package httpcodes extends the capabilities of a Representation custom type so it Representation with
// different HTTP status codes.
package httpcodes

// Representation is a custom type used to represent the state of the API when a response is returnned
type Representation struct {
	Status  int64       `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Token   interface{} `json:"token,omitempty"`
}

// BadRequest returns a representation of the state of the API with an HTTP/x.x 400 Bad Request code
func (r Representation) BadRequest() Representation {
	return Representation{
		Status:  400,
		Success: false,
		Message: r.Message,
	}
}

// Unauthorized returns a representation of the state of the API with an HTTP/x.x 401 Unauthorized code
func (r Representation) Unauthorized() Representation {
	return Representation{
		Status:  401,
		Success: false,
		Message: r.Message,
	}
}

// Forbidden returns a representation of the state of the API with an HTTP/x.x 403 Forbidden code
func (r Representation) Forbidden() Representation {
	return Representation{
		Status:  403,
		Success: false,
		Message: r.Message,
	}
}

// NotFound returns a representation of the state of the API with an HTTP/x.x 404 Not Found code
func (r Representation) NotFound() Representation {
	return Representation{
		Status:  404,
		Success: false,
		Message: r.Message,
	}
}

// NotAcceptable returns a representation of the state of the API with an HTTP/x.x 406 Not Acceptable code
func (r Representation) NotAcceptable() Representation {
	return Representation{
		Status:  406,
		Success: false,
		Message: r.Message,
	}
}

// UnsupportedMediaType returns a representation of the state of the API with an HTTP/x.x 415 Unsupported Media Type code
func (r Representation) UnsupportedMediaType() Representation {
	return Representation{
		Status:  415,
		Success: false,
		Message: r.Message,
	}
}

// InternalServerError returns a representation of the state of the API with an HTTP/x.x 500 Internal Server Error code
func (r Representation) InternalServerError() Representation {
	return Representation{
		Status:  500,
		Success: false,
		Message: r.Message,
	}
}

// Ok returns a representation of the state of the API with an HTTP/x.x 200 Ok code
func (r Representation) Ok() Representation {
	return Representation{
		Status:  200,
		Success: true,
		Message: r.Message,
		Data:    r.Data,
		Token:   r.Token,
	}
}

// Created returns a representation of the state of the API with an HTTP/x.x 201 Created code
func (r Representation) Created() Representation {
	return Representation{
		Status:  201,
		Success: true,
		Message: r.Message,
		Data:    r.Data,
		Token:   r.Token,
	}
}
