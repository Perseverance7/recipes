package models

type ErrorResponse struct {
    Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
    Code    int              `json:"code"`
    Message string           `json:"message"`
    Details []FieldErrorInfo `json:"details,omitempty"`
}

type FieldErrorInfo struct {
    Field string `json:"field"`
    Error string `json:"error"`
}

