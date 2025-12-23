# Contributing to Poster

🎉 Thank you for considering contributing to Poster!

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

## 📜 Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- A GitHub account

### Fork and Clone

```bash
# Fork the repository on GitHub, then: 
git clone https://github.com/YOUR_USERNAME/poster.git
cd poster

# Add upstream remote
git remote add upstream https://github.com/adcondev/poster.git
```

### Build and Test

```bash
# Install dependencies
go mod download

# Run tests
go test -v -race ./pkg/...

# Run linter
golangci-lint run

# Build examples
cd examples/basic
go build
```

## 🔄 Development Workflow

### 1. Create a Branch

```bash
# Update your fork
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feat/my-new-feature
```

### 2. Make Changes

- Write clean, readable code
- Add tests for new features
- Update documentation as needed
- Follow existing code patterns

### 3. Test Your Changes

```bash
# Run tests
go test -v -race ./pkg/...

# Run benchmarks
go test -bench=. ./pkg/...

# Check coverage
go test -coverprofile=coverage.txt ./pkg/...
go tool cover -html=coverage.txt
```

## 📏 Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Pass `golangci-lint` checks
- Keep functions small and focused
- Add comments for exported functions

### Example

```go
// ProcessReceipt processes a receipt document and sends it to the printer. 
// It returns an error if the document is invalid or printing fails.
func ProcessReceipt(doc *Document) error {
if err := doc.Validate(); err != nil {
return fmt.Errorf("invalid document: %w", err)
}
// ... implementation
}
```

## 💬 Commit Messages

We follow [Conventional Commits](https://conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test additions or modifications
- `ci`: CI/CD changes
- `chore`: Other changes

### Examples

```bash
feat(graphics): add dithering algorithm for image processing

fix(connection): resolve timeout issue on Windows 11

docs(readme): update installation instructions

test(composer): add unit tests for ESC/POS commands
```

See [Commit Guidelines](.github/git-commit-instructions.md) for complete details.

## 🔀 Pull Request Process

### 1. Before Submitting

- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] PR title follows conventional commits

### 2. Submit PR

1. Push your branch to your fork
2. Open a PR against `main`
3. Fill out the PR template completely
4. Link related issues

### 3. Code Review

- Address review comments
- Keep PR focused and small
- Be responsive to feedback
- Update PR description as needed

### 4. After Merge

```bash
# Update your fork
git checkout main
git pull upstream main
git push origin main

# Delete feature branch
git branch -d feat/my-new-feature
git push origin --delete feat/my-new-feature
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test ./pkg/...

# Run specific package
go test ./pkg/graphics/...

# Run with coverage
go test -coverprofile=coverage.txt ./pkg/... 

# View coverage
go tool cover -html=coverage.txt
```

### Benchmarks

```bash
# Run benchmarks
go test -bench=.  -benchmem ./pkg/... 

# Benchmark specific function
go test -bench=BenchmarkDithering ./pkg/graphics/
```

### Integration Tests

```bash
# Run integration tests (if applicable)
go test -tags=integration ./test/integration/... 
```

## 📚 Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use complete sentences
- Include examples where helpful

```go
// CreateProfile80mm creates a printer profile for standard 80mm thermal printers.
// It configures ESC/POS settings optimized for Epson TM-T88 and compatible models.
//
// Example:
//
//	prof := profile.CreateProfile80mm()
//	printer := service.NewPrinter(composer, prof, conn)
func CreateProfile80mm() *Profile {
// ... 
}
```

### README Updates

- Update README. md for user-facing changes
- Add examples for new features
- Update feature list when applicable

### API Documentation

- API changes should be documented in `/api/v1/DOCUMENT_V1.md`
- Include JSON schema updates
- Provide usage examples

## 🏷️ Issue Labels

We use labels to organize issues:

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `security` - Security-related issues

## ❓ Questions?

- Open a [Discussion](https://github.com/adcondev/poster/discussions)
- Ask in your PR/issue
- Check existing [documentation](README.md)

## 🙏 Thank You!

Your contributions make Poster better for everyone!

---

**Happy coding! ** 🚀