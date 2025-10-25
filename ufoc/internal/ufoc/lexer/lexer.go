package lexer

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var Def = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//[^\n]*`},
	{Name: "BlockComment", Pattern: `/\*[^*]*\*+(?:[^/*][^*]*\*+)*/`},
	{Name: "Docstring", Pattern: `"""[^"]*(?:"[^"][^"]*|""[^"][^"]*)*"""`},
	{Name: "Keyword", Pattern: `\b(?:version|namespace|type|enum|const|pattern|deprecated)\b`},
	{Name: "Number", Pattern: `[-+]?(?:\d*\.)?\d+`},
	{Name: "String", Pattern: `"(?:[^"\\]|\\["\\/bfnrt]|\\u[0-9a-fA-F]{4})*"`},
	{Name: "Ident", Pattern: `[a-zA-Z][a-zA-Z0-9]*`},
	{Name: "Punct", Pattern: `[{}()\[\]:=,?]`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})
