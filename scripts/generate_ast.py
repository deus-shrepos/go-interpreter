"""
generate_ast.py

A utility script for generating Abstract Syntax Tree (AST) definitions in Go for a custom interpreter.
It automates the creation of Go structs, interfaces, and visitor patterns based on a predefined grammar.

Usage:
    python generate_ast.py --path ./output
"""

import os
import subprocess
import argparse
from dataclasses import dataclass
from typing import List
from io import TextIOWrapper

@dataclass
class ASTStructField:
    field_name: str
    field_type: str


@dataclass
class ASTStruct:
    name: str
    fields: list[ASTStructField]



def load_ast_grammar(file_path: os.PathLike) -> list[ASTStruct]:
    """
    Loads and parses an Abstract Syntax Tree (AST) grammar file.

    This function reads a file containing AST grammar rules line by line,
    processes each line using the `parse_ast_lines` function, and returns
    a list of parsed AST structures.

    Args:
        file_path (os.PathLike): The relative path to the AST grammar file.
            The path is resolved to an absolute path using the current working directory.

    Returns:
        list[ASTStruct]: A list of parsed AST structures derived from the grammar rules.

    Raises:
        FileNotFoundError: If the specified file does not exist.
        IOError: If there is an issue reading the file.
    """
    file_path = os.path.join(os.getcwd(), file_path)
    rules: list[ASTStruct] = []
    with open(file_path, "r") as file:
        line = file.readline().strip()
        while line:
            parsed_line = parse_ast_lines(line)
            if parsed_line:
                rules.append(parsed_line)
            line = file.readline().strip()
    return rules

def parse_ast_lines(line: str) -> ASTStruct:
    """
    Parses a single line of AST (Abstract Syntax Tree) definition and converts it into an `ASTStruct` object.

    The input line is expected to define a structure in the format:
    `<StructName> -> <FieldName1>:<FieldType1> | <FieldName2>:<FieldType2> | ...`
    - Lines starting with `#` are treated as comments and ignored.
    - The left-hand side of the `->` specifies the structure name.
    - The right-hand side specifies the fields of the structure, separated by `|`.
      Each field is defined as `<FieldName>:<FieldType>`.

    Args:
        line (str): A single line of AST definition.

    Returns:
        ASTStruct: An object representing the parsed structure and its fields.
                   If the line is a comment (starts with `#`), `None` is returned.

    Raises:
        AttributeError: If the field definition does not match the expected pattern.
        ValueError: If the line is malformed or missing required components.

    Example:
        Input:
            "Node -> name:str | value:int"
        Output:
            ASTStruct(
                name="Node",
                fields=[
                    ASTStructField(name="name", type="str"),
                    ASTStructField(name="value", type="int")
                ]
            )
    """

    if line.startswith("#"): # ignore comments
        return
    rule = line.strip().split("->")
    struct = rule[0].strip()
    fields = rule[1].split("|")
    struct_fields: list[ASTStructField] = []
    for field in fields:   
        field_name, field_type = field.strip().split(":")
        struct_fields.append(ASTStructField(field_name, field_type))
    return ASTStruct(name=struct, fields=struct_fields)

def generate_ast(ast_path: os.PathLike | str, output_path: os.PathLike | str) -> None:
    """Generate an Abstract Syntax Tree (AST) Go file in the specified location.

    This function takes a single argument specifying the output directory
    where the generated AST file will be created. If the number of arguments
    provided is not exactly one, it will print a usage message.

    *args (ParamSpec.args): A variadic argument list, where the first
            argument should be the output directory for the generated AST file.

    Raises:
        ValueError: If the number of arguments provided is not exactly one.

    """
    if not output_path or not ast_path:
        print("Usage: python -m generate_ast <ast grammar path> <output directoy>")
        raise ValueError()
    ast_grammar = load_ast_grammar(ast_path)
    print(ast_grammar)
    define_ast(output_path, "Expr", ast_grammar)
        

def define_ast(output_dir: os.PathLike, base_name: str, ast_structs: List[ASTStruct]) -> None:
    """Generates the Abstract Syntax Tree (AST) definitions in Go based on the provided grammar rules.

    This function creates a Go source file containing the AST definitions, including the base interface,
    visitor interface, and struct types for each production rule. The generated file is written to the
    specified output directory.

    Args:
        output_dir (os.PathLike): The directory where the generated Go file will be saved.
        base_name (str): The name of the base interface representing the AST node.
        ast_structs (List[ASTStruct]): A list of AST struct rules defining the structure of the AST in golang.
    
    Returns:
        None

    """
    path: str = os.path.join(f"{output_dir}/ast.go")
    os.makedirs(os.path.dirname(path), exist_ok=True)
    
    with open(path, "w") as file:
        file.write("package ast \n\n")
        file.write("import \"Crafting-interpreters/internal/token\"\n\n")
        define_visitor(file, ast_structs)
        file.write(f"type {base_name} interface {{ \n")
        file.write("\t Accept(vistior Visitor) (any, error) \n")
        file.write("}\n\n")
        for struct in ast_structs:
            define_type(file, struct.name, struct.fields)

    print(f"Generate AST file in {path}")


def define_type(file: TextIOWrapper, struct_name: str, struct_field: list[ASTStructField]) -> None:
    """Writes the definition of a Go struct and its associated Accept method to a file.

    Args:
        file (TextIOWrapper): The file object to write the struct definition and method to.
        struct_name (str): The name of the Go struct to define.
        struct_field (str): The fields of the Go struct, provided as a string.

    Returns:
        None
        
    """
    file.write(f"type {struct_name} struct {{ \n")
    file.write("\n\t".join(f"{s.field_name} {s.field_type}" for s in struct_field))
    file.write("\n}\n")
    file.write(f"func (node {struct_name}) Accept(visitor Visitor) (any, error) {'{'} \n\t")
    file.write(f"return visitor.Visit{struct_name}(node)\n")
    file.write("}\n\n")
    file.write("\n")
    

def define_visitor(file: TextIOWrapper, structs: List[ASTStruct]) -> None:
    """Generates and writes the definition of a Visitor interface to the provided file.

    The Visitor interface includes methods for visiting each type specified in the
    `struct_fields` list. Each method is named `Visit<TypeName>` and accepts a parameter
    of the corresponding type, returning a node and an error.

    Args:

        file (TextIOWrapper): The file object to which the Visitor interface definition
                              will be written.
        structs (List[ASTStruct]): A list of Golang AST Structs loaded from the AST textfile.

    Returns:
        None

    """
    file.write("type Visitor interface {\n")
    for struct in structs:
        file.write(f"\tVisit{struct.name}(node {struct.name}) (any, error)\n")
    file.write("}\n\n")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="AST Generator")
    parser.add_argument("-a", "--ast", help="ast grammar file", type=str, required=True)
    parser.add_argument('-p', '--path', help="file path to store the ast.go", type=str, required=True)
    args = parser.parse_args()
    generate_ast(ast_path=args.ast, output_path=args.path)
    subprocess.run(["gofmt", "-w", args.path])