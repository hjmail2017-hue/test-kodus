# CLAUDE.md

本文档为Claude Code (claude.ai/code) 在本代码库中工作时提供指导。

## 项目概述

这是一个用Go实现的JSON解析器，提供以下功能：
- JSON字符串的词法分析（标记化）
- 递归下降解析为Go数据结构
- 支持JSONPath查询以提取特定字段
- 用于解析和查询的命令行界面

`jsonparser` 包提供了简单的API：`Parse()` 用于解析JSON，`Get()` 用于JSONPath查询。

## 架构

### 核心组件

**词法分析器 (`lexer.go`)**
- 将JSON字符串转换为标记（`Token`结构体，包含Type、Literal、Pos字段）
- 处理Unicode字符和转义序列
- 在标记化过程中验证JSON数字格式
- 提供详细的错误位置信息

**解析器 (`parser.go`)**
- 带有前瞻（当前标记和下一个标记）的递归下降解析器
- 返回Go类型：`map[string]interface{}`（对象）、`[]interface{}`（数组）、`int64`、`float64`、`string`、`bool`、`nil`
- 验证JSON语法（逗号、冒号、括号）
- 主要入口点：`Parse(input string) (interface{}, error)`

**JSONPath (`jsonpath.go`)**
- 将JSONPath表达式解析为片段链
- 片段类型：`rootSegment`、`childSegment`、`arrayIndexSegment`
- 根据解析后的JSON数据评估路径
- 支持：`$.key`、`$[0]`、`$["key"]`、嵌套表达式
- 主要函数：`ParseJSONPath()`、`Evaluate()`、`Get()`快捷方式

### 数据流
1. `Parse()` → 词法分析器标记化 → 解析器树构建 → Go值
2. `Get()` → 解析JSONPath → 遍历解析后的数据 → 提取的值

### 命令行工具 (`cmd/jsonparser/main.go`)
- 从参数或标准输入解析JSON
- 可选的JSONPath参数用于字段提取
- 使用标准JSON编码器美化打印结果

## 开发命令

### 构建
```bash
# 构建命令行工具
go build ./cmd/jsonparser

# 构建包
go build .
```

### 测试
```bash
# 运行所有测试
go test ./...

# 运行特定测试套件
go test -v -run TestParseLiteral
go test -v -run TestParseArray
go test -v -run TestParseObject
go test -v -run TestParseErrors
go test -v -run TestParseJSONPath
go test -v -run TestJSONPathEvaluate

# 测试覆盖率
go test -cover
```

### 代码质量
```bash
# 格式化代码
go fmt ./...

# 静态分析
go vet ./...
```

## 关键实现细节

### 标记类型
词法分析器为所有JSON元素定义了标记类型：大括号、中括号、冒号、逗号、字符串、数字、true、false、null、EOF、错误。

### 数字验证
JSON数字会验证格式是否正确：可选负号、整数部分（不能有前导零，除非是单个0）、可选小数部分（必须有数字）、可选指数部分（e/E，可选符号和数字）。

### 转义序列
字符串支持：`\"`、`\\`、`\/`、`\b`、`\f`、`\n`、`\r`、`\t`、`\uXXXX`（简化）。

### JSONPath限制
- 不支持通配符操作符（`*`、`..`）
- 不支持过滤表达式
- 不支持数组切片
- 不支持联合操作符
- 不支持括号内的空格

## API使用示例

```go
// 解析JSON
data, err := jsonparser.Parse(`{"user": {"name": "Alice", "age": 30}}`)

// JSONPath查询
name, err := jsonparser.Get(data, "$.user.name")

// 可重用的JSONPath
jp, _ := jsonparser.ParseJSONPath("$.user.age")
age, _ := jp.Evaluate(data)
```

## 错误处理
解析和JSONPath评估都会返回带有位置的详细错误信息。解析器会验证完整的JSON结构（没有多余的标记）。

## 依赖项
纯Go实现，除了标准库外没有外部依赖。

## 开发规范
