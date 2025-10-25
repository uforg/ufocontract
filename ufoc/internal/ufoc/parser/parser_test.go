package parser

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserSimpleVersion(t *testing.T) {
	input := `version 1`

	assertAST(t, input, &File{
		Version: 1,
	})
}

func TestParserEmptyNamespace(t *testing.T) {
	input := `
		version 1
		namespace Tasks {}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
				},
			},
		},
	})
}

func TestParserNamespaceWithDocstring(t *testing.T) {
	input := `
		version 1
		"""
		Documentation for Tasks namespace.
		"""
		namespace Tasks {}
	`

	docstring := "\"\"\"\n\t\tDocumentation for Tasks namespace.\n\t\t\"\"\""
	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Docstring: &docstring,
					Name:      "Tasks",
				},
			},
		},
	})
}

func TestParserSimpleType(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			type Task {
				id: string
				status: int
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Type: &TypeDef{
								Name: "Task",
								Fields: []*Field{
									{
										Name: "id",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
									{
										Name: "status",
										Type: &TypeRef{
											Named: strPtr("int"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserOptionalField(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			type Task {
				id: string
				tags?: string
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Type: &TypeDef{
								Name: "Task",
								Fields: []*Field{
									{
										Name: "id",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
									{
										Name:     "tags",
										Optional: true,
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserArrayType(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			type Task {
				tags: string[]
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Type: &TypeDef{
								Name: "Task",
								Fields: []*Field{
									{
										Name: "tags",
										Type: &TypeRef{
											Named: strPtr("string"),
											Array: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserEnum(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			enum TaskStatus {
				PENDING
				RUNNING
				COMPLETED
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Enum: &EnumDef{
								Name: "TaskStatus",
								Members: []*EnumMember{
									{Name: "PENDING"},
									{Name: "RUNNING"},
									{Name: "COMPLETED"},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserEnumWithValues(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			enum ErrorCode: int {
				UNKNOWN = 1
				TIMEOUT = 100
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Enum: &EnumDef{
								Name:     "ErrorCode",
								BaseType: strPtr("int"),
								Members: []*EnumMember{
									{
										Name: "UNKNOWN",
										Value: &Value{
											Number: strPtr("1"),
										},
									},
									{
										Name: "TIMEOUT",
										Value: &Value{
											Number: strPtr("100"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserConst(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			const MaxRetries: int = 5
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Const: &ConstDef{
								Name: "MaxRetries",
								Type: &TypeRef{
									Named: strPtr("int"),
								},
								Value: &Value{
									Number: strPtr("5"),
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserPattern(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			pattern TaskTopic = "{ns}.{taskId}.updates"
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Pattern: &PatternDef{
								Name:    "TaskTopic",
								Pattern: "\"{ns}.{taskId}.updates\"",
							},
						},
					},
				},
			},
		},
	})
}

func TestParserDeprecated(t *testing.T) {
	input := `
		version 1
		namespace Tasks {
			deprecated("Use NewTask instead")
			type OldTask {
				id: string
			}
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Namespace: &Namespace{
					Name: "Tasks",
					Children: []*NamespaceChild{
						{
							Type: &TypeDef{
								Deprecated: strPtr("\"Use NewTask instead\""),
								Name:       "OldTask",
								Fields: []*Field{
									{
										Name: "id",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestParserCompleteExample(t *testing.T) {
	input := `
		version 1

		""" This is a standalone docstring. """

		namespace FooBar {}

		// This is a standalone comment.

		/*
			This is a standalone block comment.
		*/

		"""
		This file defines all data contracts
		for the Tasks domain.
		"""
		namespace Tasks {

			""" This is a standalone docstring. """

			// This is a standalone comment.

			/*
				This is a standalone block comment.
			*/

			"""
			Maximum number of retries for a task.
			"""
			const MaxRetries: int = 5

			"""
			Defines the possible states of a Task.
			"""
			enum TaskStatus {
				PENDING
				RUNNING
				COMPLETED
				FAILED
			}

			"""
			Payload for creating a new task.
			"""
			type CreateTaskPayload {
				""" Correlation ID for tracking. """
				correlationId: string

				""" The name of the task to execute. """
				taskName: string
			}

			"""
			Defines the structure of a stored Task.
			"""
			type Task {
				""" Unique ID of the task. """
				id: string

				""" Current status of the task. """
				status: TaskStatus

				""" Optional tags for filtering. """
				tags?: string[]
			}

			"""
			Topic for updates on a specific task.
			"""
			pattern TaskUpdatesTopic = "{ns}.{taskId}.updates"
		}
	`

	assertAST(t, input, &File{
		Version: 1,
		Children: []*FileChild{
			{
				Docstring: &Docstring{Text: "\"\"\" This is a standalone docstring. \"\"\"", BlankLine: true},
			},
			{
				Namespace: &Namespace{
					Name: "FooBar",
				},
			},
			{
				Comment: &Comment{Text: "// This is a standalone comment."},
			},
			{
				BlockComment: &BlockComment{Text: "/*\n\t\t\tThis is a standalone block comment.\n\t\t*/"},
			},
			{
				Namespace: &Namespace{
					Docstring: strPtr("\"\"\"\n\t\tThis file defines all data contracts\n\t\tfor the Tasks domain.\n\t\t\"\"\""),
					Name:      "Tasks",
					Children: []*NamespaceChild{
						{
							Docstring: &Docstring{Text: "\"\"\" This is a standalone docstring. \"\"\"", BlankLine: true},
						},
						{
							Comment: &Comment{Text: "// This is a standalone comment."},
						},
						{
							BlockComment: &BlockComment{Text: "/*\n\t\t\t\tThis is a standalone block comment.\n\t\t\t*/"},
						},
						{
							Const: &ConstDef{
								Docstring: strPtr("\"\"\"\n\t\t\tMaximum number of retries for a task.\n\t\t\t\"\"\""),
								Name:      "MaxRetries",
								Type: &TypeRef{
									Named: strPtr("int"),
								},
								Value: &Value{
									Number: strPtr("5"),
								},
							},
						},
						{
							Enum: &EnumDef{
								Docstring: strPtr("\"\"\"\n\t\t\tDefines the possible states of a Task.\n\t\t\t\"\"\""),
								Name:      "TaskStatus",
								Members: []*EnumMember{
									{Name: "PENDING"},
									{Name: "RUNNING"},
									{Name: "COMPLETED"},
									{Name: "FAILED"},
								},
							},
						},
						{
							Type: &TypeDef{
								Docstring: strPtr("\"\"\"\n\t\t\tPayload for creating a new task.\n\t\t\t\"\"\""),
								Name:      "CreateTaskPayload",
								Fields: []*Field{
									{
										Docstring: strPtr("\"\"\" Correlation ID for tracking. \"\"\""),
										Name:      "correlationId",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
									{
										Docstring: strPtr("\"\"\" The name of the task to execute. \"\"\""),
										Name:      "taskName",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
								},
							},
						},
						{
							Type: &TypeDef{
								Docstring: strPtr("\"\"\"\n\t\t\tDefines the structure of a stored Task.\n\t\t\t\"\"\""),
								Name:      "Task",
								Fields: []*Field{
									{
										Docstring: strPtr("\"\"\" Unique ID of the task. \"\"\""),
										Name:      "id",
										Type: &TypeRef{
											Named: strPtr("string"),
										},
									},
									{
										Docstring: strPtr("\"\"\" Current status of the task. \"\"\""),
										Name:      "status",
										Type: &TypeRef{
											Named: strPtr("TaskStatus"),
										},
									},
									{
										Docstring: strPtr("\"\"\" Optional tags for filtering. \"\"\""),
										Name:      "tags",
										Optional:  true,
										Type: &TypeRef{
											Named: strPtr("string"),
											Array: true,
										},
									},
								},
							},
						},
						{
							Pattern: &PatternDef{
								Docstring: strPtr("\"\"\"\n\t\t\tTopic for updates on a specific task.\n\t\t\t\"\"\""),
								Name:      "TaskUpdatesTopic",
								Pattern:   "\"{ns}.{taskId}.updates\"",
							},
						},
					},
				},
			},
		},
	})
}

/*******************
* HELPER FUNCTIONS *
*******************/

func strPtr(s string) *string {
	return &s
}

func assertAST(t *testing.T, input string, expected *File) {
	t.Helper()

	ast, err := Parser.ParseString("", input)
	require.NoError(t, err)
	require.NotNil(t, ast)

	stripPositions(ast)
	stripPositions(expected)

	assert.Equal(t, expected, ast)
}

func stripPositions(v any) {
	if v == nil {
		return
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Name == "Pos" && fieldType.Type == reflect.TypeOf(lexer.Position{}) {
			field.Set(reflect.Zero(fieldType.Type))
			continue
		}

		if !field.CanInterface() {
			continue
		}

		switch field.Kind() {
		case reflect.Pointer:
			if !field.IsNil() {
				stripPositions(field.Interface())
			}
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.CanInterface() {
					stripPositions(elem.Interface())
				}
			}
		case reflect.Struct:
			stripPositions(field.Addr().Interface())
		}
	}
}
