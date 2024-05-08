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
	"strings"
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

func InitDb() (*redis.Client, error) {
	// Загрузка переменных из файла .env
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка при загрузке файла .env: %v", err)
	}

	// Получение конфигурации из переменных окружения
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	// Преобразуем значение DB из строки в число
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при преобразовании REDIS_DB в число: %v", err)
	}

	// Создаем новый клиент Redis с использованием переменных окружения
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Создаем контекст для вызова метода Ping
	ctx := context.Background()

	// Выполняем команду PING с контекстом
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к Redis: %v", err)
	}

	return rdb, nil
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
	ctx := context.Background()
	logger.Log.Info("Начало добавления данных в базу", zap.Int("totalRecords", len(records)))

	for index, record := range records {
		chatID := record[0]
		user := record[1]
		lastMsgEvent := record[3]

		userKey := "user:" + user
		chatUsersKey := "chat:" + chatID + ":users"

		// Добавляем пользователя в список пользователей чата
		if err := m.RedisClient.SAdd(ctx, chatUsersKey, user).Err(); err != nil {
			logger.Log.Error("Ошибка при добавлении пользователя в список пользователей чата", zap.String("chatUsersKey", chatUsersKey), zap.Error(err))
			return err
		}

		// Обновляем данные пользователя
		if err := m.RedisClient.HSet(ctx, userKey, "last_active", lastMsgEvent, "last_chat_id", chatID).Err(); err != nil {
			logger.Log.Error("Ошибка при обновлении информации пользователя", zap.String("userKey", userKey), zap.Error(err))
			return err
		}

		// Обновляем список чатов пользователя
		if err := m.RedisClient.SAdd(ctx, userKey+":chats", chatID).Err(); err != nil {
			logger.Log.Error("Ошибка при добавлении чата в список чатов пользователя", zap.String("userKey", userKey+":chats"), zap.Error(err))
			return err
		}

		logger.Log.Info("Данные пользователя обработаны", zap.Int("recordIndex", index+1))
	}

	logger.Log.Info("Все данные успешно обновлены в базе")
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
	rdb, err := InitDb()
	if err != nil {
		logger.Log.Error("Не удалось инициализировать базу данных", zap.Error(err))
		return err
	}

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

func (m *CsvDbManager) FindChatsWithMinUsers(minUsers int64) ([]string, error) {
	ctx := context.Background()
	var cursor uint64
	chatIDs := []string{}

	for {
		var keys []string
		var err error
		keys, cursor, err = m.RedisClient.Scan(ctx, cursor, "chat:*:users", 0).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			count, err := m.RedisClient.SCard(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			if count >= minUsers { // Сравнение с заданным минимальным количеством пользователей
				chatID := key[len("chat:") : len(key)-len(":users")]
				chatIDs = append(chatIDs, chatID)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return chatIDs, nil
}

func RunChatSearch(minUsers int64) error {
	redisClient, err := InitDb()
	if err != nil {
		return fmt.Errorf("failed to initialize Redis client: %v", err)
	}

	// Создание экземпляра CsvDbManager
	manager := NewCsvDbManager("", redisClient) // Предполагаем, что FilePath не требуется

	// Получение списка чатов с минимальным количеством пользователей
	chatIDs, err := manager.FindChatsWithMinUsers(minUsers)
	if err != nil {
		return fmt.Errorf("error retrieving chat IDs: %v", err)
	}

	// Вывод полученного списка чатов
	logger.Log.Info("Output of chat list with the specified minimum number of users", zap.Int64("minUsers", minUsers), zap.Strings("chatIDs", chatIDs))
	return nil
}

func (m *CsvDbManager) GetUsersByChatCount(targetChatCount int64) ([]string, error) {
	ctx := context.Background()
	var userIdsWithTargetChatCount []string

	// Инициализируем курсор для сканирования ключей
	var cursor uint64
	var err error
	var keys []string

	// Используем SCAN в цикле
	for {
		// Scan возвращает новый курсор и порцию ключей
		keys, cursor, err = m.RedisClient.Scan(ctx, cursor, "user:*:chats", 0).Result()
		if err != nil {
			logger.Log.Error("Ошибка при сканировании ключей пользователей", zap.Error(err))
			return nil, err
		}

		// Обрабатываем каждый ключ
		for _, key := range keys {
			chatCount, err := m.RedisClient.SCard(ctx, key).Result()
			if err != nil {
				logger.Log.Error("Ошибка при получении размера множества", zap.String("key", key), zap.Error(err))
				continue // Пропускаем в случае ошибки
			}

			// Проверяем, соответствует ли количество чатов заданному
			if chatCount == targetChatCount {
				// Извлекаем идентификатор пользователя из ключа
				parts := strings.Split(key, ":")
				if len(parts) > 1 {
					userId := parts[1]
					userIdsWithTargetChatCount = append(userIdsWithTargetChatCount, userId)
				}
			}
		}

		// Если курсор равен 0, мы обработали все ключи
		if cursor == 0 {
			break
		}
	}

	logger.Log.Info("Пользователи с заданным количеством чатов найдены", zap.Int64("targetChatCount", targetChatCount), zap.Strings("userIds", userIdsWithTargetChatCount))
	return userIdsWithTargetChatCount, nil
}

func RunUserSearch(targetChatCount int64) error {
	redisClient, err := InitDb()
	if err != nil {
		return fmt.Errorf("failed to initialize Redis client: %v", err)
	}

	// Создание экземпляра CsvDbManager
	manager := NewCsvDbManager("", redisClient) // Предполагаем, что FilePath не требуется

	// Получение списка пользователей
	userIds, err := manager.GetUsersByChatCount(targetChatCount)
	if err != nil {
		logger.Log.Error("Ошибка при получении списка пользователей по количеству чатов", zap.Error(err))
		return err
	}

	logger.Log.Info("Output of user list with the specified number of chats", zap.Int64("targetChatCount", targetChatCount), zap.Strings("userIds", userIds))

	return nil
}
