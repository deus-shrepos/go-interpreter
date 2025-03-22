"""
This script will be used to generate the AST stuff. Doing it in python is way easier. 

"""

import os
import argparse
from typing import List
from io import TextIOWrapper

# This is the interpreter's Grammar for GoLang
Grammer: List[str] = [
	"Binary: Left Expr,Operator token.Token,Right Expr",
	"Grouping: Expression Expr",
	"Literal: Value any",
	"Unary: Operator token.Token,Right Expr",
	"stmt: Expression Expr",
	"Print: Expression Expr",
    # And more to come soon...
]

def generate_ast(path: os.PathLike | str) -> None:
    """Generate an Abstract Syntax Tree (AST) Go file in the specified location.

    This function takes a single argument specifying the output directory
    where the generated AST file will be created. If the number of arguments
    provided is not exactly one, it will print a usage message.

    *args (ParamSpec.args): A variadic argument list, where the first
            argument should be the output directory for the generated AST file.

    Raises:
        ValueError: If the number of arguments provided is not exactly one.

    """
    if not path:
        print("Usage: generate_ast <output directoy>")
        raise ValueError()
 
    define_ast(path, "Expr", Grammer)
        

def define_ast(output_dir: os.PathLike, base_name: str, production_rules: List[str]) -> None:
    """Generates the Abstract Syntax Tree (AST) definitions in Go based on the provided grammar rules.

    This function creates a Go source file containing the AST definitions, including the base interface,
    visitor interface, and struct types for each production rule. The generated file is written to the
    specified output directory.

    Args:
        output_dir (os.PathLike): The directory where the generated Go file will be saved.
        base_name (str): The name of the base interface representing the AST node.
        production_rules (List[str]): A list of grammar rules defining the structure of the AST. Each rule
            should be in the format "StructName: field1 type1, field2 type2, ...".
    
    Returns:
        None

    """
    path: str = os.path.join(os.getcwd(), f"{output_dir}/ast.go")
    if not os.path.exists(path):
        print("No path found!")
    
    with open(path, "w") as file:
        file.write("package ast \n\n")
        file.write("import \"Crafting-interpreters/internal/token\"\n\n")
        define_visitor(file, production_rules)
        file.write(f"type {base_name} interface {{ \n")
        file.write("\t Accept(vistior Visitor) (any, error) \n")
        file.write("}\n\n")
        for rule in production_rules:
            struct_name = rule.split(":")[0].lstrip()
            struct_fields = rule.split(":")[1].lstrip() 
            define_type(file, struct_name, struct_fields)
    
    print(f"Generate AST file in {path}")


def define_type(file: TextIOWrapper, struct_name: str, struct_field: str) -> None:
    """Writes the definition of a Go struct and its associated Accept method to a file.

    Args:
        file (TextIOWrapper): The file object to write the struct definition and method to.
        struct_name (str): The name of the Go struct to define.
        struct_field (str): The fields of the Go struct, provided as a string.

    Returns:
        None
    """
    file.write(f"type {struct_name} struct {{ \n")
    file.write("\n\t".join(s for s in struct_field.split(",")))
    file.write("\n}\n")
    file.write(f"func (node {struct_name}) Accept(visitor Visitor) (any, error) {'{'} \n\t")
    file.write(f"return visitor.Visit{struct_name}(node)\n")
    file.write("}\n\n")
    file.write("\n")
    

def define_visitor(file: TextIOWrapper, struct_fields: List[str]) -> None:
    """Generates and writes the definition of a Visitor interface to the provided file.

    The Visitor interface includes methods for visiting each type specified in the
    `struct_fields` list. Each method is named `Visit<TypeName>` and accepts a parameter
    of the corresponding type, returning a node and an error.

    Args:
        file (TextIOWrapper): The file object to which the Visitor interface definition
                              will be written.
        struct_fields (List[str]): A list of strings representing the types and their
                                   fields. Each string is expected to be in the format
                                   "TypeName: FieldType".

    Returns:
        None

    """
    file.write("type Visitor interface {\n")
    for type_ in struct_fields:
        type_name = type_.split(":")[0].lstrip()
        file.write(f"\tVisit{type_name}(node {type_name}) (any, error)\n")
    file.write("}\n\n")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="AST Generator")
    parser.add_argument('-p', '--path', help="file path to store the ast.go", type=str)
    args = parser.parse_args()
    generate_ast(path=args.path)
    os.system(f"gofmt -w {args.path}")