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