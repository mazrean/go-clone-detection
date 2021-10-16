package values

type (
	NodeType byte
	Position struct {
		file  string
		start int64
		end   int64
	}
	ChildCount int64
)

const (
	NodeTypeBad NodeType = iota
	NodeTypeArrayType
	NodeTypeAssignStmt
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

func NewPosition(file string, start, end int64) *Position {
	return &Position{
		file:  file,
		start: start,
		end:   end,
	}
}

func (p *Position) GetFile() string {
	return p.file
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
