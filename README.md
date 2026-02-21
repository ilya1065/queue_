# queue

Telegram-бот для очередей на пары.

Бот:
- загружает расписание из API МИРЭА;
- сохраняет его в SQLite;
- показывает расписание по дням;
- ведет очередь пользователей на выбранную пару.

## Технологии

- Go 1.25
- `gopkg.in/telebot.v4` (Telegram Bot API)
- SQLite (`modernc.org/sqlite`)
- `sqlx`
- Docker / Docker Compose

## Быстрый старт

### 1) Локально

1. Создай `.env` в корне:

```env
TG_KEY=<telegram_bot_token>
```

2. Примени миграции:

```bash
go run ./cmd/migrate
```

3. Запусти бота:

```bash
go run ./cmd
```

### 2) Через Docker Compose

1. Создай `.env` в корне:

```env
TG_KEY=<telegram_bot_token>
```

2. Запусти сервисы:

```bash
docker compose up -d --build
```

В `docker-compose.yml` используется отдельный volume `queue_queue_data` для БД.

## Конфигурация

Основной файл: `internal/config/config.yml`.

Поля:
- `db_url` — DSN SQLite.
- `tg_key` — есть в конфиге, но фактически токен берется из переменной окружения `TG_KEY`.

Переменные окружения:
- `TG_KEY` — обязателен для запуска бота.
- `DB_URL` — используется мигратором в Docker-сервисе `migrate`.
- `MIGRATIONS_DIR` — путь к SQL-миграциям для `cmd/migrate`.
- `TZ` — в docker-compose зафиксирован как `Europe/Moscow`.

## Миграции

Миграции лежат в `db/migrations`.

- `000001_create_users.up.sql` — базовая схема.
- `000002_dedupe_schedule_items.up.sql` — удаление дублей расписания и уникальный индекс по `(name, start_date, end_date)`.

Ручной запуск мигратора:

```bash
go run ./cmd/migrate
```

## Структура проекта

- `cmd/main.go` — точка входа приложения.
- `cmd/migrate/main.go` — простой SQL-мигратор.
- `internal/app` — сборка зависимостей и запуск.
- `internal/tgbot` — контроллер, маршруты и клавиатуры Telegram.
- `internal/server` — сервисный слой.
- `internal/repo/sqlLiteStore` — доступ к SQLite.
- `internal/infra` — загрузка расписания из внешнего API.
- `internal/parser` — парсинг iCal.
- `internal/entity` — доменные сущности и ошибки.
- `db/migrations` — SQL-схема и изменения.

## Дополнительная документация

- `docs/architecture.md` — архитектура и потоки данных.
- `docs/bot-flow.md` — сценарии и маршруты Telegram-бота.
- `docs/database.md` — схема БД и диагностические SQL-запросы.
- `docs/operations.md` — эксплуатация, запуск, отладка, частые проблемы.

