package lexer

import (
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/require"
)

func TestLexerKeywords(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"version", "version", []lexer.Token{
			{Type: symbols["Keyword"], Value: "version"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"namespace", "namespace", []lexer.Token{
			{Type: symbols["Keyword"], Value: "namespace"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"type", "type", []lexer.Token{
			{Type: symbols["Keyword"], Value: "type"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"enum", "enum", []lexer.Token{
			{Type: symbols["Keyword"], Value: "enum"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"const", "const", []lexer.Token{
			{Type: symbols["Keyword"], Value: "const"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"pattern", "pattern", []lexer.Token{
			{Type: symbols["Keyword"], Value: "pattern"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"deprecated", "deprecated", []lexer.Token{
			{Type: symbols["Keyword"], Value: "deprecated"},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"simple", "myVar", []lexer.Token{
			{Type: symbols["Ident"], Value: "myVar"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"camelCase", "myVariable", []lexer.Token{
			{Type: symbols["Ident"], Value: "myVariable"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"PascalCase", "MyVariable", []lexer.Token{
			{Type: symbols["Ident"], Value: "MyVariable"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"with_numbers", "var123", []lexer.Token{
			{Type: symbols["Ident"], Value: "var123"},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerNumbers(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"integer", "42", []lexer.Token{
			{Type: symbols["Number"], Value: "42"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"negative", "-10", []lexer.Token{
			{Type: symbols["Number"], Value: "-10"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"positive", "+5", []lexer.Token{
			{Type: symbols["Number"], Value: "+5"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"float", "3.14", []lexer.Token{
			{Type: symbols["Number"], Value: "3.14"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"float_no_leading", ".5", []lexer.Token{
			{Type: symbols["Number"], Value: ".5"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"negative_float", "-2.718", []lexer.Token{
			{Type: symbols["Number"], Value: "-2.718"},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerStrings(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"simple", `"hello"`, []lexer.Token{
			{Type: symbols["String"], Value: `"hello"`},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"empty", `""`, []lexer.Token{
			{Type: symbols["String"], Value: `""`},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"with_spaces", `"hello world"`, []lexer.Token{
			{Type: symbols["String"], Value: `"hello world"`},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"escaped_quote", `"say \"hello\""`, []lexer.Token{
			{Type: symbols["String"], Value: `"say \"hello\""`},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"escaped_backslash", `"path\\to\\file"`, []lexer.Token{
			{Type: symbols["String"], Value: `"path\\to\\file"`},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"escaped_newline", `"line1\nline2"`, []lexer.Token{
			{Type: symbols["String"], Value: `"line1\nline2"`},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerDocstrings(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"simple", "\"\"\"Simple docstring\"\"\"", []lexer.Token{
			{Type: symbols["Docstring"], Value: "\"\"\"Simple docstring\"\"\""},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"empty", "\"\"\"\"\"\"", []lexer.Token{
			{Type: symbols["Docstring"], Value: "\"\"\"\"\"\""},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"multiline", "\"\"\"Line 1\n\t\t\tLine 2\n\t\t\tLine 3\"\"\"", []lexer.Token{
			{Type: symbols["Docstring"], Value: "\"\"\"Line 1\n\t\t\tLine 2\n\t\t\tLine 3\"\"\""},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"with_markdown", "\"\"\"# Title\n\t\t\t## Subtitle\n\t\t\tSome text\"\"\"", []lexer.Token{
			{Type: symbols["Docstring"], Value: "\"\"\"# Title\n\t\t\t## Subtitle\n\t\t\tSome text\"\"\""},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"with_single_quote", "\"\"\"It's a docstring\"\"\"", []lexer.Token{
			{Type: symbols["Docstring"], Value: "\"\"\"It's a docstring\"\"\""},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerComments(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"single_line", "// This is a comment", []lexer.Token{
			{Type: symbols["Comment"], Value: "// This is a comment"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"single_line_empty", "//", []lexer.Token{
			{Type: symbols["Comment"], Value: "//"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"multi_line", "/* Multi line comment */", []lexer.Token{
			{Type: symbols["BlockComment"], Value: "/* Multi line comment */"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"multi_line_multiline", "/* Line 1\n\t\tLine 2 */", []lexer.Token{
			{Type: symbols["BlockComment"], Value: "/* Line 1\n\t\tLine 2 */"},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerPunctuation(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{"open_brace", "{", []lexer.Token{
			{Type: symbols["Punct"], Value: "{"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"close_brace", "}", []lexer.Token{
			{Type: symbols["Punct"], Value: "}"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"open_paren", "(", []lexer.Token{
			{Type: symbols["Punct"], Value: "("},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"close_paren", ")", []lexer.Token{
			{Type: symbols["Punct"], Value: ")"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"open_bracket", "[", []lexer.Token{
			{Type: symbols["Punct"], Value: "["},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"close_bracket", "]", []lexer.Token{
			{Type: symbols["Punct"], Value: "]"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"colon", ":", []lexer.Token{
			{Type: symbols["Punct"], Value: ":"},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"equals", "=", []lexer.Token{
			{Type: symbols["Punct"], Value: "="},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"comma", ",", []lexer.Token{
			{Type: symbols["Punct"], Value: ","},
			{Type: symbols["EOF"], Value: ""},
		}},
		{"question", "?", []lexer.Token{
			{Type: symbols["Punct"], Value: "?"},
			{Type: symbols["EOF"], Value: ""},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexerComplexExamples(t *testing.T) {
	symbols := Def.Symbols()

	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{
			"version_statement",
			"version 1",
			[]lexer.Token{
				{Type: symbols["Keyword"], Value: "version"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Number"], Value: "1"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"namespace_declaration",
			"namespace \"Tasks\" {",
			[]lexer.Token{
				{Type: symbols["Keyword"], Value: "namespace"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["String"], Value: "\"Tasks\""},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Punct"], Value: "{"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"type_field",
			"fieldName: string",
			[]lexer.Token{
				{Type: symbols["Ident"], Value: "fieldName"},
				{Type: symbols["Punct"], Value: ":"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Ident"], Value: "string"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"optional_field",
			"fieldName?: string[]",
			[]lexer.Token{
				{Type: symbols["Ident"], Value: "fieldName"},
				{Type: symbols["Punct"], Value: "?"},
				{Type: symbols["Punct"], Value: ":"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Ident"], Value: "string"},
				{Type: symbols["Punct"], Value: "["},
				{Type: symbols["Punct"], Value: "]"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"const_declaration",
			"const MaxRetries: int = 5",
			[]lexer.Token{
				{Type: symbols["Keyword"], Value: "const"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Ident"], Value: "MaxRetries"},
				{Type: symbols["Punct"], Value: ":"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Ident"], Value: "int"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Punct"], Value: "="},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Number"], Value: "5"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"enum_member",
			"PENDING = 1",
			[]lexer.Token{
				{Type: symbols["Ident"], Value: "PENDING"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Punct"], Value: "="},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Number"], Value: "1"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"pattern_declaration",
			"pattern TaskTopic = \"{ns}.{taskId}.updates\"",
			[]lexer.Token{
				{Type: symbols["Keyword"], Value: "pattern"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Ident"], Value: "TaskTopic"},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["Punct"], Value: "="},
				{Type: symbols["Whitespace"], Value: " "},
				{Type: symbols["String"], Value: "\"{ns}.{taskId}.updates\""},
				{Type: symbols["EOF"], Value: ""},
			},
		},
		{
			"deprecated_with_message",
			"deprecated(\"Use new version\")",
			[]lexer.Token{
				{Type: symbols["Keyword"], Value: "deprecated"},
				{Type: symbols["Punct"], Value: "("},
				{Type: symbols["String"], Value: "\"Use new version\""},
				{Type: symbols["Punct"], Value: ")"},
				{Type: symbols["EOF"], Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
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

	symbols := Def.Symbols()
	expected := []lexer.Token{
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t"},
		{Type: symbols["Keyword"], Value: "version"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Number"], Value: "1"},
		{Type: symbols["BlankLine"], Value: "\n\n"},
		{Type: symbols["Whitespace"], Value: "\t\t"},
		{Type: symbols["Comment"], Value: "// Comment"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t"},
		{Type: symbols["Docstring"], Value: "\"\"\"\n\t\tDocumentation for namespace\n\t\t\"\"\""},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t"},
		{Type: symbols["Keyword"], Value: "namespace"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["String"], Value: "\"Tasks\""},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Punct"], Value: "{"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Docstring"], Value: "\"\"\"\n\t\t\tType documentation\n\t\t\t\"\"\""},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Keyword"], Value: "type"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "Task"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Punct"], Value: "{"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t\t"},
		{Type: symbols["Docstring"], Value: "\"\"\" Field doc \"\"\""},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t\t"},
		{Type: symbols["Ident"], Value: "id"},
		{Type: symbols["Punct"], Value: ":"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "string"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t\t"},
		{Type: symbols["Ident"], Value: "status"},
		{Type: symbols["Punct"], Value: "?"},
		{Type: symbols["Punct"], Value: ":"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "int"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Punct"], Value: "}"},
		{Type: symbols["BlankLine"], Value: "\n\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Keyword"], Value: "const"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "MaxRetries"},
		{Type: symbols["Punct"], Value: ":"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "int"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Punct"], Value: "="},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Number"], Value: "5"},
		{Type: symbols["BlankLine"], Value: "\n\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Keyword"], Value: "enum"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "TaskStatus"},
		{Type: symbols["Punct"], Value: ":"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "string"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Punct"], Value: "{"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t\t"},
		{Type: symbols["Ident"], Value: "PENDING"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t\t"},
		{Type: symbols["Ident"], Value: "RUNNING"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Punct"], Value: "}"},
		{Type: symbols["BlankLine"], Value: "\n\n"},
		{Type: symbols["Whitespace"], Value: "\t\t\t"},
		{Type: symbols["Keyword"], Value: "pattern"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Ident"], Value: "Topic"},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["Punct"], Value: "="},
		{Type: symbols["Whitespace"], Value: " "},
		{Type: symbols["String"], Value: "\"tasks.{taskId}\""},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t\t"},
		{Type: symbols["Punct"], Value: "}"},
		{Type: symbols["Newline"], Value: "\n"},
		{Type: symbols["Whitespace"], Value: "\t"},
		{Type: symbols["EOF"], Value: ""},
	}

	assertTokens(t, input, expected)
}

/*******************
* HELPER FUNCTIONS *
*******************/

var tokenSymbols = func() map[lexer.TokenType]string {
	symbols := Def.Symbols()
	inverse := make(map[lexer.TokenType]string)
	for name, typ := range symbols {
		inverse[typ] = name
	}
	return inverse
}()

func assertTokens(t *testing.T, input string, expected []lexer.Token) {
	t.Helper()

	lex, err := Def.LexString("", input)
	require.NoError(t, err)
	tokens, err := lexer.ConsumeAll(lex)
	require.NoError(t, err)

	require.Equal(t, len(expected), len(tokens), "Token count mismatch")

	for i, exp := range expected {
		require.Equal(t, exp.Type, tokens[i].Type, "Token type mismatch at index %d: expected %s, got %s", i, tokenSymbols[exp.Type], tokenSymbols[tokens[i].Type])
		require.Equal(t, exp.Value, tokens[i].Value, "Token value mismatch at index %d", i)
	}
}
