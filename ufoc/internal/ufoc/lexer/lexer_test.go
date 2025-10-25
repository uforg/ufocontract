package lexer

import (
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTokenSymbols() map[lexer.TokenType]string {
	symbols := Def.Symbols()
	inverse := make(map[lexer.TokenType]string)
	for name, typ := range symbols {
		inverse[typ] = name
	}
	return inverse
}

var tokenSymbols = getTokenSymbols()

func tokenTypeString(tok lexer.Token) string {
	return tokenSymbols[tok.Type]
}

func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"version", "version"},
		{"namespace", "namespace"},
		{"type", "type"},
		{"enum", "enum"},
		{"const", "const"},
		{"pattern", "pattern"},
		{"deprecated", "deprecated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "Keyword", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "myVar"},
		{"camelCase", "myVariable"},
		{"PascalCase", "MyVariable"},
		{"with_numbers", "var123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "Ident", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"integer", "42"},
		{"negative", "-10"},
		{"positive", "+5"},
		{"float", "3.14"},
		{"float_no_leading", ".5"},
		{"negative_float", "-2.718"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "Number", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", `"hello"`},
		{"empty", `""`},
		{"with_spaces", `"hello world"`},
		{"escaped_quote", `"say \"hello\""`},
		{"escaped_backslash", `"path\\to\\file"`},
		{"escaped_newline", `"line1\nline2"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "String", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerDocstrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", `"""Simple docstring"""`},
		{"empty", `""""""`},
		{"multiline", `"""Line 1
			Line 2
			Line 3"""`},
		{"with_markdown", `"""# Title
			## Subtitle
			Some text"""`},
		{"with_single_quote", `"""It's a docstring"""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "Docstring", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single_line", "// This is a comment", "Comment"},
		{"single_line_empty", "//", "Comment"},
		{"multi_line", "/* Multi line comment */", "BlockComment"},
		{"multi_line_multiline", `/* Line 1
		Line 2 */`, "BlockComment"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tokenTypeString(tokens[0]))
		})
	}
}

func TestLexerPunctuation(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"open_brace", "{"},
		{"close_brace", "}"},
		{"open_paren", "("},
		{"close_paren", ")"},
		{"open_bracket", "["},
		{"close_bracket", "]"},
		{"colon", ":"},
		{"equals", "="},
		{"comma", ","},
		{"question", "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)
			assert.Equal(t, "Punct", tokenTypeString(tokens[0]))
			assert.Equal(t, tt.input, tokens[0].Value)
		})
	}
}

func TestLexerComplexExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"version_statement",
			"version 1",
			[]string{"Keyword", "Whitespace", "Number"},
		},
		{
			"namespace_declaration",
			`namespace "Tasks" {`,
			[]string{"Keyword", "Whitespace", "String", "Whitespace", "Punct"},
		},
		{
			"type_field",
			"fieldName: string",
			[]string{"Ident", "Punct", "Whitespace", "Ident"},
		},
		{
			"optional_field",
			"fieldName?: string[]",
			[]string{"Ident", "Punct", "Punct", "Whitespace", "Ident", "Punct", "Punct"},
		},
		{
			"const_declaration",
			`const MaxRetries: int = 5`,
			[]string{"Keyword", "Whitespace", "Ident", "Punct", "Whitespace", "Ident", "Whitespace", "Punct", "Whitespace", "Number"},
		},
		{
			"enum_member",
			"PENDING = 1",
			[]string{"Ident", "Whitespace", "Punct", "Whitespace", "Number"},
		},
		{
			"pattern_declaration",
			`pattern TaskTopic = "{ns}.{taskId}.updates"`,
			[]string{"Keyword", "Whitespace", "Ident", "Whitespace", "Punct", "Whitespace", "String"},
		},
		{
			"deprecated_with_message",
			`deprecated("Use new version")`,
			[]string{"Keyword", "Punct", "String", "Punct"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex, err := Def.LexString("", tt.input)
			require.NoError(t, err)
			tokens, err := lexer.ConsumeAll(lex)
			require.NoError(t, err)

			var tokenTypes []string
			for _, tok := range tokens {
				typeName := tokenTypeString(tok)
				if typeName != "EOF" {
					tokenTypes = append(tokenTypes, typeName)
				}
			}

			assert.Equal(t, tt.expected, tokenTypes)
		})
	}
}

func TestLexerFullExample(t *testing.T) {
	input := `
		version 1

		// Comment
		"""
		Documentation for namespace
		"""
		namespace "Tasks" {
			"""
			Type documentation
			"""
			type Task {
				""" Field doc """
				id: string
				status?: int
			}

			const MaxRetries: int = 5

			enum TaskStatus: string {
				PENDING
				RUNNING
			}

			pattern Topic = "tasks.{taskId}"
		}
	`

	lex, err := Def.LexString("", input)
	require.NoError(t, err)
	tokens, err := lexer.ConsumeAll(lex)
	require.NoError(t, err)

	assert.Greater(t, len(tokens), 0)

	var hasKeyword, hasIdent, hasString, hasDocstring, hasComment bool
	for _, tok := range tokens {
		typeName := tokenTypeString(tok)
		switch typeName {
		case "Keyword":
			hasKeyword = true
		case "Ident":
			hasIdent = true
		case "String":
			hasString = true
		case "Docstring":
			hasDocstring = true
		case "Comment", "BlockComment":
			hasComment = true
		}
	}

	assert.True(t, hasKeyword, "Should have keywords")
	assert.True(t, hasIdent, "Should have identifiers")
	assert.True(t, hasString, "Should have strings")
	assert.True(t, hasDocstring, "Should have docstrings")
	assert.True(t, hasComment, "Should have comments")
}
