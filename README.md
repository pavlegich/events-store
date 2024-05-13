# events-store
Service for event processing and storage.

Task description [Link](https://github.com/pavlegich/events-store/blob/main/docs/SPECIFICATION.md).

## API

Для понимания работы с сервисом представлены:

- описание [API](https://github.com/pavlegich/events-store/blob/main/docs/api.yaml) в файле _docs/api.yaml_;

- [примеры](https://github.com/pavlegich/events-store/blob/main/docs/events-store.json) запросов для Postman в файле _docs/events-store.json_. 

## SQL Запросы

1. Выборка всех уникальных eventType у которых более 1000 событий.

`SELECT eventType FROM events GROUP BY eventType HAVING count(*) > 1000`

2. Выборка событий которые произошли в первый день каждого месяца.

`SELECT * FROM events WHERE toDayOfMonth(eventTime) == 1`

3. Выборка пользователей которые совершили более 3 различных eventType.

`SELECT userID FROM events GROUP BY userID HAVING count(DISTINCT eventType) > 3`

## Возникшие вопросы и описание решения

1. Пример даты из примера не соответствует существущим в Go типам для кодирования / декодирования из JSON. Написаны структура с переменной типа time.Time и методы, реализующие Marshaler и Unmarshaler интерфейсы. Также добавлен метод Scan, реализующий метод интерфейса Scanner для чтения и декодирования переменной из базы данных.

2. Скрипты для создания и миграций таблиц находятся в директории _migrations_.

3. [Пример](https://github.com/pavlegich/events-store/blob/main/configs/.env.example) файла конфигурации расположен по пути _configs/.env.example_.

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
| `ADDRESS` | `:8080` | Адрес и порт, где будет запущено приложение. |
| `DB_PASSWORD` | `default` | Пароль для базы данных. |
| `DB_USER` | `default` | Имя пользователя для базы данных. |
| `DB_NAME` | `events` | Наименование базы данных. |
| `DB_HOST` | `db` | Наименование хоста для базы данных. |
| `GOOSE_DRIVER` | | Драйвер, используемый Goose для миграции базы данных. |
| `GOOSE_DBSTRING` | | Строка подключения, используемая Goose для миграции базы данных. |
| `DATABASE_DSN` | | Строка подключения к базе данных. |

## Makefile Параметры запуска

| Наименование параметра | Начальное значение | Описание |
| ---------------------- | ------------------ | -------- |
| `SERVER_BINARY_NAME` | `server` | Наименование создаваемого бинарного файла для запуска приложения. |
| `SERVER_PACKAGE_PATH` | `./cmd/server` | Путь к бинарному файлу для запуска приложения. |
| `SERVER_ADDR` | `localhost:8080` | Адрес и порт, где будет запущено приложение. |
| `DB_ADDR` | `localhost:9000` | Адрес и порт для подключения к базе данных. |
| `DB_NAME` | `events` | Наименование базы данных. |
| `DATABASE_DSN` | | Строка подключения к базе данных. |
| `DOC_ADDR` | `localhost:6060` | Адрес и порт, где будет запущен сервис с документацией к приложению. |
