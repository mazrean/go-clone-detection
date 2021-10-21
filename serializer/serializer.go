package serializer

import (
	"context"
	"errors"
	"go/ast"
	"go/token"
	"log"

	"github.com/mazrean/go-clone-detection/domain"
	"github.com/mazrean/go-clone-detection/domain/values"
)

type Serializer struct {
	nodes []*domain.Node
}

func (s *Serializer) Serialize(ctx context.Context, root ast.Node, nodeChan chan<- *domain.Node) error {
	s.nodes = []*domain.Node{}

	visitor := &visitor{
		ctx:      ctx,
		nodeChan: nodeChan,
		stack:    []*stackValue{},
	}

	ast.Walk(visitor, root)

	return nil
}

type stackValue struct {
	node         *domain.Node
	childCounter func(values.ChildCount)
}

type visitor struct {
	ctx      context.Context
	nodeChan chan<- *domain.Node
	stack    []*stackValue
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		sValue := v.stack[len(v.stack)-1]
		v.stack = v.stack[:len(v.stack)-1]

		if sValue.childCounter != nil {
			sValue.childCounter(sValue.node.GetChildCount() + 1)
		}

		return nil
	}

	select {
	case <-v.ctx.Done():
		return nil
	default:
		nodeType, err := getNodeType(node)
		if err != nil {
			log.Printf("Error getting node type: %v", err)
			return nil
		}

		var childCounter func(count values.ChildCount)
		if len(v.stack) == 0 {
			childCounter = nil
		} else {
			childCounter = v.stack[len(v.stack)-1].node.IncrementChildCount
		}

		sValue := stackValue{
			node: domain.NewNode(
				node,
				nodeType,
				values.NewPosition(int64(node.Pos()), int64(node.End())),
				0,
				getNodeToken(node),
			),
			childCounter: childCounter,
		}

		select {
		case <-v.ctx.Done():
			return nil
		case v.nodeChan <- sValue.node:
		}

		v.stack = append(v.stack, &sValue)

		return v
	}
}

func getNodeType(node ast.Node) (values.NodeType, error) {
	switch node.(type) {
	case *ast.ArrayType:
		return values.NodeTypeArrayType, nil
	case *ast.AssignStmt:
		return values.NodeTypeAssignStmt, nil
	case *ast.BadDecl, *ast.BadExpr, *ast.BadStmt:
		return values.NodeTypeBad, nil
	case *ast.BasicLit:
		return values.NodeTypeBasicLit, nil
	case *ast.BinaryExpr:
		return values.NodeTypeBinaryExpr, nil
	case *ast.BlockStmt:
		return values.NodeTypeBlockStmt, nil
	case *ast.BranchStmt:
		return values.NodeTypeBranchStmt, nil
	case *ast.CallExpr:
		return values.NodeTypeCallExpr, nil
	case *ast.CaseClause:
		return values.NodeTypeCaseClause, nil
	case *ast.ChanType:
		return values.NodeTypeChanType, nil
	case *ast.CommClause:
		return values.NodeTypeCommClause, nil
	case *ast.Comment:
		return values.NodeTypeComment, nil
	case *ast.CommentGroup:
		return values.NodeTypeCommentGroup, nil
	case *ast.CompositeLit:
		return values.NodeTypeCompositeLit, nil
	case *ast.DeclStmt:
		return values.NodeTypeDeclStmt, nil
	case *ast.DeferStmt:
		return values.NodeTypeDeferStmt, nil
	case *ast.Ellipsis:
		return values.NodeTypeEllipsis, nil
	case *ast.EmptyStmt:
		return values.NodeTypeEmptyStmt, nil
	case *ast.ExprStmt:
		return values.NodeTypeExprStmt, nil
	case *ast.Field:
		return values.NodeTypeField, nil
	case *ast.FieldList:
		return values.NodeTypeFieldList, nil
	case *ast.File:
		return values.NodeTypeFile, nil
	case *ast.ForStmt:
		return values.NodeTypeForStmt, nil
	case *ast.FuncDecl:
		return values.NodeTypeFuncDecl, nil
	case *ast.FuncLit:
		return values.NodeTypeFuncLit, nil
	case *ast.FuncType:
		return values.NodeTypeFuncType, nil
	case *ast.GenDecl:
		return values.NodeTypeGenDecl, nil
	case *ast.GoStmt:
		return values.NodeTypeGoStmt, nil
	case *ast.Ident:
		return values.NodeTypeIdent, nil
	case *ast.IfStmt:
		return values.NodeTypeIfStmt, nil
	case *ast.ImportSpec:
		return values.NodeTypeImportSpec, nil
	case *ast.IncDecStmt:
		return values.NodeTypeIncDecStmt, nil
	case *ast.IndexExpr:
		return values.NodeTypeIndexExpr, nil
	case *ast.InterfaceType:
		return values.NodeTypeInterfaceType, nil
	case *ast.KeyValueExpr:
		return values.NodeTypeKeyValueExpr, nil
	case *ast.LabeledStmt:
		return values.NodeTypeLabeledStmt, nil
	case *ast.MapType:
		return values.NodeTypeMapType, nil
	case *ast.ParenExpr:
		return values.NodeTypeParenExpr, nil
	case *ast.RangeStmt:
		return values.NodeTypeRangeStmt, nil
	case *ast.ReturnStmt:
		return values.NodeTypeReturnStmt, nil
	case *ast.SelectStmt:
		return values.NodeTypeSelectStmt, nil
	case *ast.SelectorExpr:
		return values.NodeTypeSelectorExpr, nil
	case *ast.SendStmt:
		return values.NodeTypeSendStmt, nil
	case *ast.SliceExpr:
		return values.NodeTypeSliceExpr, nil
	case *ast.StarExpr:
		return values.NodeTypeStarExpr, nil
	case *ast.StructType:
		return values.NodeTypeStructType, nil
	case *ast.SwitchStmt:
		return values.NodeTypeSwitchStmt, nil
	case *ast.TypeAssertExpr:
		return values.NodeTypeTypeAssertExpr, nil
	case *ast.TypeSpec:
		return values.NodeTypeTypeSpec, nil
	case *ast.TypeSwitchStmt:
		return values.NodeTypeTypeSwitchStmt, nil
	case *ast.UnaryExpr:
		return values.NodeTypeUnaryExpr, nil
	case *ast.ValueSpec:
		return values.NodeTypeValueSpec, nil
	}

	return 0, errors.New("unknown node type")
}

func getNodeToken(node ast.Node) values.NodeToken {
	switch node := node.(type) {
	case *ast.AssignStmt:
		switch node.Tok {
		// AssignとDefineはコピペで置換されることが多いので区別しない
		case token.ASSIGN, token.DEFINE:
			return values.NodeTokenAssign
		case token.ADD_ASSIGN:
			return values.NodeTokenAddAssign
		case token.SUB_ASSIGN:
			return values.NodeTokenSubAssign
		case token.MUL_ASSIGN:
			return values.NodeTokenMultipleAssign
		case token.QUO_ASSIGN:
			return values.NodeTokenQuotientAssign
		case token.REM_ASSIGN:
			return values.NodeTokenRemainderAssign
		case token.AND_ASSIGN:
			return values.NodeTokenAndAssign
		case token.OR_ASSIGN:
			return values.NodeTokenOrAssign
		case token.XOR_ASSIGN:
			return values.NodeTokenXorAssign
		case token.SHL_ASSIGN:
			return values.NodeTokenShiftLeftAssign
		case token.SHR_ASSIGN:
			return values.NodeTokenShiftRightAssign
		case token.AND_NOT_ASSIGN:
			return values.NodeTokenAndNotAssign
		}

		log.Printf("unknown assign token: %v", node.Tok)
		return values.NodeTokenIllegal
	case *ast.BasicLit:
		switch node.Kind {
		case token.INT:
			return values.NodeTokenInt
		case token.FLOAT:
			return values.NodeTokenFloat
		case token.IMAG:
			return values.NodeTokenImaginary
		case token.CHAR:
			return values.NodeTokenChar
		case token.STRING:
			return values.NodeTokenString
		}

		log.Printf("unknown basic literal kind: %v", node.Kind)
		return values.NodeTokenIllegal
	case *ast.BinaryExpr:
		switch node.Op {
		case token.ADD:
			return values.NodeTokenAdd
		case token.SUB:
			return values.NodeTokenSub
		case token.MUL:
			return values.NodeTokenMultiple
		case token.QUO:
			return values.NodeTokenQuotient
		case token.REM:
			return values.NodeTokenRemainder
		case token.AND:
			return values.NodeTokenAnd
		case token.OR:
			return values.NodeTokenOr
		case token.XOR:
			return values.NodeTokenXor
		case token.SHL:
			return values.NodeTokenShiftLeft
		case token.SHR:
			return values.NodeTokenShiftRight
		case token.AND_NOT:
			return values.NodeTokenAndNot
		case token.LAND:
			return values.NodeTokenLogicalAnd
		case token.LOR:
			return values.NodeTokenLogicalOr
		case token.EQL:
			return values.NodeTokenEqual
		case token.LSS:
			return values.NodeTokenLess
		case token.GTR:
			return values.NodeTokenGreater
		case token.NEQ:
			return values.NodeTokenNotEqual
		case token.LEQ:
			return values.NodeTokenLessOrEqual
		case token.GEQ:
			return values.NodeTokenGreaterOrEqual

		}

		log.Printf("unknown binary token: %v", node.Op)
		return values.NodeTokenIllegal
	case *ast.BranchStmt:
		switch node.Tok {
		case token.BREAK:
			return values.NodeTokenBreak
		case token.CONTINUE:
			return values.NodeTokenContinue
		case token.GOTO:
			return values.NodeTokenGoto
		case token.FALLTHROUGH:
			return values.NodeTokenFallthrough
		}

		log.Printf("unknown branch token: %v", node.Tok)
		return values.NodeTokenIllegal
	case *ast.GenDecl:
		switch node.Tok {
		case token.IMPORT:
			return values.NodeTokenImport
		case token.CONST:
			return values.NodeTokenConst
		case token.TYPE:
			return values.NodeTokenType
		case token.VAR:
			return values.NodeTokenVar
		}

		log.Printf("unknown gen decl token: %v", node.Tok)
		return values.NodeTokenIllegal
	case *ast.IncDecStmt:
		switch node.Tok {
		case token.INC:
			return values.NodeTokenIncrement
		case token.DEC:
			return values.NodeTokenDecrement
		}

		log.Printf("unknown inc dec token: %v", node.Tok)
		return values.NodeTokenIllegal
	case *ast.UnaryExpr:
		switch node.Op {
		case token.AND:
			return values.NodeTokenAnd
		case token.NOT:
			return values.NodeTokenNot
		case token.ARROW:
			return values.NodeTokenArrow
		case token.SUB:
			return values.NodeTokenSub
		}

		log.Printf("unknown unary token: %v", node.Op)
		return values.NodeTokenIllegal
	}

	return values.NodeTokenNone
}
