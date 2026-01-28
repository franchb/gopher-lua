# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build          # Build library (runs go-inline, go fmt, go build)
make test           # Run all tests
make glua           # Build standalone interpreter
go test -run TestX  # Run specific test
go test -v ./...    # Run all tests with verbose output
```

**Important**: Run `make build` before committing. This generates code from template files.

## Code Generation

Files prefixed with `_` are templates processed by `_tools/go-inline` (Python 3):
- `_state.go` → `state.go`
- `_vm.go` → `vm.go`

**Always edit the `_` prefixed template files, never the generated files.**

## Architecture

### Compilation Pipeline
```
Source → Lexer → Parser → AST → Compiler → Bytecode → VM
         parse/   parse/   ast/   compile.go         vm.go
         lexer.go parser.go
```

### Type System (value.go)
All Lua values implement `LValue` interface:
- `LNilType`, `LBool`, `LNumber` (float64), `LString`
- `LFunction` - Lua or Go function with optional upvalues
- `LTable` - hybrid array/hashmap with metatable support
- `LUserData` - wrapper for Go values
- `LState` - VM thread/coroutine
- `LChannel` - Go channel (GopherLua extension)

### Virtual Machine
- Register-based bytecode (32-bit instructions, 42 opcodes)
- `mainLoop()` - standard execution
- `mainLoopWithContext()` - context-aware with cancellation support
- Configurable registry (data stack) and call stack sizes

### Standard Libraries (linit.go)
base, package, table, io, os, string, math, debug, coroutine, channel

## Key Files

| Purpose | File |
|---------|------|
| VM state | `_state.go` (template), `state.go` (generated) |
| Bytecode execution | `_vm.go` (template), `vm.go` (generated) |
| Compiler | `compile.go` |
| Opcodes | `opcode.go` |
| Type system | `value.go` |
| Table impl | `table.go` |
| Function/Proto | `function.go` |
| Lexer | `parse/lexer.go` |
| Parser | `parse/parser.go` (generated from `parser.go.y`) |
| AST types | `ast/ast.go`, `ast/expr.go`, `ast/stmt.go` |

## Testing

- `_glua-tests/` - GopherLua-specific test scripts
- `_lua5.1-tests/` - Standard Lua 5.1 test suite
- `script_test.go` - Lua script test runner
- `*_test.go` - Go unit tests

## Fork Notes

This fork (github.com/franchb/gopher-lua) adds:
- Debug hooks support (including hook removal via `debug.sethook(nil, "crl")`)
- Performance optimizations (LNumber strconv, allocator improvements)
- Targets latest Go versions (currently Go 1.25) using all latest Go features and performance optimizations
