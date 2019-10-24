package httpcodes

type Response struct {
	Status  int64       `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Token   interface{} `json:"token,omitempty"`
}

func (r Response) BadRequest() Response {
	return Response{
		Status:  400,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) Unauthorized() Response {
	return Response{
		Status:  401,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) Forbidden() Response {
	return Response{
		Status:  403,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) NotFound() Response {
	return Response{
		Status:  404,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) NotAcceptable() Response {
	return Response{
		Status:  406,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) UnsupportedMediaType() Response {
	return Response{
		Status:  415,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) InternalServerError() Response {
	return Response{
		Status:  500,
		Success: false,
		Message: r.Message,
	}
}

func (r Response) Ok() Response {
	return Response{
		Status:  200,
		Success: true,
		Message: r.Message,
		Data:    r.Data,
		Token:   r.Token,
	}
}

func (r Response) Created() Response {
	return Response{
		Status:  201,
		Success: true,
		Message: r.Message,
		Data:    r.Data,
		Token:   r.Token,
	}
}
