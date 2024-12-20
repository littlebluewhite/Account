// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Wilson",
            "url": "https://github.com/littlebluewhite",
            "email": "wwilson008@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/account/user/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login"
                ],
                "summary": "login with username and password",
                "parameters": [
                    {
                        "description": "username and password",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/logs": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Log"
                ],
                "summary": "get logs history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "start time",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "stop time",
                        "name": "stop",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/e_log.Log"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "e_log.Log": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "api_url": {
                    "type": "string"
                },
                "content_length": {
                    "type": "integer"
                },
                "datetime": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "module": {
                    "type": "string"
                },
                "referer": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                },
                "timestamp": {
                    "type": "number"
                },
                "token": {
                    "type": "string"
                },
                "user_agent": {
                    "type": "string"
                },
                "web_path": {
                    "type": "string"
                }
            }
        },
        "user.Login": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0.0",
	Host:             "127.0.0.1:9600",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Schedule-Task-Command swagger API",
	Description:      "This is a account server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
