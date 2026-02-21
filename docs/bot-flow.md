# Telegram Flow

## Категории роутов

### 1) Onboarding пользователя

- `/start` -> `registerStartHandler`
- `name` -> `registerNameHandler`
- `OnText` (в режиме ожидания имени) -> `registerTextHandler`

### 2) Главное меню

- `back` -> `registerBackHandler`
- `record` -> `registerRecordHandler`

### 3) Навигация по неделе/дням

- `week_current` -> `registerWeekCurrentHandler`
- `week_next` -> `registerWeekNextHandler`
- `back_week` -> `registerBackWeekHandler`
- `day` -> `registerDayHandler`
- `back_days` -> `registerBackDaysHandler`

### 4) Пара и очередь

- `lesson` -> `registerLessonHandler`
- `join` -> `registerJoinHandler`
- `leave` -> `registerLeaveHandler`
- `back_lessons` -> `registerBackLessonsHandler`

### 5) Админские действия

- `admin_menu` -> `registerAdminMenuHandler`
- `reload` -> `registerReloadHandler`

## Группировка файлов

- `routes_user.go`
- `routes_schedule.go`
- `routes_queue.go`
- `routes_admin.go`

`RegisterRoutes()` в `controller.go` регистрирует группы в фиксированном порядке.

## Клавиатуры

Клавиатуры формируются в `internal/tgbot/keyboard.go`:
- `MainMenu`
- `WeekMenu`
- `DaysMenu`
- `LessonsMenu`
- `LessonActions`
- `AdminMenu`

