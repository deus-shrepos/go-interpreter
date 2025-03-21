package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Grammar = []string{
	"Binary: Left Expr,Operator internal.Token,Right Expr",
	"Grouping: Expression Expr",
	"Literal: Value interface{}", // TODO: enforce a type constraint. Obviously, I don't want just any type.
	"Unary: Operator internal.Token,Right Expr",
}

// GenerateAst This function will generate an AST .go file for each operation
func GenerateAst(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: generate_ast <output directory>")
	}
	outputDir := args[0]
	// types => (arg type)
	err := defineAst(outputDir, "Expr", Grammar)
	if err != nil {
		return err
	}
	return nil
}

// defineAst This function will be used to create AST structs for types of operations
func defineAst(outputDir string, baseName string, types []string) error {
	path := strings.Join([]string{outputDir, "/", baseName, ".go"}, "")
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString("package ast") // TODO: do I want to keep it as main?
	_, _ = writer.WriteString("\n\n")
	_, _ = writer.WriteString("import \"crafting-interpreters/internal\"")
	_, _ = writer.WriteString("\n\n")
	_, _ = writer.WriteString(strings.Join([]string{"type", " ", baseName, " ", "interface{}"}, ""))
	_, _ = writer.WriteString("\n")
	defineVisitor(writer, baseName, types)
	for _, type_ := range types {
		structName := strings.Trim(strings.Split(type_, ":")[0], " ")
		structFields := strings.Trim(strings.Split(type_, ":")[1], " ")
		defineType(writer, structName, structFields)
	}
	_ = writer.Flush()
	return nil
}

// defineType Define all the Production symbols with their properties
func defineType(writer *bufio.Writer, structName, structFields string) {
	_, _ = writer.WriteString(fmt.Sprintf("type %s struct {\n\t", structName))
	_, _ = writer.WriteString(strings.Join(strings.Split(structFields, ","), "\n\t"))
	_, _ = writer.WriteString(fmt.Sprintf("\n}\n"))
	_, _ = writer.WriteString(fmt.Sprintf("func (node %s) Accept(visitor Visitor) interface{} {\n\t", structName))
	_, _ = writer.WriteString(fmt.Sprintf("return visitor.Visit%s(node)\n}\n", structName))
	_, _ = writer.WriteString("\n")
}

// defineVisitor generates a Visitor interface for a given struct name and its fields and writes it to the provided writer.
func defineVisitor(writer *bufio.Writer, structName string, structFields []string) {
	_, _ = writer.WriteString(fmt.Sprintf("type Visitor interface{\n\t"))
	for _, type_ := range structFields {
		typeName := strings.Trim(strings.Split(type_, ":")[0], " ")
		_, _ = writer.WriteString(fmt.Sprintf("Visit%s(node %s) interface{}", strings.Trim(strings.Split(type_, ":")[0], " "), typeName))
		_, _ = writer.WriteString(fmt.Sprintf("\n\t"))
	}
	_, _ = writer.WriteString(fmt.Sprintf("\n}\n"))
}
