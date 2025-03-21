# Go (Tree-Walk) Interpreter

This project is an implementation of a Tree-Walk interpreter written in Go. The interpreter processes source code by scanning, parsing, and evaluating expressions, supporting basic arithmetic, logical operations, and string manipulation. A lot more to come soon!

### Features
- **Token Scanner**: Converts source code into tokens based on language grammar.
- **Parser**: Constructs an Abstract Syntax Tree (AST) from tokens.
- **Interpreter**: Evaluates the AST to produce results.

### Current Status
This project is a **work in progress** and is far from complete. Many features are still under development, and the codebase is subject to significant changes. 

### Future Plans
- Add support for more complex expressions and statements.
- Implement error recovery for better user feedback.
- Expand the language's grammar to include loops, conditionals, and functions.
- Improve performance and add comprehensive tests.

### Getting Started
To run the interpreter, clone the repository and build the project using Go:
```bash
git clone https://github.com/your-username/go-interpreter.git
cd go-interpreter
make run
```