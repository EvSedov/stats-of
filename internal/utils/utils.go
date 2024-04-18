package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stats-of/internal/logger"

	"go.uber.org/zap"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	// Логирование попытки преобразования payload в JSON
	logger.Log.Info("Attempting to marshal payload to JSON", zap.Any("payload", payload))

	response, err := json.Marshal(payload)
	if err != nil {
		// Логирование ошибки маршалинга
		logger.Log.Error("Failed to marshal payload", zap.Error(err))

		// Попытка отправить ответ об ошибке
		respondErr := RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		if respondErr != nil {
			// Логирование ошибки при отправке HTTP ответа об ошибке
			logger.Log.Error("Failed to send error response", zap.Error(respondErr))
			return err
		}
		return fmt.Errorf("failed to marshall payload: %w", err)
	}

	// Установка заголовков ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Отправка JSON ответа
	if _, err = w.Write(response); err != nil {
		// Логирование ошибки при записи ответа
		logger.Log.Error("Failed to write JSON response", zap.Error(err))
		return fmt.Errorf("failed to write response: %w", err)
	}

	// Логирование успешной отправки ответа
	logger.Log.Info("JSON response sent successfully", zap.Int("statusCode", code), zap.ByteString("response", response))
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) error {
	return RespondWithJSON(w,
		code,
		struct {
			Error string `json:"error"`
		}{
			Error: message,
		})
}

func RespondWith400(w http.ResponseWriter, message string) error {
	return RespondWithError(w,
		http.StatusBadRequest,
		message)
}

func RespondWith404(w http.ResponseWriter) error {
	return RespondWithError(w,
		http.StatusNotFound,
		http.StatusText(http.StatusNotFound))
}

func RespondWith500(w http.ResponseWriter) error {
	return RespondWithError(w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError))
}

func SuccessRespondWith200(w http.ResponseWriter, payload interface{}) error {
	return RespondWithJSON(w,
		http.StatusOK,
		payload)
}

func SuccessRepondWith201(w http.ResponseWriter, payload interface{}) error {
	return RespondWithJSON(w,
		http.StatusCreated,
		payload)
}
