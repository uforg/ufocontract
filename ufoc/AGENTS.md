# UFOC LLM Agent guidelines

This directory contains the CLI for the UFO Contract project. You MUST follow these guidelines STRICTLY.

## General Guidelines

- All project communication MUST be in English (it includes the code, documentation, and comments).
- Dont overuse comments, use them only when necessary because the code should be self-explanatory.
- Dont overcomplicate the code, keep it simple and readable. Follow the KISS principle.

## Testing Guidelines

- Use the `github.com/stretchr/testify/require` package to assert expectations.
- Use table-driven tests when it makes the test more readable and maintainable.
- When testing a function, only create one single test and for each test case or edge case create a sub-test.
