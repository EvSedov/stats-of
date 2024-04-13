## Пример реализации функций поиска данных в redis


```
// ---------------------------------------------------
	redisService := storage.NewRedisService()
	keys, err := redisService.FindKeysByPattern("*pattern*")
	if err != nil {
		logger.Log.Fatal("Ошибка при поиске ключей", zap.Error(err))
	}
	for _, key := range keys {
		logger.Log.Info("Найден ключ", zap.String("key", key))
	}

	key := "chat:5481:user:121190002:type:CHANNEL"

	// Получение значения по ключу
	value, err := storage.NewRedisService().FindKeyByGetRequest(key)
	if err != nil {
		logger.Log.Fatal("Ошибка при получении значения из Redis", zap.Error(err))
	}

	// Вывод полученного значения
	logger.Log.Info("Полученное значение", zap.String("value", value))
	// ---------------------------------------------------

  ```