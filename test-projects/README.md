# Test Projects

This directory contains sample projects for different test frameworks supported by Aligned. These serve as:

1. **Examples** - Demonstrate how to structure projects for each framework
2. **Manual Testing** - Test Aligned connectors against real projects
3. **Documentation** - Show test discovery output formats

## Available Projects

### go-test-sample

**Language:** Go
**Connector:** `go-test`
**Test Framework:** Go's built-in `testing` package

A minimal Go project demonstrating:
- Package-based test organization
- Standard Go test conventions
- Table-driven test patterns

**Setup:**
```bash
cd test-projects/go-test-sample
go mod download
```

**Test Aligned:**
```bash
go test -list=. ./...
../../bin/align check spec.md
```

### pytest-sample

**Language:** Python
**Connector:** `python-pytest`
**Test Framework:** pytest

A minimal Python project demonstrating:
- Class-based test organization
- Pytest fixtures and assertions
- Exception testing patterns

**Setup:**
```bash
cd test-projects/pytest-sample
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

**Test Aligned:**
```bash
pytest --collect-only -q
../../bin/align check spec.md
```

### exunit-sample

**Language:** Elixir
**Connector:** `elixir-exunit`
**Test Framework:** ExUnit

A minimal Elixir project demonstrating:
- ExUnit describe blocks
- Pattern matching in tests
- Standard Mix project structure

**Setup:**
```bash
cd test-projects/exunit-sample
mix deps.get
mix compile
```

**Test Aligned:**
```bash
mix test --list
../../bin/align check spec.md
```

### gleeunit-sample

**Language:** Gleam
**Connector:** `gleam-gleeunit`
**Test Framework:** gleeunit

A minimal Gleam project demonstrating:
- Public test function conventions (`_test` suffix)
- Gleam's type system in tests
- Result type patterns

**Setup:**
```bash
cd test-projects/gleeunit-sample
gleam deps download
gleam build
```

**Test Aligned:**
```bash
gleam test
../../bin/align check spec.md
```

### vitest-sample

**Language:** JavaScript/TypeScript
**Connector:** `javascript-vitest`
**Test Framework:** Vitest

A minimal Vitest project demonstrating:
- ES module configuration
- Nested describe blocks
- Multiple test files structure
- Test name hierarchy

**Setup:**
```bash
cd test-projects/vitest-sample
npm install
```

**Test Aligned:**
```bash
npx vitest list --json
../../bin/align check spec.md
```

## Adding New Test Projects

When adding support for a new test framework, include a minimal sample project here following this structure:

```
test-projects/
  framework-sample/
    .gitignore          # Exclude dependencies
    package.json/etc    # Framework config
    src/                # Example tests
    README.md           # Setup instructions
```
