# Task Kaspi

## Описание
Task Kaspi - это REST API сервис на Go, разработанный с использованием чистой архитектуры. Проект использует PostgreSQL для хранения данных и Redis для кэширования, а также реализует механизмы пагинации и ограничения запросов (Rate Limiting).

## Стек технологий
- **Язык программирования:** Go
- **Фреймворк:** Gin
- **База данных:** PostgreSQL
- **Кэш:** Redis
- **Логирование:** Logrus
- **Документация:** Swagger

## Архитектурный подход
Проект построен на принципах **чистой архитектуры**, что позволяет:
- Четко разделять уровни приложения (handler, service, repository).
- Улучшать тестируемость кода.
- Обеспечивать гибкость и расширяемость системы.

### Структура проекта
```
/task-kaspi
│── cmd/
│   ├── server/
│       ├── main.go      # Точка входа
│── internal/
│   ├── handler/        # Контроллеры API
│   ├── service/        # Бизнес-логика
│   ├── repository/     # Работа с БД
│── pkg/
│   ├── database/       # Подключение к БД
│── docs/               # Swagger-документация
│── go.mod              # Модули Go
│── go.sum              # Хеши зависимостей
```

## Установка и запуск
### 1. Клонирование репозитория
```sh
git clone https://github.com/ddProgerGo/task-kaspi.git
cd task-kaspi
```
### 2. Запуск контейнеров PostgreSQL и Redis (docker-compose)
```sh
docker-compose up -d
```
### 3. Запуск сервера
```sh
go run cmd/server/main.go
```

## API
### 1. Получение списка людей по имени с пагинацией
**GET /people/info/phone/{name}?page=1&limit=10**
```json
Response:
{
    "data": [
        { "name": "John Doe", "iin": "123456789012", "phone": "+77011234567" }
    ],
    "page": 1,
    "limit": 10,
    "total": 50
}
```
### 2. Проверка ИИН
**GET /iin_check/{iin}**
```json
Response:
{ "valid": true }
```
### 3. Добавление нового человека
**POST /people/info**
```json
Request:
{
    "name": "John Doe",
    "iin": "123456789012",
    "phone": "+77011234567"
}
```
### 4. Получение человека по ИИН
**GET /people/info/iin/{iin}**
```json
Response:
{ "name": "John Doe", "iin": "123456789012", "phone": "+77011234567" }
```

## Доступ к Swagger UI
После запуска приложения документация доступна по адресу: ``` http://localhost:8080/swagger/index.html ```

## Оптимизация SQL-запросов
В таблице **people** были добавлены индексы для ускорения поиска:
```sql
CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_iin ON people(iin);
```
**Преимущества:**
- **idx_people_name** ускоряет поиск по `name ILIKE '%value%'`.
- **idx_people_iin** ускоряет поиск по `iin`.

## Кэширование
- Используется Redis для хранения данных о людях.
- При создании нового человека его ИИН удаляется из кеша, чтобы избежать устаревших данных.

## Валидация ИИН
- Валидация ИИН реализована на основе алгоритма, описанного в [Wikipedia](https://ru.wikipedia.org/wiki/%D0%98%D0%BD%D0%B4%D0%B8%D0%B2%D0%B8%D0%B4%D1%83%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B9_%D0%B8%D0%B4%D0%B5%D0%BD%D1%82%D0%B8%D1%84%D0%B8%D0%BA%D0%B0%D1%86%D0%B8%D0%BE%D0%BD%D0%BD%D1%8B%D0%B9_%D0%BD%D0%BE%D0%BC%D0%B5%D1%80):

## Заключение
Этот проект демонстрирует принципы чистой архитектуры, оптимизацию БД и производительность за счет кеширования. 🚀

