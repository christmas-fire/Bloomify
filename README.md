# 🌸 Bloomify

Bloomify - это RESTful API сервис для управления цветочным магазином, разработанный на Go.

## Функционал

- **Авторизация/Аутентификация пользователей**:
  - Регистрация нового пользователя.
  - Вход в систему (аутентификация) и получение JWT токена.
  
- **CRUD операции с сущностями пользователей и цветов**:
  - Создание, чтение, обновление и удаление пользователей.
  - Создание, чтение, обновление и удаление цветов.
  - Поиск цветов по имени, описанию, цене и наличию.

## 🛠 Технологический стек

- **Go** - основной язык программирования.
- **PostgreSQL** - основная база данных.
- **Gin** - фреймворк для создания HTTP сервера и маршрутизации запросов.
- **SQLx** - библиотека для работы с PostgreSQL.
- **Swagger** - генерация и документация API.
- **Mockgen** - генерация моков для тестирования.

## Установка зависимостей

Убедитесь, что у вас установлен Go и PostgreSQL. Затем выполните следующие команды:

### Клонируйте репозиторий:
```bash
git clone https://github.com/christmas-fire/Bloomify.git
cd Bloomify
```

### Установите зависимости:
```bash
go mod download
```

### Настройте конфигурацию базы данных в файле `.env`.

## Конфигурация

Создайте файл `.env` в корне проекта и добавьте следующие переменные окружения:

```yaml
DATABASE_HOST=db
DATABASE_PORT=5432
DATABASE_USER=your_db_user
DATABASE_PASSWORD=your_db_password
DATABASE_NAME=bloomify

POSTGRES_USER=your_db_user
POSTGRES_PASSWORD=your_db_password
POSTGRES_DB=bloomify
```

## Локальный запуск
Для локального запуска приложения используйте следующие команды:

### Запуск приложения :
```bash
make run
```
### Запуск тестов :
```bash
make test
```

### Проверка кодстайла :
```bash
make lint
```

### Форматирование кода :
```bash
make fmt
```

## Документация
Документация API доступна по ссылке [interactive Swagger documentation](http://localhost:8080/swagger/index.html#/) после запуска приложения.
