# Go Tree-Walk Interpreter

A work-in-progress interpreter for a custom programming language, written in Go. It uses a tree-walk approach to scan,
parse, and evaluate source code.

## Features

- **Token Scanner**: Converts source code into tokens based on language grammar.
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens.
- **Interpreter**: Evaluates the AST, supporting arithmetic, logical operations, string manipulation, and explicit
  variable assignment.
- **Variable Assignment**: Supports updating variable values after declaration (e.g., `x = 2`).
- **Control Flow Signals**: Internal support for `break` and `continue` via control signal types.
- **Error Handling**: Reports runtime and syntax errors with line and character information.

## Usage

Clone the repository and run the interpreter (this run the default file in the :

```bash
git clone https://github.com/shahnawaz-lang/go-interpreter.git
cd go-interpreter
make run
