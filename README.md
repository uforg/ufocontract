# UFO Contract

The universal source of truth for your applications. Define types, constants, and string patterns in one place, generate them everywhere. Modern, cross-language, and built for developer experience.

<p>
  <a href="https://github.com/uforg/ufocontract/actions/workflows/ci.yaml?query=branch%3Amain">
    <img src="https://github.com/uforg/ufocontract/actions/workflows/ci.yaml/badge.svg" alt="CI Status"/>
  </a>
  <a href="https://github.com/uforg/ufocontract/releases/latest">
    <img src="https://img.shields.io/github/release/uforg/ufocontract.svg" alt="Release Version"/>
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/uforg/ufocontract.svg" alt="License"/>
  </a>
  <a href="https://github.com/uforg/ufocontract">
    <img src="https://img.shields.io/github/stars/uforg/ufocontract?style=flat&label=github+stars"/>
  </a>
</p>

> [!WARNING]  
> `UFO Contract` is currently in active development and may introduce breaking changes until reaching a stable v1.0 release. Feedback and contributions are welcome!

## What is UFO Contract?

In any project with multiple services, applications, or even a frontend and backend, keeping shared data structures, constants, and communication contracts in sync is a major challenge. Data becomes inconsistent, and simple changes can lead to breaking bugs.

`UFO Contract` solves this problem by providing a single, human-readable **Domain-Specific Language (DSL)** in `.ufoc` files. These files act as the **universal source of truth** for your data contracts.

A single command-line tool (`ufoc`) compiles these contracts into type-safe code for all your applications, ensuring consistency everywhere.

- **Define** contracts in a simple, intuitive DSL.
- **Compile** contracts into native code for your stack.
- **Explore** contracts in a self-generated documentation playground.

`UFO Contract` generates code for:

- **Go** (Structs, Constants, Builder Functions)
- **TypeScript** (Interfaces, Enums, Consts, Builder Functions)
- **Static HTML Playground** (A searchable website to explore your contracts)

## Key Features

- **Simple, Human-Readable DSL:** A syntax designed for clarity, not complexity.
- **Type-Safe Generation:** Define `types` (structs) in your DSL and get native, type-safe `structs` in Go and `interfaces` in TypeScript.
- **Centralized Constants:** Define `consts` (like error codes or retry counts) and use them across your entire stack.
- **String Pattern Builders:** Define `patterns` (like `tasks.{taskID}.updates`) to generate type-safe builder functions, eliminating magic strings for topics, subjects, or routes.
- **Enumerations:** First-class `enum` support (string or int-based) for defining states and categories.
- **Self-Generating Documentation:** Generates a static HTML playground from your DSL, providing a "living" documentation site that never goes out of date.
- **Single Binary, Zero Dependencies:** The `ufoc` compiler is a single executable. No complex tooling, no `protoc` plugins, no dependency hell.

## Use Cases

`UFO Contract` is a general-purpose tool for maintaining consistency. It is ideal for:

- **Event-Driven Architectures:** Share `type` payloads and `pattern` topics for event queues (like NATS, Kafka, or RabbitMQ) between your publisher and subscriber services.
- **Frontend/Backend Contracts:** Define the `type` for your API DTOs (Data Transfer Objects) and generate identical interfaces for your Go backend and TypeScript frontend.
- **Centralized Configuration:** Define `const` values (e.g., `MaxRetries: int = 5`, `DefaultPageSize: int = 25`) and use the same static values across multiple applications.
- **Shared API Constants:** Define `enum` definitions (like error codes or status types) and `pattern` builders for API routes, ensuring all services use the correct values.
- **Team Documentation:** The generated static playground acts as the single source of truth for all data contracts in the organization, replacing outdated wikis.

## Example: `tasks.ufoc`

Here is a quick look at the `UFO Contract` DSL in action.

```ufoc
version 1

"""
This namespace defines all contracts for
the asynchronous task management system.
"""
namespace "Tasks" {

    """ Possible states for a Task. """
    enum TaskStatus {
        PENDING
        RUNNING
        COMPLETED
        FAILED
    }

    """ Payload for creating a new task. """
    type CreateTaskPayload {
        """ The name of the task to execute. """
        taskName: string

        """ Arbitrary data for the task (stringified JSON). """
        payload: string
    }

    """ Maximum number of retries for any task. """
    const MaxRetries: int = 3

    """
    The NATS subject pattern for task-specific updates.
    (Generates a builder function: BuildTaskUpdatesTopic(taskID string))
    """
    pattern TaskUpdatesTopic = "tasks.{taskID}.updates"

}
```

## Contributing

Contributions are welcome\! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the **AGPL-3.0 License**. See the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.
