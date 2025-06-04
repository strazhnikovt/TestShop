# TestShop

Это TestShop - реализация CRUD по ТЗ, с некоторыми доработками, так как мне казалось нелогичным то, что обычный пользователь может изменять состояние продуктов, поэтому я добавил ролевую систему. 

## Основные фичи:

- **Регистрация и вход (JWT)**. Пользователь может зарегистрироваться (от 18 лет, пароль от 8 символов) и зайти по логину/паролю, получив JWT.
- **Роли**: обычный пользователь (user) и админ (admin).
- **Продукты**:
  - Любой может посмотреть список (`GET /products`).
  - Только админ может добавлять, править и удалять товары.
- **Заказы**:
  - Пользователь формирует корзину из товаров и отправляет заказ.
  - Цена каждого товара в заказе сохраняется на момент покупки (чтобы хранилась история).
  - Если на складе недостаточно товаров, заказ запрещается.

В коде всё разделено на слои:
- **delivery**: HTTP-обработчики, роуты, middleware.
- **service**: бизнес-логика.
- **repository**: доступ к базе через sqlx.
- **domain**: структуры данных.
- **config**: конфигурация (читать `default.env` или переменные окружения).
- **pkg/auth**: JWT-менеджер для выдачи и проверки токенов.
- **pkg/logging**: простой логгер.

---


## Переменные окружения


Переменные окружения лежат в `docker-compose.yml` и в `default.env`

---

## Запуск через Docker

1. Клонируем репозиторий:
   ```bash
   git clone https://github.com/strazhnikovt/TestShop
   cd TestShop
   ```

2. Запускаем:
   ```bash
   docker compose down -v
   docker compose up -d --build
   ```

---

## Роуты

### Регистрация

- **POST /register**  
  Отправляем JSON вида:
  ```json
  {
    "first_name":"Boris",
    "last_name":"Britva",
    "login":"boris@britva.com",
    "age":25,
    "is_married":false,
    "password":"borisbritva"
  }
  ```
  Валидация:
  - `age` не меньше 18.
  - `password` от 8 символов.
  - Логин уникальный.

  **Ожидается**:
  ```json
  {"id": "user_id"}
  ```
  **Ошибки**:
  - `400 age must be at least 18`
  - `400 password must be at least 8 characters`
  - `400 user already exists`

### Логин

- **POST /login**
  ```json
  {
    "login":"boris@britva.com",
    "password":"borisbritva"
  }
  ```
  **Ожидается**:
  ```json
  {"token":"<JWT_TOKEN>"}
  ```
  **Ошибки**:
  - `401 user not found`
  - `401 invalid credentials`

### Список товаров

- **GET /products**  
  Без токена.  
  **Ожидается**:
  - Если в базе пусто: `[]`
  - Если есть товары:
    ```json
    [
      {
        "id":1,
        "description":"Товар-А",
        "tags":["category1","tagA"],
        "quantity":5,
        "price":50.00
      },
      {
        "id":2,
        "description":"Товар-Б",
        "tags":["category2","tagB"],
        "quantity":2,
        "price":30.00
      }
    ]
    ```

### CRUD для товаров (админ)

#### Создать товар

- **POST /admin/products**  
  Заголовок `Authorization: Bearer <ADMIN_TOKEN>`  
  Тело:
  ```json
  {
    "description":"Товар-А",
    "tags":["category1","tagA"],
    "quantity":5,
    "price":50.00
  }
  ```
  **Успех**:
  ```json
  {"id":"product_id"}
  ```
  **Ошибки**:
  - `401` без токена
  - `403` если не админ

#### Обновить товар

- **PUT /admin/products/{id}**  
  `Authorization: Bearer <ADMIN_TOKEN>`  
  Тело:
  ```json
  {
    "description":"Товар-А (обновлён)",
    "tags":["category1","updated"],
    "quantity":3,
    "price":60.00
  }
  ```
  **Ожидается**: возвращается статус `200`

#### Удалить товар

- **DELETE /admin/products/{id}**  
  `Authorization: Bearer <ADMIN_TOKEN>`  
  **Ожидается**: `204 No Content`

### Заказы (пользователь)

- **POST /orders**  
  `Authorization: Bearer <USER_TOKEN>`
  ```json
  {
    "items":[
      {"product_id":1,"quantity":2},
      {"product_id":2,"quantity":1}
    ]
  }
  ```
  Логика:
  - Проверяем наличие на складе. Если не хватает — `400 insufficient product quantity`.
  - Сохраняем в `order_items.price_at_order` текущую цену из `products`.
  - Уменьшаем количество в таблице `products`.

  **Ожидается**:
  ```json
  {"id":"order_id"}
  ```
  **Ошибки**:
  - `401` без токена
  - `400 product not found`
  - `400 insufficient product quantity`

---

## Примеры тестовых запросов (перед этим желательно почистить базу чтоб прям эти курлы подставлять)

1. **Проверяем что товаров нет:**
   ```bash
   curl -i http://localhost:8080/products
   ```
   Ответ: `200 []`

2. **Логин админа:**
   ```bash
   curl -i -X POST http://localhost:8080/login      -H "Content-Type: application/json"      -d '{
       "login":"archibald",
       "password":"archibaldpass"
     }'
   ```
   Ответ: `200 {"token":"<ADMIN_TOKEN>"}`

3. **Добавляем Товар-А:**
   ```bash
   curl -i -X POST http://localhost:8080/admin/products      -H "Content-Type: application/json"      -H "Authorization: Bearer <ADMIN_TOKEN>"      -d '{
       "description":"Товар-А",
       "tags":["category1","tagA"],
       "quantity":5,
       "price":50.00
     }'
   ```
   Ответ: `201 {"id":"1"}`

4. **Добавляем Товар-Б:**
   ```bash
   curl -i -X POST http://localhost:8080/admin/products      -H "Content-Type: application/json"      -H "Authorization: Bearer <ADMIN_TOKEN>"      -d '{
       "description":"Товар-Б",
       "tags":["category2","tagB"],
       "quantity":2,
       "price":30.00
     }'
   ```
   Ответ: `201 {"id":"2"}`

5. **Убеждаемся, что товары есть:**
   ```bash
   curl -i http://localhost:8080/products
   ```
   Ответ: два товара

6. **Регистрация: возраст < 18 (ошибка):**
   ```bash
   curl -i -X POST http://localhost:8080/register      -H "Content-Type: application/json"      -d '{
       "first_name":"Bingo",
       "last_name":"Bolshoy",
       "login":"bolshoy@big.com",
       "age":16,
       "is_married":false,
       "password":"BolshoyMalenkiy"
     }'
   ```
   Ответ: `400 age must be at least 18`

7. **Регистрация: пароль < 8 (ошибка):**
   ```bash
   curl -i -X POST http://localhost:8080/register      -H "Content-Type: application/json"      -d '{
       "first_name":"Calvin",
       "last_name":"Strazhnikov",
       "login":"strazhnikov@gmail.com",
       "age":20,
       "is_married":false,
       "password":"Stra"
     }'
   ```
   Ответ: `400 password must be at least 8 characters`

8. **Успешная регистрация:**
   ```bash
   curl -i -X POST http://localhost:8080/register      -H "Content-Type: application/json"      -d '{
       "first_name":"Boris",
       "last_name":"Britva",
       "login":"boris@britva.com",
       "age":25,
       "is_married":false,
       "password":"borisblade"
     }'
   ```
   Ответ: `201 {"id":"3"}`

9. **Логин обычного пользователя:**
   ```bash
   curl -i -X POST http://localhost:8080/login      -H "Content-Type: application/json"      -d '{
       "login":"boris@britva.com",
       "password":"borisblade"
     }'
   ```
   Ответ: `200 {"token":"<USER_TOKEN>"}`

10. **Заказ 2го Товар-А (id=1):**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[{"product_id":1,"quantity":2}]
      }'
    ```
    Ответ: `201 {"id":"1"}`  
    (осталось `3`. Товар-А)

11. **Заказ 1го. Товар-Б (id=2):**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[{"product_id":2,"quantity":1}]
      }'
    ```
    Ответ: `201 {"id":"2"}`  
    (осталось `1`  Товар-Б)

12. **Заказ одновременно Товар-А и Товара-Б:**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[
          {"product_id":1,"quantity":1},
          {"product_id":2,"quantity":1}
        ]
      }'
    ```
    Ответ: `201 {"id":"3"}`  
    (осталось `2` Товар-А и `0` Товар-Б)

13. **Попытка заказать Товар-Б (id=2) снова:**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[{"product_id":2,"quantity":1}]
      }'
    ```
    Ответ: `400 insufficient product quantity`

14. **Добавляем Товар-В (id=3, qty=3, price=100):**
    ```bash
    curl -i -X POST http://localhost:8080/admin/products       -H "Content-Type: application/json"       -H "Authorization: Bearer <ADMIN_TOKEN>"       -d '{
        "description":"Товар-В",
        "tags":["category3"],
        "quantity":3,
        "price":100.00
      }'
    ```
    Ответ: `201 {"id":"3"}`

15. **Заказ 2го Товар-В (id=3):**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[{"product_id":3,"quantity":2}]
      }'
    ```
    Ответ: `201 {"id":"4"}`  
    (осталось `1` Товар-В, в истории цена = 100)

16. **Админ меняет цену Товара-В на 150 (qty=1):**
    ```bash
    curl -i -X PUT http://localhost:8080/admin/products/3       -H "Content-Type: application/json"       -H "Authorization: Bearer <ADMIN_TOKEN>"       -d '{
        "description":"Товар-В (обновлён)",
        "tags":["category3","updated"],
        "quantity":1,
        "price":150.00
      }'
    ```
    Ответ: `200`

17. **Заказ 1го Товар-В (цена = 150):**
    ```bash
    curl -i -X POST http://localhost:8080/orders       -H "Content-Type: application/json"       -H "Authorization: Bearer <USER_TOKEN>"       -d '{
        "items":[{"product_id":3,"quantity":1}]
      }'
    ```
    Ответ: `201 {"id":"5"}`  
    (осталось `0` Товар-В)

18. **Проверяем историю цен для Товар-В:**
    ```bash
    docker exec -it testshop-db-1 psql -U postgres -d appdb
    ```
    Внутри:
    ```sql
    SELECT order_id, product_id, quantity, price_at_order
    FROM order_items
    WHERE product_id = 3
    ORDER BY order_id;
    \q
    ```
    Результат:
    ```
     order_id | product_id | quantity | price_at_order
    ----------+------------+----------+---------------
            4 |          3 |        2 |         100.00
            5 |          3 |        1 |         150.00
    ```

---
