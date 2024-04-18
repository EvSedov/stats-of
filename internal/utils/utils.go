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
	// Логирование попытки отправить ошибочный ответ
	logger.Log.Info("Attempting to respond with error", zap.Int("statusCode", code), zap.String("errorMessage", message))

	err := RespondWithJSON(w,
		code,
		struct {
			Error string `json:"error"`
		}{
			Error: message,
		})

	if err != nil {
		// Логирование ошибки при попытке отправить JSON ответ об ошибке
		logger.Log.Error("Failed to respond with error JSON", zap.Error(err))
		return err
	}

	// Логирование успешной отправки ошибочного ответа
	logger.Log.Info("Error response sent successfully", zap.Int("statusCode", code))
	return nil
}

func RespondWith400(w http.ResponseWriter, message string) error {
	logger.Log.Info("Responding with 400 Bad Request", zap.String("message", message))
	err := RespondWithError(w, http.StatusBadRequest, message)
	if err != nil {
		logger.Log.Error("Failed to send 400 response", zap.Error(err))
	}
	return err
}

func RespondWith404(w http.ResponseWriter) error {
	logger.Log.Info("Responding with 404 Not Found")
	err := RespondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	if err != nil {
		logger.Log.Error("Failed to send 404 response", zap.Error(err))
	}
	return err
}

func RespondWith500(w http.ResponseWriter) error {
	logger.Log.Info("Responding with 500 Internal Server Error")
	err := RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	if err != nil {
		logger.Log.Error("Failed to send 500 response", zap.Error(err))
	}
	return err
}

func SuccessRespondWith200(w http.ResponseWriter, payload interface{}) error {
	logger.Log.Info("Responding with 200 OK", zap.Any("payload", payload))
	err := RespondWithJSON(w, http.StatusOK, payload)
	if err != nil {
		logger.Log.Error("Failed to send 200 response", zap.Error(err))
	}
	return err
}

func SuccessRespondWith201(w http.ResponseWriter, payload interface{}) error {
	logger.Log.Info("Responding with 201 Created", zap.Any("payload", payload))
	err := RespondWithJSON(w, http.StatusCreated, payload)
	if err != nil {
		logger.Log.Error("Failed to send 201 response", zap.Error(err))
	}
	return err
}
