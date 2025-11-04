# Test Projects

This directory contains sample projects for different test frameworks supported by Aligned. These serve as:

1. **Examples** - Demonstrate how to structure projects for each framework
2. **Manual Testing** - Test Aligned connectors against real projects
3. **Documentation** - Show test discovery output formats

## Available Projects

### vitest-sample

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
# From test-projects/vitest-sample/
npx vitest list --json
npx vitest list

# From repo root with align binary built
./bin/align check test-projects/vitest-sample/
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
