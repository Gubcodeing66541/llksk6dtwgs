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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/service/auth/login": {
            "post": {
                "tags": [
                    "客服信息"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "账号",
                        "name": "member",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service/create_live_token/2": {
            "post": {
                "tags": [
                    "客服信息"
                ],
                "summary": "创建视频通话",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客服名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客服头像",
                        "name": "head",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/api/service/info": {
            "post": {
                "tags": [
                    "客服信息"
                ],
                "summary": "获取客服基本信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service/rooms/detail": {
            "post": {
                "tags": [
                    "房间信息"
                ],
                "summary": "房间-获取用户房间详细",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service/rooms/list": {
            "post": {
                "tags": [
                    "房间信息"
                ],
                "summary": "房间-获取用户房间列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "user_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "all 所有 user_no_read  用户未读 server_read 已回复 server_no_read 未回复 top 置顶 black 拉黑",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "指定页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "指定每页数量",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service/update": {
            "post": {
                "tags": [
                    "客服信息"
                ],
                "summary": "修改客服头像和昵称",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客服名称",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客服头像",
                        "name": "head",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/api/service_message/create": {
            "post": {
                "tags": [
                    "快捷消息"
                ],
                "summary": "快捷消息-创建",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "消息内容",
                        "name": "msg_info",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "消息类型 text文本 image图片 video视频 link链接",
                        "name": "msg_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "类型 hello打招呼 quick_reply快捷回复 leave离线消息",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service_message/delete": {
            "post": {
                "tags": [
                    "快捷消息"
                ],
                "summary": "快捷消息-删除",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "删除的消息指定ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service_message/list": {
            "post": {
                "tags": [
                    "快捷消息"
                ],
                "summary": "快捷消息-列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "条目数",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "分页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "类型 hello打招呼 quick_reply快捷回复 leave离线消息",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service_message/swap": {
            "post": {
                "tags": [
                    "快捷消息"
                ],
                "summary": "快捷消息-位置交换",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "来自于的交换ID",
                        "name": "form",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "来给予交换ID",
                        "name": "to",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/service_message/update": {
            "post": {
                "tags": [
                    "快捷消息"
                ],
                "summary": "快捷消息-修改",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "修改ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "消息内容",
                        "name": "msg_info",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "消息类型 text文本 image图片 video视频 link链接",
                        "name": "msg_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "类型 hello打招呼 quick_reply快捷回复 leave离线消息",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "enable 启用 un_enable 禁用",
                        "name": "status",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/system/upload": {
            "post": {
                "tags": [
                    "公共接口"
                ],
                "summary": "系统默认文件上传",
                "parameters": [
                    {
                        "type": "string",
                        "description": "认证token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "文件参数",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0`",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "客服系统API文档`",
	Description:      "客服系统api `\n客服后端：`",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
