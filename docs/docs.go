// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v0/auth/login": {
            "post": {
                "description": "sign in account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in",
                "parameters": [
                    {
                        "description": "login params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/v0/auth/logout": {
            "post": {
                "description": "logout from account and destroy the session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/v0/auth/register": {
            "post": {
                "description": "create an account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up",
                "parameters": [
                    {
                        "description": "register params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/v0/notes": {
            "get": {
                "description": "returning all notes of current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Fetch all notes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/notes.UserNotesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Create note",
                "parameters": [
                    {
                        "description": "create note params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/notes.CreateNoteRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/notes.NoteResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Update note",
                "parameters": [
                    {
                        "description": "update note params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/notes.UpdateNoteRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/notes.NoteResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        },
        "/v0/notes/{id}": {
            "get": {
                "description": "returning specific note",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Fetch note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Note ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/notes.NoteResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Delete note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Note ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3,
                    "example": "login"
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8,
                    "example": "password"
                }
            }
        },
        "auth.RegisterRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 3,
                    "example": "login"
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8,
                    "example": "password"
                }
            }
        },
        "note.Note": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "string",
                    "example": "07f3c5a1-70ea-4e3f-b9b5-110d29891673"
                },
                "content": {
                    "type": "string",
                    "example": "some content"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-10-03T14:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "6b30e5df-5add-42e1-be60-62b6f98afab1"
                },
                "title": {
                    "type": "string",
                    "example": "some title"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-10-04T04:20:00Z"
                }
            }
        },
        "notes.CreateNoteRequest": {
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "example": "some content"
                },
                "title": {
                    "type": "string",
                    "maxLength": 32,
                    "example": "some title"
                }
            }
        },
        "notes.NoteResponse": {
            "type": "object",
            "properties": {
                "note": {
                    "$ref": "#/definitions/note.Note"
                }
            }
        },
        "notes.UpdateNoteRequest": {
            "type": "object",
            "required": [
                "content",
                "note_id",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "example": "new content"
                },
                "note_id": {
                    "type": "string",
                    "example": "07f3c5a1-70ea-4e3f-b9b5-110d29891673"
                },
                "title": {
                    "type": "string",
                    "maxLength": 32,
                    "example": "new title"
                }
            }
        },
        "notes.UserNotesResponse": {
            "type": "object",
            "properties": {
                "notes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/note.Note"
                    }
                }
            }
        },
        "responses.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 500
                },
                "developer": {},
                "message": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Notes Service API",
	Description:      "Simple server to demonstrate some features",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}