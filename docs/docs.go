// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Михаил Токмачев",
            "url": "https://t.me/Zatrasz",
            "email": "zatrasz@ya.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "post": {
                "description": "Принимает обязательные поля group , song .",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Добавление новой песни",
                "parameters": [
                    {
                        "description": "Данные структуры песни",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Songs"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно созданная запись \u003c Ok \u003e",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса или не верно заполнены обязательные поля",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Не получены детальные данные из API или ошибка при сохранении в бд",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.SongDetail": {
            "type": "object",
            "properties": {
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.Songs": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "detail": {
                    "$ref": "#/definitions/models.SongDetail"
                },
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8586",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Реализация онлайн библиотеки песен",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
