# events-store
Service for event processing and storage.

Task description [Link](https://github.com/pavlegich/events-store/blob/main/SPECIFICATION.md).

## Запуск

Вывести список всех возможных команд:

`make help`

Запустить приложение с помощью docker-compose:

`make run`

Запустить приложение локально:

`make run/local`

Сформировать документацию:

`make doc`

> [!NOTE]
> Изменить значения флагов:
> - для локального запуска - в первых строках Makefile;
> - для запуска с помощью docker-compose - в файле .env.

## .env Файл конфигурации

| Наименование параметра | Начальное значение | Описание |
| ---------------------- | ------------------ | -------- |
| `DB_PASSWORD` | `postgres` | Пароль для базы данных. |
| `DATABASE_DSN` | `postgresql://postgres:postgres@db:5432/postgres` | Строка подключения к базе данных. |
| `ADDRESS` | `:8080` | Адрес и порт, где будет запущено приложение. |

## Makefile Параметры запуска

| Наименование параметра | Начальное значение | Описание |
| ---------------------- | ------------------ | -------- |
| `DATABASE_DSN` | `postgresql://localhost:5432/postgres` | Строка подключения к базе данных. |
| `DOC_ADDR` | `localhost:6060` | Адрес и порт, где будет запущен сервис с документацией к приложению. |
| `SERVER_BINARY_NAME` | `server` | Наименование создаваемого бинарного файла для запуска приложения. |
| `SERVER_PACKAGE_PATH` | `./cmd/server` | Путь к бинарному файлу для запуска приложения. |
| `SERVER_ADDR` | `localhost:8080` | Адрес и порт, где будет запущено приложение. |
