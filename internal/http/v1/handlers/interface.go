package handlers

// Интерфейс для обязательного формата ответа
type ResponseFormatter interface {
	FormatResponse() interface{}
}

// Структура для успешного ответа
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Структура для ошибки
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
