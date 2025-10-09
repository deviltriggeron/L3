# Сервив оповещений

В данном сервисе создаются оповещения, через указанное время приходит уведомление

## Запуск

```
    go run cmd/main.go
```

- Ожидаемый ответ "Listen and running :8080", означает что всё хорошо и сервис запущен

## Примеры:

Создание оповещения:
```
    curl -X POST http://localhost:8080/notify \
        -H "Content-Type: application/json" \
        -d '{
            "message": "Напомнить о встрече",
            "send_at": "2025-10-09T20:36:00+09:00"
        }'
```

- При успешном создании {"Notify succesfully create, ID":{id}}

- Если что то не получилось получите ошибку парсинга json

Узнать статус оповещения

```
    curl -X GET http://localhost:8080/notify/{id}
```

- При успешном получении {"Notify status {id}":STATUS}

- Если что то не получилось "cannot parse ID: {id}" или "not found ID: {id}"

```
    curl -X DELETE http://localhost:8080/notify/1
```

- При успешном получении {"Notify {id} deleted}

- Если что то не получилось "cannot parse ID: {id}" или "not found ID: {id}"
   

    