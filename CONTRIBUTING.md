# Contributing to Sarvam Go SDK

Thank you for contributing to the Sarvam Go SDK. Contributions should keep the
client small, idiomatic, well tested, and compatible with the public API.

## Development setup

You will need Go 1.27 or newer. Clone the repository and run the test suite:

```bash
git clone https://github.com/abhinayjangde/sarvam.git
cd sarvam
go test .
```

An API key is not required to run the tests. HTTP behavior should be tested with
`httptest` rather than live requests to the Sarvam API.

## Making changes

1. Create a focused branch from the default branch.
2. Make the smallest change that solves the problem.
3. Add or update tests for behavior that changes.
4. Format changed Go files with `gofmt`.
5. Run the checks listed below before opening a pull request.

Keep exported types, functions, and fields documented. Preserve context
cancellation, configurable HTTP clients, and the existing authentication
behavior unless a change explicitly requires otherwise.

## Checks

Run these commands from the repository root:

```bash
gofmt -w *.go examples/basic/*.go examples/cli-chatbot/*.go
go test .
go vet .
```

Review the resulting diff to make sure formatting did not change unrelated
files.

The `examples/cli-chatbot` example currently imports `github.com/joho/godotenv`,
which is not declared in the module dependencies. If you work on that example,
resolve its dependency before building or testing it.

## Pull requests

Pull requests should include:

- A concise description of the problem and the change.
- Tests covering new or changed behavior.
- Any public API or README documentation updates that are needed.
- A note about checks that were run locally.

Keep pull requests focused. Explain compatibility or behavioral trade-offs in
the description when they affect SDK users.

## Security

Never commit API keys, tokens, or other credentials. If you discover a security
issue, do not open a public issue with sensitive details; contact the project
maintainer privately first.

## License

By contributing, you agree that your contributions will be licensed under the
MIT License used by this project.
