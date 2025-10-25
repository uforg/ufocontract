package parser

import "github.com/alecthomas/participle/v2/lexer"

type Comment struct {
	Pos  lexer.Position `parser:""`
	Text string         `parser:"@Comment"`
}

type BlockComment struct {
	Pos  lexer.Position `parser:""`
	Text string         `parser:"@BlockComment"`
}

type File struct {
	Pos      lexer.Position `parser:""`
	Version  int            `parser:"'version' @Number"`
	Children []*FileChild   `parser:"@@*"`
}

type FileChild struct {
	Pos          lexer.Position `parser:""`
	Docstring    *Docstring     `parser:"@@"`
	Comment      *Comment       `parser:"| @@"`
	BlockComment *BlockComment  `parser:"| @@"`
	Namespace    *Namespace     `parser:"| @@"`
}

type Docstring struct {
	Pos       lexer.Position `parser:""`
	Text      string         `parser:"@Docstring"`
	BlankLine bool           `parser:"@BlankLine"`
}

type Namespace struct {
	Pos           lexer.Position    `parser:""`
	Comments      []*Comment        `parser:"@@*"`
	BlockComments []*BlockComment   `parser:"@@*"`
	Docstring     *string           `parser:"@Docstring?"`
	Name          string            `parser:"'namespace' @Ident '{'"`
	Children      []*NamespaceChild `parser:"@@* '}'"`
}

type NamespaceChild struct {
	Pos lexer.Position `parser:""`

	Docstring    *Docstring    `parser:"@@"`
	Comment      *Comment      `parser:"| @@"`
	BlockComment *BlockComment `parser:"| @@"`
	Type         *TypeDef      `parser:"| @@"`
	Enum         *EnumDef      `parser:"| @@"`
	Const        *ConstDef     `parser:"| @@"`
	Pattern      *PatternDef   `parser:"| @@"`
}

type TypeDef struct {
	Pos        lexer.Position `parser:""`
	Docstring  *string        `parser:"@Docstring?"`
	Deprecated *string        `parser:"( 'deprecated' ( '(' @String ')' )? )?"`
	Name       string         `parser:"'type' @Ident '{'"`
	Fields     []*Field       `parser:"@@* '}'"`
}

type Field struct {
	Pos           lexer.Position  `parser:""`
	Comments      []*Comment      `parser:"@@*"`
	BlockComments []*BlockComment `parser:"@@*"`
	Docstring     *string         `parser:"@Docstring?"`
	Name          string          `parser:"@Ident"`
	Optional      bool            `parser:"@'?'?"`
	Type          *TypeRef        `parser:"':' @@"`
}

type TypeRef struct {
	Pos lexer.Position `parser:""`

	Inline *InlineType `parser:"  @@"`
	Named  *string     `parser:"| @Ident"`
	Array  bool        `parser:"  @( '[' ']' )?"`
}

type InlineType struct {
	Pos    lexer.Position `parser:""`
	Fields []*Field       `parser:"'{' @@* '}'"`
}

type EnumDef struct {
	Pos        lexer.Position `parser:""`
	Docstring  *string        `parser:"@Docstring?"`
	Deprecated *string        `parser:"( 'deprecated' ( '(' @String ')' )? )?"`
	Name       string         `parser:"'enum' @Ident"`
	BaseType   *string        `parser:"( ':' @Ident )?"`
	Members    []*EnumMember  `parser:"'{' @@* '}'"`
}

type EnumMember struct {
	Pos           lexer.Position  `parser:""`
	Comments      []*Comment      `parser:"@@*"`
	BlockComments []*BlockComment `parser:"@@*"`
	Docstring     *string         `parser:"@Docstring?"`
	Name          string          `parser:"@Ident"`
	Value         *Value          `parser:"( '=' @@ )?"`
}

type ConstDef struct {
	Pos        lexer.Position `parser:""`
	Docstring  *string        `parser:"@Docstring?"`
	Deprecated *string        `parser:"( 'deprecated' ( '(' @String ')' )? )?"`
	Name       string         `parser:"'const' @Ident"`
	Type       *TypeRef       `parser:"':' @@"`
	Value      *Value         `parser:"'=' @@"`
}

type PatternDef struct {
	Pos        lexer.Position `parser:""`
	Docstring  *string        `parser:"@Docstring?"`
	Deprecated *string        `parser:"( 'deprecated' ( '(' @String ')' )? )?"`
	Name       string         `parser:"'pattern' @Ident"`
	Pattern    string         `parser:"'=' @String"`
}

type Value struct {
	Pos lexer.Position `parser:""`

	String *string `parser:"  @String"`
	Number *string `parser:"| @Number"`
	Ident  *string `parser:"| @Ident"`
}
