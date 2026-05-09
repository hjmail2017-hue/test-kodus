# JSON 解析器（Go 实现）

一个用 Go 语言编写的 JSON 解析器，支持完整的 JSON 语法规范，包含词法分析和语法解析功能。

## 功能特性

- ✅ 完整的 JSON 语法支持：
  - 基本类型：`null`、布尔值（`true`/`false`）、数字、字符串
  - 复合类型：对象（`{}`）、数组（`[]`）
  - 嵌套结构：对象嵌套、数组嵌套、混合嵌套
- ✅ 转义字符支持：`\"`、`\\`、`\/`、`\b`、`\f`、`\n`、`\r`、`\t`
- ✅ 数字格式支持：整数、浮点数、科学计数法（如 `1.2e3`）
- ✅ 严格的 JSON 规范验证
- ✅ 详细的错误信息（包含错误位置）
- ✅ JSONPath 查询支持（提取指定字段）
- ✅ 同时提供库 API 和命令行工具

## 安装

### 作为库使用

```bash
go get github.com/yourusername/json-parser
```

或者在 `go.mod` 中添加依赖：
```go
require json-parser v0.1.0
```

### 构建命令行工具

```bash
cd json-parser
go build ./cmd/jsonparser
```

## 使用方法

### 作为库使用

```go
package main

import (
    "fmt"
    "json-parser"
)

func main() {
    // 解析简单的 JSON 对象
    jsonStr := `{"name": "Alice", "age": 30, "active": true}`
    result, err := jsonparser.Parse(jsonStr)
    if err != nil {
        fmt.Printf("解析错误: %v\n", err)
        return
    }
    
    // result 是 interface{} 类型，根据 JSON 结构可能是：
    // map[string]interface{}  // 对象
    // []interface{}           // 数组
    // string                  // 字符串
    // int64                   // 整数
    // float64                 // 浮点数
    // bool                    // 布尔值
    // nil                     // null
    
    fmt.Printf("解析结果: %v\n", result)
    
    // 类型断言使用
    if obj, ok := result.(map[string]interface{}); ok {
        fmt.Printf("姓名: %v\n", obj["name"])
        fmt.Printf("年龄: %v\n", obj["age"])
    }
}
```

### 命令行工具

```bash
# 解析整个 JSON
./jsonparser '{"hello": "world", "count": 42}'

# 使用 JSONPath 提取字段
./jsonparser '{"hello": "world", "count": 42}' '$.hello'

# 提取嵌套字段
./jsonparser '{"user": {"name": "Alice", "age": 30}}' '$.user.name'

# 提取数组元素
./jsonparser '{"items": [1, 2, 3]}' '$.items[0]'

# 从标准输入读取
echo '{"key": "value"}' | ./jsonparser -

# 标准输入 + JSONPath
echo '{"data": {"value": 42}}' | ./jsonparser - '$.data.value'
```

输出示例：
```bash
$ ./jsonparser '{"user": {"name": "Alice", "age": 30}}' '$.user.name'
"Alice"

$ ./jsonparser '{"array": [1, 2, 3], "nested": {"foo": "bar"}}'
{
  "array": [
    1,
    2,
    3
  ],
  "nested": {
    "foo": "bar"
  }
}
```

## API 参考

### `Parse(input string) (interface{}, error)`

解析 JSON 字符串并返回对应的 Go 值。

**参数：**
- `input`：要解析的 JSON 字符串

**返回值：**
- `interface{}`：解析后的 Go 值
- `error`：解析错误，如果解析成功则为 `nil`

**返回类型映射：**
| JSON 类型 | Go 类型 |
|-----------|---------|
| object    | `map[string]interface{}` |
| array     | `[]interface{}` |
| string    | `string` |
| number (整数) | `int64` |
| number (浮点数) | `float64` |
| boolean   | `bool` |
| null      | `nil` |

### `Get(data interface{}, path string) (interface{}, error)`

使用 JSONPath 表达式从 JSON 数据中提取指定字段。

**参数：**
- `data`：解析后的 JSON 数据（通常来自 `Parse()` 的返回值）
- `path`：JSONPath 表达式（必须以 `$` 开头）

**返回值：**
- `interface{}`：提取的值
- `error`：JSONPath 解析或查询错误

**支持的 JSONPath 语法：**
- `$` - 根元素
- `$.key` - 子元素（点表示法）
- `$["key"]` - 子元素（括号表示法）
- `$[0]` - 数组索引
- `$.key1.key2[0].key3` - 组合表达式

### `ParseJSONPath(path string) (*JSONPath, error)`

解析 JSONPath 表达式并返回可重用的 `JSONPath` 对象。

### `(jp *JSONPath) Evaluate(data interface{}) (interface{}, error)`

使用已解析的 JSONPath 对象查询数据。

## 示例

### 基本示例

```go
// 解析字符串
result, _ := jsonparser.Parse(`"hello world"`)
// result: "hello world"

// 解析数字
result, _ := jsonparser.Parse(`42`)
// result: int64(42)

result, _ := jsonparser.Parse(`3.14`)
// result: 3.14

// 解析布尔值
result, _ := jsonparser.Parse(`true`)
// result: true

result, _ := jsonparser.Parse(`false`)
// result: false

// 解析 null
result, _ := jsonparser.Parse(`null`)
// result: nil
```

### 数组示例

```go
// 解析数组
result, _ := jsonparser.Parse(`[1, 2, "three", true]`)
// result: []interface{}{int64(1), int64(2), "three", true}

// 嵌套数组
result, _ := jsonparser.Parse(`[[1, 2], [3, 4]]`)
// result: []interface{}{
//     []interface{}{int64(1), int64(2)},
//     []interface{}{int64(3), int64(4)},
// }
```

### 对象示例

```go
// 解析对象
result, _ := jsonparser.Parse(`{"name": "Alice", "age": 30}`)
// result: map[string]interface{}{"name": "Alice", "age": int64(30)}

// 嵌套对象
result, _ := jsonparser.Parse(`{
    "user": {
        "name": "Bob",
        "preferences": {
            "theme": "dark",
            "notifications": true
        }
    },
    "tags": ["go", "json", "parser"]
}`)
```

### JSONPath 查询示例

```go
// 解析 JSON 数据
data := `{
    "store": {
        "book": [
            {
                "title": "Book 1",
                "price": 8.95
            },
            {
                "title": "Book 2", 
                "price": 12.99
            }
        ],
        "location": "Beijing"
    }
}`

parsed, _ := jsonparser.Parse(data)

// 提取单个字段
title, _ := jsonparser.Get(parsed, "$.store.book[0].title")
// title: "Book 1"

price, _ := jsonparser.Get(parsed, "$.store.book[1].price")
// price: 12.99

location, _ := jsonparser.Get(parsed, "$.store.location")
// location: "Beijing"

// 复用 JSONPath 对象
jp, _ := jsonparser.ParseJSONPath("$.store.book[0]")
book, _ := jp.Evaluate(parsed)
// book: map[string]interface{}{"title": "Book 1", "price": 8.95}
```

### 错误处理

```go
// 语法错误
result, err := jsonparser.Parse(`{"key": "value"`)  // 缺少右括号
if err != nil {
    fmt.Println(err)  // 输出: unexpected token EOF at position 14
}

// 无效数字
result, err := jsonparser.Parse(`12.`)  // 无效数字格式
if err != nil {
    fmt.Println(err)  // 输出: invalid number 12. at position 0
}

// 无效转义
result, err := jsonparser.Parse(`"\x"`)  // 无效转义序列
if err != nil {
    fmt.Println(err)  // 输出: unexpected token ERROR(invalid escape sequence \x) at position 2
}
```

## 测试

项目包含全面的测试用例：

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestParseLiteral
go test -v -run TestParseArray
go test -v -run TestParseObject
go test -v -run TestParseErrors

# 测试覆盖率
go test -cover

# 运行所有包的测试
go test ./...
```

测试覆盖以下场景：
- ✅ 字面量解析（null、布尔值、数字、字符串）
- ✅ 数组解析（空数组、数字数组、混合数组、嵌套数组）
- ✅ 对象解析（空对象、简单对象、多键对象、嵌套对象）
- ✅ 错误处理（语法错误、无效转义、格式错误）

## 实现细节

### 词法分析器（lexer.go）
- 将 JSON 字符串转换为 tokens
- 支持 Unicode 字符处理
- 实现转义字符解析
- 提供详细的错误位置信息

### 解析器（parser.go）
- 递归下降解析器设计
- 支持前瞻（lookahead）token
- 严格的 JSON 语法验证
- 类型安全的返回值

### 数字验证
```go
// JSON 数字格式规则：
// 1. 可选负号
// 2. 整数部分（不能以 0 开头，除非是单个 0）
// 3. 可选小数部分（必须包含至少一位数字）
// 4. 可选指数部分（e/E 后接可选符号和至少一位数字）
```

### 转义序列支持
```go
// 支持的转义序列：
// \"   双引号
// \\   反斜杠
// \/   斜杠
// \b   退格
// \f   换页
// \n   换行
// \r   回车
// \t   制表符
// \uXXXX Unicode 字符（支持 4 位十六进制）
```

### JSONPath 查询（jsonpath.go）
- **解析器**：将 JSONPath 表达式解析为段序列（根、子节点、数组索引）
- **段接口**：统一的操作接口，支持不同类型的选择器
- **查询引擎**：递归应用段操作，逐层提取数据
- **错误处理**：类型检查、边界检查、路径验证
- **语法支持**：点表示法（`.key`）、括号表示法（`["key"]`）、数组索引（`[n]`）

## 项目结构

```
json-parser/
├── go.mod                    # Go 模块定义
├── lexer.go                  # 词法分析器
├── parser.go                 # 解析器
├── parser_test.go            # 解析器测试用例
├── jsonpath.go               # JSONPath 查询实现
├── jsonpath_test.go          # JSONPath 测试用例
├── README.md                 # 本文档
├── cmd/
│   └── jsonparser/
│       └── main.go           # 命令行工具
└── jsonparser                # 可执行文件（构建后）
```

## 开发指南

### 添加新功能

1. **修改词法分析器**：在 `lexer.go` 中添加新的 token 类型
2. **修改解析器**：在 `parser.go` 中更新解析逻辑
3. **添加测试**：在 `parser_test.go` 中添加测试用例
4. **运行测试**：确保所有测试通过

### 代码规范

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 运行所有测试
go test ./...
```

## 限制

- 当前实现不支持大数字（超出 int64/float64 范围）
- Unicode 转义序列 `\uXXXX` 处理为简化版本
- 没有流式解析支持（适合小到中等大小 JSON）

## 许可证

MIT License