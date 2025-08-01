//evaluator/macro_expansion.go

package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func DefineMacros(program *ast.Program, environment *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, environment)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex+1:]...)
	}
}

func addMacro(stmt ast.Statement, environment *object.Environment) {
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:       environment,
		Body:      macroLiteral.Body,
	}
	environment.Set(letStatement.Name.Value, macro)
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStatement, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}


func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
    return ast.Modify(program, func(node ast.Node) ast.Node {
        callExpression, ok := node.(*ast.CallExpression)
        if !ok {
            return node
        }

        macro, ok := isMacroCall(callExpression, env)
        if !ok {
            return node
        }

        args := quoteArgs(callExpression)
        evalEnv := extendMacroEnv(macro, args)

        evaluated := Eval(macro.Body, evalEnv)

        quote, ok := evaluated.(*object.Quote)
        if !ok {
            panic("we only support returning AST-nodes from macros")
        }

        return quote.Node
    })
}

func isMacroCall(
    exp *ast.CallExpression,
    env *object.Environment,
) (*object.Macro, bool) {
    identifier, ok := exp.Function.(*ast.Identifier)
    if !ok {
        return nil, false
    }

    obj, ok := env.Get(identifier.Value)
    if !ok {
        return nil, false
    }

    macro, ok := obj.(*object.Macro)
    if !ok {
        return nil, false
    }

    return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
    args := []*object.Quote{}

    for _, a := range exp.Arguments {
        args = append(args, &object.Quote{Node: a})
    }

    return args
}

func extendMacroEnv(
    macro *object.Macro,
    args []*object.Quote,
) *object.Environment {
    extended := object.NewEnclosedEnvironment(macro.Env)

    for paramIdx, param := range macro.Parameters {
        extended.Set(param.Value, args[paramIdx])
    }

    return extended
}