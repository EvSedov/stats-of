## Для добавления данных в redis необходимо указать filepath и количество строк в функции HandleCsvToDb и добавить в main.go следующий код:

```
// Initialize and handle database operations
	if err := storagetestsutils.HandleCsvToDb(); err != nil {
		logger.Log.Fatal("Ошибка при обработке данных CSV", zap.Error(err))
	}
```

<!-- end of code -->

<!-- Вход в командную строку redis -->
```
redis-cli -h localhost -p 6379
```