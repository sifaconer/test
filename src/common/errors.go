package common

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Type     string `json:"type"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Details  any    `json:"details,omitempty"`
	Internal error  `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

// NotFoundError para recursos que no existen
func NotFoundError(resource string) AppError {
	return AppError{
		Type:    "not_found",
		Code:    http.StatusNotFound,
		Message: resource,
	}
}

// ValidationError para errores de validación
func ValidationError(details any) AppError {
	return AppError{
		Type:    "validation_error",
		Code:    http.StatusBadRequest,
		Message: "Error de validación",
		Details: details,
	}
}

// UnauthorizedError para problemas de autenticación
func UnauthorizedError(message string) AppError {
	if message == "" {
		message = "No autorizado"
	}
	return AppError{
		Type:    "unauthorized",
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// ForbiddenError para problemas de permisos
func ForbiddenError(message string) AppError {
	if message == "" {
		message = "Acceso prohibido"
	}
	return AppError{
		Type:    "forbidden",
		Code:    http.StatusForbidden,
		Message: message,
	}
}

// BadRequestError para peticiones incorrectas
func BadRequestError(message string) AppError {
	return AppError{
		Type:    "bad_request",
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// ConflictError para conflictos en operaciones
func ConflictError(message string) AppError {
	return AppError{
		Type:    "conflict",
		Code:    http.StatusConflict,
		Message: message,
	}
}

// DatabaseError para errores relacionados con la base de datos
func DatabaseError(err error) AppError {
	return AppError{
		Type:     "database_error",
		Code:     http.StatusInternalServerError,
		Message:  "Error en la base de datos",
		Internal: err,
	}
}

// ThirdPartyError para errores en servicios externos
func ThirdPartyError(service string, err error) AppError {
	return AppError{
		Type:     "third_party_error",
		Code:     http.StatusServiceUnavailable,
		Message:  fmt.Sprintf("Error en servicio externo: %s", service),
		Internal: err,
	}
}

// InternalServerError para errores internos genéricos
func InternalServerError(err error) AppError {
	return AppError{
		Type:     "internal_server_error",
		Code:     http.StatusInternalServerError,
		Message:  "Error interno del servidor",
		Internal: err,
	}
}
