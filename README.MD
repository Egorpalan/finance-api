### Шаги для запуска

1. **Клонируйте репозиторий:**

   ```bash
   git clone git@github.com:Egorpalan/finance-api.git
   cd finance-api

2. Запустите приложение
    ```bash
   make run

### Для запуска тестов выполните

```bash
   make test
```

### Примеры запросов

1. Создание пользователя
    ```bash
    curl -X POST "http://localhost:8080/user" -H "Content-Type: application/json" -d '{
    "balance": 100.00
    }'



2. Пополнение баланса
    ```bash
   curl -X POST "http://localhost:8080/balance/top-up" -H "Content-Type: application/json" -d '{
    "user_id": 1,
    "amount": 50.00
   }'

3. Перевод денег
    ```bash
   curl -X POST "http://localhost:8080/transfer" -H "Content-Type: application/json" -d '{      
    "sender_id": 3,  
    "receiver_id": 4,
    "amount": 30.00
   }'


4. Получение транзакций
    ```bash
   curl -X GET "http://localhost:8080/transactions/3"
