package apperror

import "fmt"

const (
	noChangesError          = "you have no changes"
	timeoutError            = "timeout error"
	internalServerError     = "internal server error"
	invalidCredentialsError = "invalid credentials"
	badRequestError         = "bad request"
)

// Ошибка при попытке создать существующий ресурс
type AlreadyExistsError struct {
	Err    error
	Entity string
	Field  string
	Value  string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with %s '%s' is already exists", e.Entity, e.Field, e.Value)
}

func (e *AlreadyExistsError) Unwrap() error {
	return e.Err
}

// Ошибка при попытке изменения ресурса на идентичный
type NoChangesError struct {
	Err error
}

func (e *NoChangesError) Error() string {
	return noChangesError
}

func (e *NoChangesError) Unwrap() error {
	return e.Err
}

// Ошбика при истечении дедлайна контекста
type TimeoutError struct {
	Err error
}

func (e *TimeoutError) Error() string {
	return timeoutError
}

// Ошибка при отсутствии искомого ресурса
type NotFoundError struct {
	Err    error
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// Ошибка при внутренней проблеме сервера
type InternalServerError struct {
	Err error
}

func (e *InternalServerError) Error() string {
	return internalServerError
}

func (e *InternalServerError) Unwrap() error {
	return e.Err
}

// Ошибка при вводе неправильных данных пользователя
type InvalidCredentialsError struct {
	Err error
}

func (e *InvalidCredentialsError) Error() string {
	return invalidCredentialsError
}

func (e *InvalidCredentialsError) Unwrap() error {
	return e.Err
}

// Ошибка при вводе некорректных данных
type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return badRequestError
}

func (e *BadRequestError) Unwrap() error {
	return e.Err
}

// Ошибка при работе с JWT токеном
type TokenError struct {
	Err     error
	Message string
}

func (e *TokenError) Error() string {
	return e.Message
}

func (e *TokenError) Unwrap() error {
	return e.Err
}
