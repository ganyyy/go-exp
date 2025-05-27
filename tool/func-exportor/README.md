# func-exportor

# func-exportor

A powerful Go tool for analyzing Go source files and extracting exported symbols. This tool can generate header-like files, provide comprehensive package overviews, and export function signatures in various formats.

## Features

🔍 **Symbol Analysis**
- Extract all exported symbols (functions, variables, constants, types)
- Support for methods and interfaces
- Documentation comment extraction
- Position tracking for all symbols

📄 **Header File Generation**
- Generate Go header files similar to C header files
- Convert functions to function variable declarations
- Preserve package imports and structure
- Include documentation comments

📊 **Multiple Output Formats**
- Human-readable format with emojis
- JSON format for programmatic use
- Statistics summary
- Comprehensive package overview

✅ **Validation**
- Validate generated header files compile correctly
- Verbose output for debugging
- Error handling and recovery

## Installation

```bash
# Clone or download the tool
cd func-exportor
go build -o func-exportor main.go
```

## Usage

### Basic Usage

```bash
# Analyze a Go file and show all exported symbols
./func-exportor main.go

# Show statistics only
./func-exportor --stats main.go

# Output in JSON format
./func-exportor --json main.go

# Generate comprehensive package overview
./func-exportor --overview main.go
```

### Header File Generation

```bash
# Generate header file
./func-exportor --header header.go main.go

# Generate and validate header file
./func-exportor --header header.go --validate main.go

# Verbose output during generation
./func-exportor --header header.go --verbose main.go
```

### Command Line Options

| Option | Description |
|--------|-------------|
| `--json` | Output results in JSON format |
| `--stats` | Show only statistics summary |
| `--header <file>` | Generate header file with function signatures |
| `--validate` | Validate that generated header file compiles |
| `--verbose` | Enable verbose output for debugging |
| `--overview` | Generate comprehensive package overview with documentation |
| `--version` | Show version information |

## Examples

### Example 1: Basic Analysis

```bash
./func-exportor example.go
```

Output:
```
在文件 example.go 中找到 14 个导出的符号:

Function (4个):
  - NewUser (example.go:41:1)
  - ProcessRequest (example.go:49:1)
  - String (example.go:57:1)
  - Validate (example.go:61:1)

Type (4个):
  - User (example.go:24:6)
  - Handler (example.go:31:6)
  - Config (example.go:35:6)
  - Response (example.go:72:6)

Constant (3个):
  - DefaultTimeout (example.go:11:2)
  - MaxRetries (example.go:12:2)
  - Version (example.go:13:2)

Variable (3个):
  - GlobalConfig (example.go:18:2)
  - Logger (example.go:19:2)
  - DefaultUser (example.go:20:2)
```

### Example 2: Package Overview

```bash
./func-exportor --overview example.go
```

Output:
```
Package Overview for example.go
===============================

📊 Summary
----------
Total exported symbols: 14
  Constants: 3
  Variables: 3
  Types: 4
  Functions: 4

🏗️ Types (4)
---------------
  • User - 类型声明 (example.go:24:6)
    type User struct {
        ID       int64  `json:"id"`
        Name     string `json:"name"`
        Email    string `json:"email"`
        CreateAt time.Time
    }

⚡ Functions (4)
-------------------
  • NewUser - 函数 (example.go:41:1)
    var NewUser func(name string, email string) *User
```

### Example 3: Header File Generation

```bash
./func-exportor --header user_header.go user.go
```

This generates a header file like:
```go
package user

import (
    "context"
    "time"
)

// Types
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    CreateAt time.Time
}

// Constants
const DefaultTimeout = 30 * time.Second

// Variables  
var GlobalConfig map[string]interface{}

// Functions (as function variables)
var NewUser func(name string, email string) *User
var ProcessRequest func(ctx context.Context, data []byte) (*Response, error)
```

## How It Works

### Symbol Detection
The tool uses Go's AST (Abstract Syntax Tree) parser to analyze source code and identify:
- Exported functions (including methods)
- Exported types (structs, interfaces)
- Exported constants and variables
- Documentation comments

### Function Transformation
Functions are converted to function variable declarations:
- `func NewUser(name string) *User` becomes `var NewUser func(name string) *User`
- Methods are converted to functions with receiver as first parameter
- Return types and parameters are preserved exactly

### Header File Structure
Generated header files maintain:
- Original package declaration
- All necessary imports
- Organized sections: Types, Constants, Variables, Functions
- Documentation comments
- Proper Go syntax

## Testing

Run the comprehensive test suite:

```bash
./test_all.sh
```

This tests all major functionality including:
- Symbol extraction
- JSON output
- Statistics mode
- Package overview
- Header generation and validation

## Version History

### v1.1.0
- ✅ Added package overview functionality with emojis
- ✅ Integrated all command-line options
- ✅ Enhanced documentation extraction
- ✅ Fixed function return type handling
- ✅ Added comprehensive testing suite

### v1.0.0
- ✅ Basic symbol extraction
- ✅ Header file generation
- ✅ JSON and statistics output
- ✅ Validation functionality

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is provided as-is for educational and development purposes.

// 非导出符号（不会被检测）
const internal = "secret"
func helper() {}
```

运行分析：
```bash
./func-exportor demo.go
```

## 应用场景

- 📝 API文档生成
- 🔍 代码审查和分析
- 📊 代码质量度量
- 🛠️ IDE插件开发
- 📈 项目重构分析

## 许可证

MIT License
