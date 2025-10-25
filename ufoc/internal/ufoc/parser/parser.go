package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/ufocontract/internal/ufoc/lexer"
)

var Parser = participle.MustBuild[File](
	participle.Lexer(lexer.Def),
	participle.Elide("Whitespace", "Newline", "BlankLine"),
	participle.UseLookahead(2),
)
