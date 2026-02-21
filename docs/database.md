# База Данных

## Движок

SQLite (`modernc.org/sqlite`).

Рекомендуемый DSN (по умолчанию в проекте):

```text
file:./data/app.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)
```

## Таблицы

### `users`

- `id INTEGER PRIMARY KEY` (Telegram user id)
- `name TEXT NOT NULL`

### `schedule_items`

- `id INTEGER PRIMARY KEY`
- `name TEXT NOT NULL`
- `description TEXT`
- `start_date timestamp NOT NULL`
- `end_date timestamp NOT NULL`
- `external_id TEXT NOT NULL`

Ограничения:
- историческое: `UNIQUE (external_id, start_date)`
- дополнительный индекс: `idx_schedule_items_name_start_end` на `(name, start_date, end_date)`

### `records`

- `id INTEGER PRIMARY KEY`
- `user_id INTEGER NOT NULL`
- `schedule_item_id INTEGER NOT NULL`
- `createdAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP`

Ограничения:
- `FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- `FOREIGN KEY (schedule_item_id) REFERENCES schedule_items(id) ON DELETE CASCADE`
- `UNIQUE (user_id, schedule_item_id)`

## Полезные SQL-запросы

### Количество пар в БД

```sql
SELECT COUNT(*) FROM schedule_items;
```

### Проверка дублей по естественному ключу пары

```sql
SELECT name, start_date, end_date, COUNT(*) AS c
FROM schedule_items
GROUP BY name, start_date, end_date
HAVING c > 1
ORDER BY c DESC;
```

### Расписание на конкретный день

```sql
SELECT id, name, start_date, end_date
FROM schedule_items
WHERE substr(start_date, 1, 10) = '2026-02-25'
ORDER BY start_date;
```

### Проверка целостности очереди

```sql
SELECT r.id, r.user_id, r.schedule_item_id
FROM records r
LEFT JOIN users u ON u.id = r.user_id
LEFT JOIN schedule_items s ON s.id = r.schedule_item_id
WHERE u.id IS NULL OR s.id IS NULL;
```

