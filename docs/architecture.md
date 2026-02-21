# Архитектура

## Общая схема

Поток данных:
1. `cmd/main.go` создает конфиг и запускает `app.Start`.
2. `internal/app/bot.go` поднимает SQLite, создает `Store`, `Infra`, `Server`.
3. На старте вызывается `Infra.LoadDBScheduleItem()` для загрузки расписания.
4. Запускается Telegram-бот (`internal/tgbot/start.go`).
5. Обработчики из `Controller` читают/пишут данные через `Server`.

## Слои

### 1) Transport (`internal/tgbot`)

Отвечает за:
- Telegram callback/text handlers;
- формирование inline-клавиатур;
- пользовательский flow (неделя -> день -> пара -> очередь).

### 2) Service (`internal/server`)

Тонкий слой бизнес-логики:
- операции с пользователем;
- получение расписания на день;
- запись/выход из очереди.

### 3) Repository (`internal/repo/sqlLiteStore`)

Работа с SQLite:
- `users`
- `schedule_items`
- `records`

### 4) Infra/Parser (`internal/infra`, `internal/parser`)

- HTTP-запрос к API расписания;
- парсинг iCal контента;
- upsert в `schedule_items`.

## Идемпотентность загрузки расписания

В `LoadDBScheduleItem()` используется upsert:
- ключ конфликта: `(name, start_date, end_date)`;
- при конфликте выполняется `DO UPDATE`.

Это снижает риск дублей при повторной загрузке расписания.

## Таймзона

В docker-среде задается `TZ=Europe/Moscow`.

Это важно для:
- корректной навигации по дням;
- согласованного отображения расписания.

