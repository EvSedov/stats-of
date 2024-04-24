package storagetestsutils

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"stats-of/internal/logger"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// CsvDbManager управляет операциями с CSV файлами и базой данных.
type CsvDbManager struct {
	FilePath    string
	RedisClient *redis.Client
}

// NewCsvDbManager создает новый экземпляр CsvDbManager.
func NewCsvDbManager(filePath string, redisClient *redis.Client) *CsvDbManager {
	return &CsvDbManager{
		FilePath:    filePath,
		RedisClient: redisClient,
	}
}

func InitDb() *redis.Client {
	// Загрузка переменных из файла .env
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Ошибка при загрузке файла .env", zap.Error(err))
	}

	// Получение конфигурации из переменных окружения
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB") // Получаем значение как строку

	// Преобразуем значение DB из строки в число
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		logger.Log.Fatal("Ошибка при преобразовании REDIS_DB в число", zap.Error(err))
	}

	// Создаем новый клиент Redis с использованием переменных окружения
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, // Используем преобразованное значение
	})

	// Создаем контекст для вызова метода Ping
	ctx := context.Background()

	// Выполняем команду PING с контекстом
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Log.Info("Ошибка при подключении к Redis", zap.Error(err))
	} else {
		logger.Log.Info("Ответ от Redis:", zap.String("response", pong))
	}

	return rdb
}

// CountCsvRows подсчитывает количество строк в CSV файле.
func (m *CsvDbManager) CountCsvRows() (int, error) {
	file, err := os.Open(m.FilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	rowCount := 0

	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		rowCount++
	}

	return rowCount, nil
}

// ReadCsvData читает заданное количество строк из CSV-файла.
func (m *CsvDbManager) ReadCsvData(rowLimit int) ([][]string, error) {
	totalRows, err := m.CountCsvRows()
	if err != nil {
		logger.Log.Error("Ошибка при подсчете строк в файле", zap.String("path", m.FilePath), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Общее количество строк в файле", zap.String("path", m.FilePath), zap.Int("totalRows", totalRows))

	file, err := os.Open(m.FilePath)
	if err != nil {
		logger.Log.Error("Ошибка при открытии файла", zap.String("path", m.FilePath), zap.Error(err))
		return nil, err
	}
	defer file.Close()

	logger.Log.Info("Файл успешно открыт", zap.String("path", m.FilePath))

	reader := csv.NewReader(bufio.NewReader(file))

	var records [][]string

	for rowCount := 0; rowCount < rowLimit; rowCount++ {
		record, err := reader.Read()
		if err == io.EOF {
			logger.Log.Info("Достигнут конец файла", zap.String("path", m.FilePath))
			break
		}
		if err != nil {
			logger.Log.Error("Ошибка при чтении данных из файла", zap.String("path", m.FilePath), zap.Error(err))
			return nil, err
		}

		records = append(records, record)
	}

	logger.Log.Info("Данные успешно прочитаны", zap.String("path", m.FilePath), zap.Int("rowsRead", len(records)))

	return records, nil
}

// AddDataToDb добавляет данные из слайса слайсов строк в базу данных Redis.
func (m *CsvDbManager) AddDataToDb(records [][]string) error {
	ctx := context.Background() // Создаем контекст для операций с Redis

	for _, record := range records {
		chatID := record[0]
		user := record[1]
		messageType := record[2]
		lastMsgEvent := record[3]

		key := "chat:" + chatID + ":user:" + user + ":type:" + messageType
		value := lastMsgEvent

		// Использование context.Background() и time.Duration(0) для установки ключей без истечения срока действия
		err := m.RedisClient.Set(ctx, key, value, 0*time.Second).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CsvDbManager) ReadCsvHeaders() ([]string, error) {
	file, err := os.Open(m.FilePath)
	if err != nil {
		logger.Log.Error("Ошибка при открытии файла", zap.String("path", m.FilePath), zap.Error(err))
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// Чтение только первой строки для заголовков
	headers, err := reader.Read()
	if err == io.EOF {
		logger.Log.Info("Файл пуст", zap.String("path", m.FilePath))
		return nil, err
	}
	if err != nil {
		logger.Log.Error("Ошибка при чтении заголовков из файла", zap.String("path", m.FilePath), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Заголовки успешно прочитаны", zap.String("path", m.FilePath))

	return headers, nil
}

// handleCsvToDb initializes the database, reads CSV data, and stores it in the database.
func HandleCsvToDb() error {
	rdb := InitDb()

	filePath := "/home/sergey/Development/Sfera/testmeetdb/asap/Result_1.csv"
	rowLimit := 150000

	manager := NewCsvDbManager(filePath, rdb)

	records, err := manager.ReadCsvData(rowLimit)
	if err != nil {
		logger.Log.Error("Не удалось прочитать данные из CSV", zap.Error(err))
		return err // Return error instead of fatal log to allow for easier testing and error handling
	}

	err = manager.AddDataToDb(records)
	if err != nil {
		logger.Log.Error("Не удалось добавить данные в базу данных", zap.Error(err))
		return err
	}

	logger.Log.Info("Данные успешно добавлены в базу данных")
	return nil
}

// AddUsersData использует горутины для параллельной записи данных пользователей в Redis.
func (m *CsvDbManager) AddUsersData(userCount int) error {
	var wg sync.WaitGroup
	wg.Add(userCount)

	for i := 0; i < userCount; i++ {
		go func(userID int) {
			defer wg.Done()
			key := fmt.Sprintf("chat:5481:user:%d:type:CHANNEL", 121000000+userID)
			value := time.Now().Format("2006-01-02 15:04:05.000000 +00:00")
			err := m.RedisClient.Set(context.Background(), key, value, 0).Err() // Assuming you want the key to never expire.
			if err != nil {
				logger.Log.Error("Error setting value for user", zap.Int("userID", userID), zap.Error(err))
			} else {
				logger.Log.Info("Value set for user", zap.Int("userID", userID))
			}
		}(i)
	}

	wg.Wait()
	logger.Log.Info("All goroutines have finished executing")
	return nil
}
