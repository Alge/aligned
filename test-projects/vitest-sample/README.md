# Vitest Sample Project

This is a minimal Vitest project used for testing and demonstrating Aligned's Vitest connector.

## Setup

```bash
npm install
```

## Running Tests

```bash
# Run tests
npm test

# List tests (what Aligned uses)
npx vitest list
npx vitest list --json
```

## Test Structure

The project demonstrates:
- **Nested describe blocks** - Shows hierarchy in test names
- **Multiple test suites** - Math operations and String operations
- **ES modules** - Using `type: "module"` in package.json

## Test Discovery Output

Running `npx vitest list --json` produces:

```json
[
  {
    "name": "Math operations > Addition > adds 1 + 2 to equal 3",
    "file": "/absolute/path/to/src/example.test.js"
  },
  ...
]
```

Aligned converts this to:
```
src/example.test.js > Math operations > Addition > adds 1 + 2 to equal 3
```

## Using with Aligned

From the repository root:

```bash
# Build Aligned
make build

# Initialize Aligned config for this project
cd test-projects/vitest-sample
../../bin/align init javascript-vitest .

# Update .align.yml to use local vitest installation
# Change executable from "vitest" to "node_modules/.bin/vitest"

# Check specification coverage
../../bin/align check spec.md
```

**Note:** Vitest is typically installed locally in `node_modules/.bin/`, not globally. The `.align.yml` file should specify the relative path `node_modules/.bin/vitest` as the executable.
