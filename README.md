# biocad-tsv-service

---
Сервис для:
- Периодического сканирования директории на наличие .tsv файлов 
- Парсинга данных и сохранения в PostgreSQL 
- Логирования ошибок парсинга 
- Генерации PDF отчётов по unit_guid 
- Предоставления HTTP API для получения сообщений с пагинацией.

---
## Архитектура
```shell
.
├── Dockerfile
├── README.md
├── cmd
│   └── app
│       └── main.go
├── config.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
├── input
│   └── test_data_01.tsv
├── internal
│   ├── api
│   │   └── server.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── postgres.go
│   ├── migrations
│   │   ├── 001_create_messages.sql
│   │   ├── 002_create_processed_files.sql
│   │   └── 003_create_parse_errors.sql
│   ├── models
│   │   ├── message.go
│   │   ├── parse_error.go
│   │   └── processed_file.go
│   ├── parser
│   │   └── tsv_parser.go
│   ├── pdf
│   │   └── pdf.go
│   ├── queue
│   │   ├── manager.go
│   │   └── scanner.go
│   ├── repository
│   │   ├── message_repo.go
│   │   ├── parse_error_repo.go
│   │   └── processed_file_repo.go
│   └── util
│       └── util.go
└── output
```

---
## Конфигурация

Файл `config.yaml`:
```yaml
server:
port: "8080"

db:
host: "db"
port: 5432
user: "postgres"
password: "postgres"
name: "tsv_service"

dirs:
input: "./input"
output: "./output"
```

---
## Запуск через Docker
### Сборка и запуск
```shell
docker compose up --build
```
Сервис будет доступен на:

http://localhost:8080

PostgreSQL:
```shell
localhost:5432
user: postgres
password: postgres
db: tsv_service
```

---
## Работа сервиса
### 1. Сканирование

Каждые 30 секунд сервис:
- ищет `.tsv` файлы в папке `input`
- проверяет, обрабатывался ли файл ранее 
- добавляет в очередь на обработку

### 2. Парсинг

Каждая строка TSV:
- валидируется 
- сохраняется в таблицу `messages`
- при ошибке — записывается в `parse_errors`

### 3. Генерация PDF

После обработки файла:
- для каждого `unit_guid`
- создаётся PDF отчёт 
- сохраняется в папку `output`

---
## API
`GET /messages`

Получение сообщений по `unit_guid` с пагинацией.

**Параметры:**

| Параметр  | Обязательный | Описание                                 |
| --------- | ------------ | ---------------------------------------- |
| `unit_guid` | ✅            | UUID устройства                          |
| `page`      | ❌            | номер страницы (по умолчанию 1)          |
| `limit`     | ❌            | размер страницы (1–100, по умолчанию 50) |

Пример запроса:
```shell
GET http://localhost:8080/messages?unit_guid=11111111-1111-1111-1111-111111111111&page=1&limit=20
```

Ответ:
```shell
{
  "page": 1,
  "limit": 20,
  "total": 134,
  "data": [
    {
      "id": "...",
      "unit_guid": "...",
      "text": "...",
      ...
    }
  ]
}
```

---
## Структура БД
- **`messages`** – хранит успешно распарсенные сообщения.
- **`parse_errors`** – ошибки парсинга (битые строки).
- **`processed_files`** – статус обработки файлов.

---
## Graceful Shutdown

При `SIGTERM`/`SIGINT`:
- останавливается scanner 
- завершаются worker’ы 
- корректно закрывается HTTP сервер 
- закрывается пул соединений к БД

---
## Тестирование

Пример тестового файла:

```tsv
mqtt	unit_guid	msg_id	text	context	class	level	area	addr	block	type	bit	invert_bit
topic/alpha	11111111-1111-1111-1111-111111111111	1001	Test message	System	Alarm	2	A1	10	B1	sensor	1	false
```

Положить в папку: `input/test.tsv`.

---