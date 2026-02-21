# Эксплуатация И Отладка

## Локальный запуск

```bash
go run ./cmd/migrate
go run ./cmd
```

Перед запуском убедись, что в окружении задан `TG_KEY`.

## Запуск через Docker Compose

```bash
docker compose up -d --build
```

Проверка статуса:

```bash
docker compose ps
```

Логи приложения:

```bash
docker compose logs -f queue
```

## Подключение к БД в Docker volume

Имя volume обычно: `queue_queue_data`.

Просмотр таблиц:

```bash
docker run --rm -it -v queue_queue_data:/data alpine sh -lc \
"apk add --no-cache sqlite >/dev/null && sqlite3 /data/app.db '.tables'"
```

## Типовые проблемы

### 1) `TG_KEY is not set`

Причина: не передан токен бота.

Решение:
- добавить `TG_KEY` в `.env`;
- перезапустить сервис.

### 2) `FOREIGN KEY constraint failed` при `join`

Возможные причины:
- пользователь еще не создан в `users`;
- `schedule_item_id` устарел (нажата старая кнопка из старого сообщения).

Решение:
- пройти `/start` и указать имя;
- заново открыть актуальное расписание и выбрать пару из нового сообщения.

### 3) Расхождение локальной и docker-БД

Причина:
- локально используется `./data/app.db`;
- в docker используется volume `queue_queue_data`.

Решение:
- проверять нужный файл/volume отдельно;
- при необходимости сбросить docker state:

```bash
docker compose down -v
docker compose up -d --build
```

### 4) Неверный день в расписании

Проверить:
- `TZ=Europe/Moscow` в docker-compose;
- SQL выборку по дате через `substr(start_date, 1, 10)`.

## Обновление расписания вручную

В боте (для админа) кнопка:
- `админ панель` -> `Обновить расписание`

ID админа сейчас захардкожен в `internal/tgbot/keyboard.go`.

