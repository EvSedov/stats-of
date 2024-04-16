# Пример реализации функций поиска данных в redis

```go
// ---------------------------------------------------
opts, err := redis.CreateOptions()
	if err != nil {
		logger.Log.Fatal("Не удалось загрузить настройки Redis", zap.Error(err))
	}

	// Создание клиента Redis
	redisClient := redis.NewRedisClient(opts)

	// Пример использования: Пинг до Redis для проверки соединения
	err = redisClient.Ping(context.Background())
	if err != nil {
		logger.Log.Fatal("Ошибка соединения с Redis", zap.Error(err))
	} else {
		logger.Log.Info("Успешное соединение с Redis")
	}

	// Использование FindKeysByPattern для поиска ключей по шаблону
	pattern := "*Pattern*" // Можно заменить на любой другой шаблон
	keys, err := redisClient.FindKeysByPattern(pattern)
	if err != nil {
		logger.Log.Fatal("Ошибка при поиске ключей", zap.Error(err))
	}

	// Логгирование найденных ключей
	for _, key := range keys {
		logger.Log.Info("Найден ключ", zap.String("key", key))
	}

	key := "key"

	// Получение значения по ключу
	value, err := redisClient.FindKeyByGetRequest(key)
	if err != nil {
		logger.Log.Fatal("Ошибка при получении значения из Redis", zap.Error(err))
	} else if value == "" {
		logger.Log.Info("Ключ не найден", zap.String("key", key))
	} else {
		logger.Log.Info("Полученное значение", zap.String("key", key), zap.String("value", value))
	}
// ---------------------------------------------------
```
