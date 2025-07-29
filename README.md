# Go (Tree-Walk) Interpreter

This project is an implementation of a Tree-Walk interpreter written in Go. The interpreter processes source code by scanning, parsing, and evaluating expressions, supporting basic arithmetic, logical operations, and string manipulation. A lot more to come soon!

### Features
- **Token Scanner**: Converts source code into tokens based on language grammar.
- **Parser**: Constructs an Abstract Syntax Tree (AST) from tokens.
- **Interpreter**: Evaluates the AST to produce results.

### Current Status
This project is a **work in progress** and is far from complete. Many features are still under development, and the codebase is subject to significant changes. 

## Features

- **Token Scanner**: Converts source code into tokens.
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens.
- **Interpreter**: Evaluates the AST, supporting arithmetic, logical operations, and string manipulation.
- **Error Handling**: Reports runtime and syntax errors with line and character information.

## Usage

Clone the repository and run the interpreter:

```bash
git clone https://github.com/shahnawaz-lang/go-interpreter.git
cd go-interpreter
make run
```