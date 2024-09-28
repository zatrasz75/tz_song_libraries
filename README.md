## **Приступая к работе**

### **Предварительные требования**
- Перейти на версию Golang 1.22 или более позднюю
- PostgreSQL latest или более поздней версии

### **Установка**

**Клонировать репозиторий**\:

```
https://github.com/zatrasz75/tz_song_libraries.git
cd tz_song_libraries
```
- Если необходимо отредактируйте файл .env
```azure
APP_IP="localhost"
APP_PORT="8586"
CORS_ALLOWED_ORIGINS="http://localhost:8586;https://api.example.com"

# Хост для обагащения данных из внешнего api
URL_EXTERNAL="https://api.example.com"

DB_CONNECTION_STRING="postgres://zatrasz:postgrespw@localhost:49878/db_song?sslmode=disable"
```
- Файл миграции находится в migrations
- Для создания новых файлов и запуска приложения используйте Makefile
```azure
run:
	go run cmd/main.go

up:
	sql-migrate new up

down:
	sql-migrate down

swag:
	swag init -d internal/handlers/ -g router.go --parseDependency --parseDepth 3
```
- По умолчанию приложение:
```azure
[INFO] postgres.go:83 Применена 0 миграция!
[INFO] app.go:51 Запуск сервера на http://localhost:8586
[INFO] app.go:52 Документация Swagger API: http://localhost:8586/swagger/index.html
```

#### Примечание
**При обогащении текста песен из внешнего API , думаю будет выглядеть так :**

- Иначе не знаю как поведет себя получение текста песни с пагинацией по куплетам
```azure
Ooh, baby, don't you know I suffer?
Ooh, baby, can't you hear me moan?
You caught me under false pretenses
How long before you let me go?

Ooh, you set my soul alight
Ooh, you set my soul alight

I thought I was a fool for no one
But ooh, baby I'm a fool for you
You're the queen of the superficial
And how long before you tell the truth?

Ooh, you set my soul alight
Ooh, you set my soul alight
Ooh, you set my soul alight
Ooh, you set my soul

Supermassive black hole
Supermassive black hole
Supermassive black hole
Supermassive black hole

Ooh, baby, can't you hear me moan?
Ooh, baby, can't you hear me moan?
```
### REST методы:
- /songs [post] Добавление новой песни
- /songs [patch] Изменение данных песни из библиотеки по ID
- /songs [get] Получение данных с фильтрацией по всем полям и пагинацией
- /songs [del] Удаление песни из библиотеки по ID
- /songs/lyrics [get] Получение текста песни с пагинацией по куплетам

### Задание:
Необходимо реализовать следующее

1. Выставить rest методы
2. Получение данных библиотеки с фильтрацией по всем полям и пагинацией
3.   Получение текста песни с пагинацией по куплетам
4.   Удаление песни
5.   Изменение данных песни
6.   Добавление новой песни в формате

JSON
```
{
"group": "Muse",
"song": "Supermassive Black Hole"
}
```

2. При добавлении сделать запрос в АПИ, описанного сваггером

```
openapi: 3.0.3
info:
  title: Music info
  version: 0.0.1
paths:
  /info:
    get:
      parameters:
        - name: group
          in: query
          required: true
          schema:
            type: string
        - name: song
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
                  description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SongDetail'
        '400':
          description: Bad request
        '500':
          description: Internal server error
components:
  schemas:
    SongDetail:
      required:
        - releaseDate
        - text
        - link
      type: object
      properties:
        releaseDate:
          type: string
          example: 16.07.2006
        text:
          type: string
          example: Ooh baby, don't you know I suffer?
        patronymic:
          type: string
          example: https://www.youtube.com/watch?v=Xsp3_a-PMTw

```