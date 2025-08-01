//ast/modify.go

package ast

type ModifyFunc func(Node) Node

func Modify(node Node, modifier ModifyFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i, _ := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, elem := range node.Elements {
			node.Elements[i], _ = Modify(elem, modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			modifiedKey, _ := Modify(key, modifier).(Expression)
			modifiedValue, _ := Modify(value, modifier).(Expression)
			newPairs[modifiedKey] = modifiedValue
		}
		node.Pairs = newPairs
	}
	return modifier(node)
}
