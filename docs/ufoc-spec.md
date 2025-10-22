# UFO Contract DSL Specification (.ufoc)

## 1. Overview

The UFO Contract DSL (.ufoc) is a Domain-Specific Language (DSL) designed to define data contracts, enumerations, constants, and communication patterns for multi-language code generation.  
The primary goal of UFO Contract is to offer an intuitive, human-readable format that ensures the best possible Developer Experience (DX), maintaining consistency across different services and applications (e.g., Go, TypeScript).

## 2. General Syntax (.ufoc)

This is the base syntax for the UFO Contract DSL.

```text
version <number>

// Single-line comment

/*
  Multi-line
  comment
*/

"""
Documentation for the namespace.
"""
namespace "NamespaceName" {
  """
  <Type Documentation>
  """
  type TypeName {
    """ <Field Documentation> """
    fieldName[?]: <Type>
  }

  """
  <Enum Documentation>
  """
  enum EnumName[: <BaseType>] {
    MEMBER_ONE [= <value>]
    MEMBER_TWO [= <value>]
  }

  """
  <Constant Documentation>
  """
  const ConstantName: <Type> = <Value>

  """
  <Pattern Documentation>
  """
  pattern PatternName = "<string_pattern_with_{placeholderName}>"
}
```

## 3. Namespaces

All definitions (types, enums, consts, patterns) must be contained within a namespace block. This is the top-level logical grouper.  
It controls the grouping of related definitions (usually by business domain).  
Generates: Packages (Go), Modules/Files (TS), Documentation Sections.

```text
"""
Defines all contracts related to
task management.
"""
namespace "Tasks" {
  // ... all other definitions go here ...
}
```

## 4. Types

Types are the building blocks of your data contracts. They define the structure of the data being exchanged (e.g., DTOs, payloads).

### 4.1 Primitive Types

Primitive types are the types built into the DSL.

| DSL      | JSON Type | Description                           |
| -------- | --------- | ------------------------------------- |
| string   | string    | UTF-8 text string                     |
| int      | integer   | 64-bit integer                        |
| float    | number    | 64-bit floating-point number          |
| bool     | boolean   | Boolean value (true or false)         |
| datetime | string    | Date and time value (ISO 8601 format) |

### 4.2 Composite Types

Composite types are types composed of other types.

```text
// Array
ElementType[]  // E.g.: string[]

// Inline object
{
  fieldOne: Type
  fieldTwo: Type
}
```

### 4.3 Custom Types

You can define custom types that can be reused throughout your contract.

```text
"""
<Type Documentation>
"""
type TypeName {
  """ <Field Documentation> """
  fieldName[?]: <Type>
}
```

#### 4.3.1 Custom Type Documentation

You can add documentation to your custom types. They support Markdown syntax that will be rendered in the generated documentation (see Section 8).

#### 4.3.2 Type Composition

To reuse fields from other types, use composition by including the type as a field:

```text
type BaseEntity {
  id: string
  createdAt: datetime
  updatedAt: datetime
}

type User {
  base: BaseEntity
  email: string
  name: string
}
```

#### 4.3.3 Optional Fields

All fields in a type are required by default. To make a field optional, use the ? suffix.

```text
// Optional field
fieldName?: Type
```

#### 4.3.4 Field Documentation

You can add documentation to your fields to help the developer understand their use.

```text
type User {
  """ The user's email address """
  email: string

  """ The user's full name """
  name: string
}
```

## 5. Enums

Controls type-safe sets of named values (e.g., states, categories).

### 5.1 Overview

- Enums can optionally define a base type: `string` or `int`.
- If a base type is omitted, it defaults to `string`.

### 5.2 Rules for the `string` base (default or explicit)

- Members without an assigned value automatically use their own name as a string literal (e.g., `PENDING` becomes `"PENDING"`).
- Members can be assigned a custom string literal.

### 5.3 Rules for the `int` base

- If the base type is `int`, all members must be explicitly assigned an integer value.
- No implicit numeric assignment is allowed when the base is `int`.

### 5.4 Examples

```text
"""
Defines the possible states of a Task.
Defaults to 'string'. Values are string literals.
"""
enum TaskStatus {
  PENDING   // = "PENDING"
  RUNNING   // = "RUNNING"
  COMPLETED // = "COMPLETED"
  FAILED    // = "FAILED"
}

"""
Defines payment methods using custom string values.
"""
enum PaymentMethod: string {
  CREDIT_CARD = "credit_card"
  PAYPAL = "paypal"
  BANK_TRANSFER = "bank_transfer"
}

"""
Defines specific error codes.
"""
enum ErrorCode: int {
  UNKNOWN = 1
  TIMEOUT = 100
  INVALID_AUTH = 101
}
```

## 6. Constants

Controls static literal values (numbers, strings).

```text
""" Maximum number of retries for a task. """
const MaxRetries: int = 5

""" Queue where failed tasks are published. """
const ErrorQueue: string = "tasks.failed.queue"
```

## 7. String Patterns

Controls the definition of static or dynamic strings, commonly used for messaging topics, NATS subjects, API routes, etc.  
Placeholders are defined using {camelCaseName}.  
Special reserved placeholders, {namespace} and {ns}, will be automatically replaced by the name of the namespace the pattern is defined in.

```text
// Static Pattern
""" Topic for broadcasting events to all workers. """
pattern BroadcastTopic = "tasks.broadcast"

// Dynamic Pattern
"""
Topic for updates on a specific task.
{ns} will be replaced with "Tasks" (from the parent namespace).
Resulting pattern: "Tasks.{taskId}.updates"
"""
pattern TaskTopic = "{ns}.{taskId}.updates"
```

## 8. Documentation (Docstrings)

### 8.1 Docstrings

Docstrings ("""...""") can be used in two ways: associated with specific elements (namespace, type, enum, const, pattern, or fields) or as standalone documentation.

#### Associated Docstrings

Placed immediately before a definition.

```text
"""
This is documentation for MyType.
"""
type MyType {
  """ This is documentation for myField. """
  myField: string
}
```

#### Standalone Docstrings

Provide general documentation. Ensure there is at least one blank line between the docstring and any following element.

```text
namespace "MyApi" {
  """
  This is general documentation for the 'types' section.
  It can include Markdown.
  """

  // At least one blank line here

  type MyType {
    // ...
  }
}
```

### 8.2 Multi-line Docstrings and Normalization

Docstrings support Markdown. The DSL automatically normalizes indentation: the indentation of the first non-empty line is considered the baseline and is stripped from all other lines, preserving the relative indentation of the Markdown.

### 8.3 External Documentation Files

For extensive documentation, you can reference external Markdown files:

```text
version 1

namespace "Tasks" {
  // Standalone documentation
  """ ./docs/tasks-overview.md """

  // Associated documentation
  """ ./docs/task-payload.md """
  type TaskPayload {
    // ...
  }
}
```

## 9. Deprecation

UFO Contract provides a mechanism to mark definitions as deprecated.

### 9.1 Basic Deprecation

Use the deprecated keyword before the definition:

```text
deprecated type OldTask {
  // ...
}

deprecated enum OldStatus {
  // ...
}
```

### 9.2 Deprecation with Message

Provide additional information in parentheses:

```text
deprecated("Replaced by TaskPayloadV2")
type TaskPayload {
  // ...
}

deprecated("Use ErrorQueue instead")
const FailureQueue: string = "tasks.failed"
```

### 9.3 Placement

The deprecated keyword must be placed between any docstring and the element definition.

```text
"""
Documentation for TaskPayload
"""
deprecated("Use TaskPayloadV2 instead")
type TaskPayload {
  // ...
}
```

## 10. Complete Example (.ufoc)

```text
version 1

"""
This file defines all data contracts
and messaging patterns for the 'Tasks' domain.
"""
namespace "Tasks" {

  // --- Standalone Documentation ---
  """
  ## Data Types (DTOs)

  These are the primary payloads used in
  our tasking system.
  """

  // --- Constants ---

  """ Maximum number of retries for a task. """
  const MaxRetries: int = 5

  """ Queue where failed tasks are published. """
  const ErrorQueue: string = "tasks.failed.queue"

  // --- Enumerations ---

  """
  Defines the possible states of a Task.
  Defaults to 'string' type. Values will be "PENDING", "RUNNING", etc.
  """
  enum TaskStatus {
    PENDING
    RUNNING
    COMPLETED
    FAILED
  }

  // --- Types ---

  """
  Payload for creating a new task.
  """
  type CreateTaskPayload {
    """ Correlation ID for tracking. """
    correlationId: string

    """ The name of the task to execute. """
    taskName: string

    """ Arbitrary data for the task (e.g., stringified JSON). """
    payload: string
  }

  """
  Defines the structure of a stored Task.
  """
  type Task {
    """ Unique ID of the task. """
    id: string

    """ Current status of the task. """
    status: TaskStatus

    """ Number of execution attempts. """
    attempts: int

    """ The original payload used to create the task. """
    payload: CreateTaskPayload

    """ Optional tags for filtering. """
    tags?: string[]
  }

  // --- Messaging Patterns ---

  """
  Topic for broadcasting general events to all workers.
  (Generates a constant)
  """
  pattern BroadcastTopic = "tasks.broadcast"

  """
  Topic for updates on a specific task.
  (Generates a builder function)
  {ns} will be replaced with "Tasks".
  """
  pattern TaskUpdatesTopic = "{ns}.{taskId}.updates"

  """
  Obsolete topic, use TaskUpdatesTopic instead.
  """
  deprecated("Use TaskUpdatesTopic instead")
  pattern OldTaskTopic = "{namespace}.{taskId}.status"
}
```

## 11. Known Limitations

- DSL keywords (e.g., type, namespace) cannot be used as identifiers.
- Circular type dependencies are not allowed.
