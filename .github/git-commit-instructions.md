# Conventional Commit Message Generation Guidelines

## Overview

You are tasked with generating Conventional Commit messages for code changes. These commit messages MUST be in English
and follow the specification below.

## Text

```
<type>[(scope)][!]: <description>

[body]

[optional footer(s)]
```

## Types

- `feat`: A new feature for the user/application (MUST be used when adding new functionality)
- `fix`: A bug fix (MUST be used when fixing bugs)
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit
- `deps`: Changes to dependencies (updating, adding, removing)

## Scope

Scope MUST be a noun describing a section of the codebase surrounded by parentheses.

### Core Components

- `(encoding)` - Text encoding and character set handling
- `(connector)` - Printer connection and communication
- `(printer)` - Main printer functionality
- `(composer)` - Command composition and assembly
- `(poster)` - Library name, used when changes affect the entire library or are broad in scope

### Printing Features

- `(commands)` - ESC/POS commands
- `(executor)` - JSON Commands execution and printing
- `(builder)` - JSON Commands programmatic building
- `(graphics)` - Image and graphics handling
- `(barcode)` - Barcode generation
- `(qrcode)` - QR code generation
- `(bitimage)` - Bitmap image processing
- `(tables)` - Table formatting and layout

### Formatting

- `(character)` - Character formatting
- `(linespacing)` - Line spacing control
- `(print)` - Print operations
- `(printposition)` - Cursor and print position

### Configuration

- `(profiles)` - Printer profiles
- `(protocol)` - Communication protocols
- `(config)` - Configuration management
- `(mechanismcontrol)` - Hardware mechanism control

### Infrastructure

- `(api)` - API endpoints and interfaces
- `(service)` - Service layer
- `(models)` - Data models
- `(errors)` - Error handling
- `(logs)` - Logging functionality
- `(utils)` - Utility functions
- `(common)` - Shared/common code

### Development

- `(github)` - GitHub templates and configurations
- `(connection)` - Connection management (consider merging with `connector`)
- `(gomod)` - Go module management
- `(npm)` - NPM package management
- `(gh-actions)` - GitHub Actions workflows

## Description

- MUST be imperative, present tense: "change" not "changed" nor "changes"
- MUST be lowercase
- MUST NOT end with a period

## Breaking Changes

Breaking changes MUST be indicated in one of two ways:

1. Adding `!` before the colon: `feat(api)!: remove deprecated endpoints`
2. Adding `BREAKING CHANGE:` in the footer:
   `BREAKING CHANGE: environment variables now take precedence over config files`

## Body

- Use the body to explain WHAT and WHY (not HOW)
- Separate paragraphs with blank lines
- Use bullet points with hyphens (`-`)

## Footers

- Footers MUST be separated from the body by a blank line
- Each footer MUST consist of a token, followed by either `: ` or ` #`
- Common footers include:
    - `Fixes: #123`
    - `Reviewed-by: Person Name`
    - `Refs: #456`
    - `BREAKING CHANGE: description of breaking change`

## Instructions for AI

When analyzing changes in this project:

1. **Identify the affected component** from the scope list above
2. **Determine the type** based on the nature of the change:
    - New ESC/POS commands → `feat`
    - Printer connectivity issues → `fix(connector)`
    - Performance optimizations → `perf`
3. **Focus on user impact** in the description, not implementation details
4. **Include technical context** in the body when necessary
5. **Reference issues** using `Fixes: #123` or `Refs: #456`

When presented with code changes, first identify the primary purpose of the change to select the appropriate type and
scope, then follow this Format to generate a conventional commit message.