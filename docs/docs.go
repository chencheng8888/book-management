// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/api/v1/book/stock/add": {
            "post": {
                "description": "添加库存接口，参数的where是可选参数",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "库存"
                ],
                "summary": "添加库存",
                "parameters": [
                    {
                        "description": "增加库存请求",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.AddStockReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.AddStockResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/stock/fuzzy_query": {
            "get": {
                "description": "模糊查询库存信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "库存"
                ],
                "summary": "模糊查询库存信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "入库时间",
                        "name": "add_stock_time",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "作者",
                        "name": "author",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "类别",
                        "name": "category",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "书本名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "每页大小",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.FuzzyQueryBookStockResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/stock/list": {
            "get": {
                "description": "列出所有库存信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "库存"
                ],
                "summary": "列出所有库存信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "每页大小",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ListBookStockResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/stock/searchByID": {
            "get": {
                "description": "根据ID查询库存信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "库存"
                ],
                "summary": "根据ID查询库存信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "书本ID",
                        "name": "book_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.SearchStockByBookIDResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.AddStockReq": {
            "type": "object",
            "required": [
                "author",
                "category",
                "name",
                "publisher",
                "quantity_added"
            ],
            "properties": {
                "author": {
                    "description": "作者",
                    "type": "string"
                },
                "category": {
                    "description": "类别",
                    "type": "string"
                },
                "name": {
                    "description": "书本名称",
                    "type": "string"
                },
                "publisher": {
                    "description": "出版社",
                    "type": "string"
                },
                "quantity_added": {
                    "description": "添加的库存数目",
                    "type": "integer"
                },
                "where": {
                    "description": "库存位置",
                    "type": "string"
                }
            }
        },
        "controller.AddStockResp": {
            "type": "object",
            "required": [
                "code",
                "msg"
            ],
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object",
                    "required": [
                        "book_id"
                    ],
                    "properties": {
                        "book_id": {
                            "description": "书本ID",
                            "type": "integer"
                        }
                    }
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "controller.Book": {
            "type": "object",
            "required": [
                "author",
                "book_id",
                "category",
                "created_at",
                "name",
                "publisher",
                "stock",
                "stock_status",
                "stock_where"
            ],
            "properties": {
                "author": {
                    "description": "作者",
                    "type": "string"
                },
                "book_id": {
                    "description": "书本ID",
                    "type": "integer"
                },
                "category": {
                    "description": "类别",
                    "type": "string"
                },
                "created_at": {
                    "description": "入库时间",
                    "type": "string"
                },
                "name": {
                    "description": "书本名称",
                    "type": "string"
                },
                "publisher": {
                    "description": "出版社",
                    "type": "string"
                },
                "stock": {
                    "description": "库存数量",
                    "type": "integer"
                },
                "stock_status": {
                    "description": "库存状态",
                    "type": "string"
                },
                "stock_where": {
                    "description": "库存位置",
                    "type": "string"
                }
            }
        },
        "controller.FuzzyQueryBookStockResp": {
            "type": "object",
            "required": [
                "code",
                "data",
                "msg"
            ],
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "description": "数据",
                    "type": "object",
                    "required": [
                        "books",
                        "current_page",
                        "total_page"
                    ],
                    "properties": {
                        "books": {
                            "description": "数据",
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controller.Book"
                            }
                        },
                        "current_page": {
                            "description": "当前页",
                            "type": "integer"
                        },
                        "total_page": {
                            "description": "总数",
                            "type": "integer"
                        }
                    }
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "controller.ListBookStockResp": {
            "type": "object",
            "required": [
                "code",
                "msg"
            ],
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object",
                    "required": [
                        "books",
                        "current_page",
                        "total_page"
                    ],
                    "properties": {
                        "books": {
                            "description": "数据",
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controller.Book"
                            }
                        },
                        "current_page": {
                            "description": "当前页",
                            "type": "integer"
                        },
                        "total_page": {
                            "description": "总数",
                            "type": "integer"
                        }
                    }
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "controller.SearchStockByBookIDResp": {
            "type": "object",
            "required": [
                "code",
                "data",
                "msg"
            ],
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "description": "数据",
                    "allOf": [
                        {
                            "$ref": "#/definitions/controller.Book"
                        }
                    ]
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8989",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Book Management API",
	Description:      "This is a sample server for a book management system.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
