package values

type (
	NodeType  byte
	NodeToken byte
	Position  struct {
		start int64
		end   int64
	}
	ChildCount int64
)

const (
	NodeTypeBad NodeType = iota
	NodeTypeArrayType
	NodeTypeAssignStmt
	NodeTypeBasicLit
	NodeTypeBinaryExpr
	NodeTypeBlockStmt
	NodeTypeBranchStmt
	NodeTypeCallExpr
	NodeTypeCaseClause
	NodeTypeChanType
	NodeTypeCommClause
	NodeTypeComment
	NodeTypeCommentGroup
	NodeTypeCompositeLit
	NodeTypeDeclStmt
	NodeTypeDeferStmt
	NodeTypeEllipsis
	NodeTypeEmptyStmt
	NodeTypeExprStmt
	NodeTypeField
	NodeTypeFieldList
	NodeTypeFile
	NodeTypeForStmt
	NodeTypeFuncDecl
	NodeTypeFuncLit
	NodeTypeFuncType
	NodeTypeGenDecl
	NodeTypeGoStmt
	NodeTypeIdent
	NodeTypeIfStmt
	NodeTypeImportSpec
	NodeTypeIncDecStmt
	NodeTypeIndexExpr
	NodeTypeInterfaceType
	NodeTypeKeyValueExpr
	NodeTypeLabeledStmt
	NodeTypeMapType
	NodeTypeMergeMode
	NodeTypePackage
	NodeTypeParenExpr
	NodeTypeRangeStmt
	NodeTypeReturnStmt
	NodeTypeSelectStmt
	NodeTypeSelectorExpr
	NodeTypeSendStmt
	NodeTypeSliceExpr
	NodeTypeStarExpr
	NodeTypeStructType
	NodeTypeSwitchStmt
	NodeTypeTypeAssertExpr
	NodeTypeTypeSpec
	NodeTypeTypeSwitchStmt
	NodeTypeUnaryExpr
	NodeTypeValueSpec
)

const (
	NodeTokenNone NodeToken = iota
	NodeTokenIllegal

	NodeTokenInt
	NodeTokenFloat
	NodeTokenImaginary
	NodeTokenChar
	NodeTokenString

	NodeTokenAdd
	NodeTokenSub
	NodeTokenMultiple
	NodeTokenQuotient
	NodeTokenRemainder

	NodeTokenLogicalAnd
	NodeTokenLogicalOr
	NodeTokenEqual
	NodeTokenNotEqual
	NodeTokenLess
	NodeTokenGreater
	NodeTokenLessOrEqual
	NodeTokenGreaterOrEqual
	NodeTokenAnd
	NodeTokenOr
	NodeTokenXor
	NodeTokenShiftLeft
	NodeTokenShiftRight
	NodeTokenAndNot

	NodeTokenAssign
	NodeTokenAddAssign
	NodeTokenSubAssign
	NodeTokenMultipleAssign
	NodeTokenQuotientAssign
	NodeTokenRemainderAssign
	NodeTokenAndAssign
	NodeTokenOrAssign
	NodeTokenXorAssign
	NodeTokenShiftLeftAssign
	NodeTokenShiftRightAssign
	NodeTokenAndNotAssign

	NodeTokenIncrement
	NodeTokenDecrement

	NodeTokenNot
	NodeTokenArrow

	NodeTokenBreak
	NodeTokenContinue
	NodeTokenFallthrough
	NodeTokenGoto

	NodeTokenConst

	NodeTokenImport
	NodeTokenType
	NodeTokenVar
)

func NewPosition(start, end int64) *Position {
	return &Position{
		start: start,
		end:   end,
	}
}

func (p *Position) GetStart() int64 {
	return p.start
}

func (p *Position) GetEnd() int64 {
	return p.end
}

func NewChildCount(childCount int64) ChildCount {
	return ChildCount(childCount)
}
