# UFOC LLM Agent guidelines

This directory contains the CLI for the UFO Contract project. You MUST follow these guidelines STRICTLY.

## General Guidelines

- All project communication MUST be in English (it includes the code, documentation, and comments).
- Dont overuse comments, use them only when necessary because the code should be self-explanatory.
- Dont overcomplicate the code, keep it simple and readable. Follow the KISS principle.
- Every time you add a new feature or want to make sure the code is working as expected, run the command `task ci` to run the tests, lint, build and other checks. Make sure all the checks pass before committing your changes.

## Testing Guidelines

- Use the `github.com/stretchr/testify/require` package to assert expectations.
- Use table-driven tests when it makes the test more readable and maintainable.
- When testing a function, only create one single test and for each test case or edge case create a sub-test.
