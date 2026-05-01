package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/enw860/covanalyze/internal/models"
	"github.com/golang/glog"
)

// exprToString converts an AST expression to a string representation
func exprToString(fset *token.FileSet, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return buf.String()
}

// parseSourceFile parses a Go source file and returns a list of relevant AST nodes.
// It handles parsing errors gracefully by logging warnings and returning an error.
func parseSourceFile(filepath string) ([]models.ASTNode, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		glog.Warningf("Source file not found: %s", filepath)
		return nil, err
	}

	// Create a new token file set
	fset := token.NewFileSet()

	// Parse the source file
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		glog.Warningf("Failed to parse source file %s: %v", filepath, err)
		return nil, err
	}

	var nodes []models.ASTNode

	// Traverse the AST to find function declarations
	ast.Inspect(file, func(n ast.Node) bool {
		// Check if this is a function declaration
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Get the function name
		funcName := funcDecl.Name.Name
		glog.V(2).Infof("Entering function: %s at line %d", funcName, fset.Position(funcDecl.Pos()).Line)

		// Now inspect the function body for relevant constructs
		if funcDecl.Body != nil {
			ast.Inspect(funcDecl.Body, func(bodyNode ast.Node) bool {
				if bodyNode == nil {
					return false
				}

				var node models.ASTNode
				var shouldAdd bool

				switch stmt := bodyNode.(type) {
				case *ast.IfStmt:
					// If statement
					node = models.ASTNode{
						NodeType:     models.NodeTypeIf,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
						Condition:    exprToString(fset, stmt.Cond),
					}
					shouldAdd = true

					// Handle else block if present
					if stmt.Else != nil {
						// Check if it's an else-if or a plain else
						if _, isElseIf := stmt.Else.(*ast.IfStmt); !isElseIf {
							// It's a plain else block
							elseNode := models.ASTNode{
								NodeType:     models.NodeTypeElse,
								StartLine:    fset.Position(stmt.Else.Pos()).Line,
								EndLine:      fset.Position(stmt.Else.End()).Line,
								FunctionName: funcName,
							}
							nodes = append(nodes, elseNode)
						}
						// else-if blocks will be caught naturally by ast.Inspect
					}

				case *ast.SwitchStmt:
					// Switch statement
					node = models.ASTNode{
						NodeType:     models.NodeTypeSwitch,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
						Condition:    exprToString(fset, stmt.Tag),
					}
					shouldAdd = true

				case *ast.TypeSwitchStmt:
					// Type switch statement
					node = models.ASTNode{
						NodeType:     models.NodeTypeTypeSwitch,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
					}
					shouldAdd = true

				case *ast.ForStmt:
					// For loop
					node = models.ASTNode{
						NodeType:     models.NodeTypeFor,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
						Condition:    exprToString(fset, stmt.Cond),
					}
					shouldAdd = true

				case *ast.RangeStmt:
					// Range loop
					node = models.ASTNode{
						NodeType:     models.NodeTypeRange,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
						Condition:    exprToString(fset, stmt.X),
					}
					shouldAdd = true

				case *ast.CaseClause:
					// Case clause in switch statement
					var condition string
					if len(stmt.List) > 0 {
						// Build condition from case expressions
						var buf bytes.Buffer
						for i, expr := range stmt.List {
							if i > 0 {
								buf.WriteString(", ")
							}
							buf.WriteString(exprToString(fset, expr))
						}
						condition = buf.String()
					}
					node = models.ASTNode{
						NodeType:     models.NodeTypeCase,
						StartLine:    fset.Position(stmt.Pos()).Line,
						EndLine:      fset.Position(stmt.End()).Line,
						FunctionName: funcName,
						Condition:    condition,
					}
					shouldAdd = true
				}

				if shouldAdd {
					nodes = append(nodes, node)
				}

				return true
			})
		}

		return true
	})

	return nodes, nil
}
