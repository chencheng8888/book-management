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
        "/api/v1/book/borrow/add": {
            "post": {
                "description": "借书接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "借书"
                ],
                "summary": "借书",
                "parameters": [
                    {
                        "description": "借书请求",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.BorrowBookReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.BorrowBookResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/borrow/query": {
            "get": {
                "description": "查询借书记录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "借书"
                ],
                "summary": "查询借书记录",
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
                    },
                    {
                        "type": "string",
                        "description": "借阅状态的查询条件",
                        "name": "query_status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.QueryBookBorrowRecordResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/borrow/query_statistics": {
            "get": {
                "description": "获取统计借阅记录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "借书"
                ],
                "summary": "获取统计借阅记录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "示例值 \"week\" \"month\" \"year\"",
                        "name": "pattern",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.QueryStatisticsBorrowRecordsResp"
                        }
                    }
                }
            }
        },
        "/api/v1/book/borrow/update_status": {
            "put": {
                "description": "更新借阅状态",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "借书"
                ],
                "summary": "更新借阅状态",
                "parameters": [
                    {
                        "description": "更新请求",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateBorrowStatusReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateBorrowStatusResp"
                        }
                    }
                }
            }
        },
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
                "description": "模糊查询库存信息,没有任何查询条件就是直接列出数据",
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
                        "description": "入库地点",
                        "name": "add_stock_where",
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
                        "in": "query"
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
        "controller.BookBorrowRecord": {
            "type": "object",
            "required": [
                "book_id",
                "copy_id",
                "return_status",
                "should_return_time",
                "user_id",
                "user_name"
            ],
            "properties": {
                "book_id": {
                    "description": "书本ID【这个你可以理解为一类书，比如《高等数学》】",
                    "type": "integer"
                },
                "copy_id": {
                    "description": "副本ID\t【这个你可以理解为具体一本书，比如《高等数学》的第一本",
                    "type": "integer"
                },
                "return_status": {
                    "description": "归还状态",
                    "type": "string"
                },
                "should_return_time": {
                    "description": "应该归还的时间",
                    "type": "string"
                },
                "user_id": {
                    "description": "用户ID",
                    "type": "string"
                },
                "user_name": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "controller.BorrowBookReq": {
            "type": "object",
            "required": [
                "book_id",
                "borrower_id",
                "expected_return_time"
            ],
            "properties": {
                "book_id": {
                    "description": "书本ID【这个你可以理解为一类书，比如《高等数学》】",
                    "type": "integer"
                },
                "borrower_id": {
                    "description": "借阅者ID",
                    "type": "string"
                },
                "expected_return_time": {
                    "description": "预计归还时间,格式为\"2006-01-02\"",
                    "type": "string"
                }
            }
        },
        "controller.BorrowBookResp": {
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
                    "type": "object",
                    "required": [
                        "book_id",
                        "copy_id"
                    ],
                    "properties": {
                        "book_id": {
                            "description": "实际借阅的书本ID",
                            "type": "integer"
                        },
                        "copy_id": {
                            "description": "副本ID\t【这个你可以理解为具体一本书，比如《高等数学》的第一本】",
                            "type": "integer"
                        }
                    }
                },
                "msg": {
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
        "controller.QueryBookBorrowRecordResp": {
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
                    "type": "object",
                    "required": [
                        "borrow_records",
                        "current_page",
                        "total_page"
                    ],
                    "properties": {
                        "borrow_records": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controller.BookBorrowRecord"
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
        "controller.QueryStatisticsBorrowRecordsResp": {
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
                        "art_enlightenment_num",
                        "children_story_num",
                        "science_knowledge_num"
                    ],
                    "properties": {
                        "art_enlightenment_num": {
                            "type": "integer"
                        },
                        "children_story_num": {
                            "type": "integer"
                        },
                        "science_knowledge_num": {
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
        },
        "controller.UpdateBorrowStatusReq": {
            "type": "object",
            "required": [
                "book_id",
                "copy_id",
                "status"
            ],
            "properties": {
                "book_id": {
                    "description": "书本ID【这个你可以理解为一类书，比如《高等数学》】",
                    "type": "integer"
                },
                "copy_id": {
                    "description": "副本ID\t【这个你可以理解为具体一本书，比如《高等数学》的第一本】",
                    "type": "integer"
                },
                "status": {
                    "description": "要更新的状态,取值有[waiting_return,returned,overdue]",
                    "type": "string"
                }
            }
        },
        "controller.UpdateBorrowStatusResp": {
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
                "data": {},
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
