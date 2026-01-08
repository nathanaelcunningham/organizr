# Contributing to Organizr

Thank you for your interest in contributing to Organizr! This guide will help you get started with development, understand our conventions, and submit contributions.

---

## 1. Getting Started

### Prerequisites

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Node 20+** - [Download](https://nodejs.org/)
- **qBittorrent** - [Download](https://www.qbittorrent.org/download.php)
  - Enable Web UI in qBittorrent settings
  - Note your Web UI port and credentials

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/organizr.git
cd organizr

# Backend setup
cd backend
make build    # Build the backend
make run      # Start the backend server

# Frontend setup (in a new terminal)
cd frontend
npm install           # Install dependencies
npm run dev           # Start development server
```

### Running Tests

```bash
# Backend tests
cd backend
make test                    # Run all tests
go test -race ./...          # Run with race detection

# Frontend tests
cd frontend
npm test                     # Run Vitest tests
```

---

## 2. Project Structure Overview

```
organizr/
├── backend/          # Go API server
│   ├── cmd/api/      # Application entry point
│   └── internal/     # Business logic, services, handlers
│       ├── config/         # Configuration service
│       ├── downloads/      # Download domain
│       ├── models/         # Domain models
│       ├── persistence/    # Repository pattern
│       ├── qbittorrent/    # qBittorrent client
│       ├── search/         # Search service
│       └── server/         # HTTP handlers, routes
├── frontend/         # React SPA
│   └── src/
│       ├── api/            # HTTP client
│       ├── components/     # React components
│       ├── pages/          # Page components
│       ├── stores/         # Zustand state stores
│       ├── types/          # TypeScript definitions
│       └── utils/          # Utility functions
└── docs/             # Documentation
    └── architecture/ # Architecture Decision Records
```

**For detailed documentation:**
- **Architecture decisions:** See [`docs/architecture/ADR.md`](docs/architecture/ADR.md)
- **Codebase structure:** See [`.planning/codebase/STRUCTURE.md`](.planning/codebase/STRUCTURE.md)
- **Conventions:** See [`.planning/codebase/CONVENTIONS.md`](.planning/codebase/CONVENTIONS.md)
- **Architecture patterns:** See [`.planning/codebase/ARCHITECTURE.md`](.planning/codebase/ARCHITECTURE.md)

---

## 3. Development Workflow

1. **Create a feature branch** from `main`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding conventions (see below)

3. **Run tests locally**
   ```bash
   # Backend
   cd backend && make test

   # Frontend
   cd frontend && npm test
   ```

4. **Ensure race detection passes** (critical for concurrent code)
   ```bash
   cd backend && go test -race ./...
   ```

5. **Commit with conventional commit format**
   ```bash
   git commit -m "feat(scope): add new feature"
   ```

6. **Push and create a pull request**
   ```bash
   git push origin feature/your-feature-name
   ```

---

## 4. Coding Conventions

### Backend (Go)

- **Follow standard Go conventions** - use `gofmt`, run `go vet`
- **Repository pattern** - Data access through repository interfaces
- **Service layer** - Business logic in service structs (ConfigService, DownloadService, etc.)
- **Error handling** - Always check and wrap errors: `return fmt.Errorf("context: %w", err)`
- **Concurrency** - Use goroutines with context.Context for cancellation
- **Testing** - Test files alongside source (`*_test.go`), table-driven tests preferred

**Example:**
```go
// Service with injected dependencies
func NewDownloadService(repo persistence.DownloadRepository, qbt *qbittorrent.Client) *DownloadService {
    return &DownloadService{
        repo: repo,
        qbt:  qbt,
    }
}

// Error wrapping
if err != nil {
    return fmt.Errorf("failed to create download: %w", err)
}
```

### Frontend (TypeScript)

- **ESLint rules enforced** - Run `npm run lint` before committing
- **2-space indentation** - Consistent across all TypeScript files
- **PascalCase for components** - `SearchBar.tsx`, `DownloadCard.tsx`
- **camelCase for functions/variables** - `formatFileSize`, `isLoading`
- **Zustand for state** - Create stores in `src/stores/`
- **API clients** - Centralized in `src/api/` with typed responses

**Example:**
```typescript
// Component naming
export function SearchBar() {
  // camelCase for variables
  const [searchQuery, setSearchQuery] = useState('');

  // API calls through centralized client
  const results = await searchApi.search(searchQuery);
}
```

### Import Organization

**Backend (Go):**
1. Standard library imports
2. External packages
3. Internal packages
4. Blank line between groups, alphabetical within groups

**Frontend (TypeScript):**
1. External packages (`react`, `zustand`, etc.)
2. Internal modules (relative imports)
3. Type imports

### Comments

- **Explain why, not what** - Code should be self-explanatory
- **Document business logic** - Explain non-obvious decisions
- **Document edge cases** - Note special handling or assumptions

---

## 5. Testing Standards

### Backend Testing

- **Unit tests alongside source** - `*_test.go` files next to implementation
- **Mock dependencies via interfaces** - Use repository/client interfaces for testability
- **Table-driven tests** - Test multiple cases efficiently
- **Race detection required** - All concurrent code must pass `go test -race`

**Example:**
```go
func TestDownloadService_CreateDownload(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateDownloadRequest
        wantErr bool
    }{
        {"valid request", CreateDownloadRequest{...}, false},
        {"invalid URL", CreateDownloadRequest{...}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Use mock repository
            mockRepo := &MockDownloadRepository{}
            service := NewDownloadService(mockRepo, mockQbt)

            err := service.CreateDownload(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("want error %v, got %v", tt.wantErr, err)
            }
        })
    }
}
```

### Frontend Testing

- **Vitest for unit tests** - Test stores, utilities, and components
- **Test stores and utilities** - Focus on logic, not implementation
- **Avoid testing implementation details** - Test behavior, not internals

**Example:**
```typescript
import { describe, it, expect } from 'vitest';
import { formatFileSize } from './formatters';

describe('formatFileSize', () => {
  it('formats bytes correctly', () => {
    expect(formatFileSize(1024)).toBe('1.0 KB');
    expect(formatFileSize(1048576)).toBe('1.0 MB');
  });
});
```

### Coverage Focus

Focus test coverage on **critical paths:**
- HTTP handlers (API contract verification)
- Services (business logic)
- File organization logic
- qBittorrent integration
- State management stores

---

## 6. Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/) with phase-plan scoping:

**Format:** `type(scope): description`

### Types

- `feat` - New feature, endpoint, component, functionality
- `fix` - Bug fix, error correction
- `docs` - Documentation changes
- `refactor` - Code cleanup, no behavior change
- `test` - Test-only changes
- `chore` - Config, tooling, dependencies

### Scope

Use phase-plan format for work within milestone phases (e.g., `10-01`, `11-02`) or general scope for standalone work.

### Examples

```bash
# New feature
git commit -m "feat(10-01): add ADR documentation"

# Bug fix
git commit -m "fix(08-02): correct race condition in monitor"

# Documentation
git commit -m "docs(10-01): add contribution guidelines"

# Refactoring
git commit -m "refactor(07-02): extract series grouping logic"

# Test addition
git commit -m "test(06-01): add handler integration tests"

# Dependency update
git commit -m "chore: update Go dependencies"
```

**Commit body (optional):**
```bash
git commit -m "feat(08-02): add batch download endpoint

- Accepts array of download requests
- Returns partial success with detailed errors
- Enforces 50-item limit
"
```

---

## 7. Pull Request Guidelines

### PR Title

Follow commit message format: `type(scope): description`

### PR Description

- **What:** Describe the changes clearly
- **Why:** Explain the motivation or problem being solved
- **Testing:** Describe how you tested the changes
- **Related issues:** Link to relevant issues (e.g., "Closes #123")

### Before Submitting

- ✅ All tests pass (`make test` for backend, `npm test` for frontend)
- ✅ Race detection passes (`go test -race ./...` for backend)
- ✅ No linting errors (`npm run lint` for frontend)
- ✅ Documentation updated if needed
- ✅ Commit messages follow conventional format

### Example PR

**Title:** `feat(11-02): standardize error response format`

**Description:**
```markdown
## What
Standardizes all API error responses to use consistent ErrorResponse structure.

## Why
Improves frontend error handling by providing predictable error format.

## Changes
- Created ErrorResponse type with code, message, details fields
- Updated all handlers to use standardized error responses
- Added tests for error response format

## Testing
- All existing handler tests pass
- Added new tests for error response structure
- Manually tested error scenarios in frontend

Closes #45
```

---

## 8. Code Review Expectations

### For Authors

- **Be responsive** - Address review comments promptly
- **Explain decisions** - Provide context for non-obvious changes
- **Accept feedback** - Reviews improve code quality for everyone

### For Reviewers

- **Check for correctness** - Does the code do what it's supposed to?
- **Verify readability** - Is the code clear and maintainable?
- **Test coverage** - Are critical paths tested?
- **Security considerations** - Path sanitization, error handling, input validation
- **Performance considerations** - No N+1 queries, avoid unnecessary allocations
- **Adherence to conventions** - Follow established patterns

### Review Focus Areas

1. **Correctness** - Logic errors, edge cases, error handling
2. **Testing** - Adequate test coverage, race detection for concurrent code
3. **Security** - Input validation, path sanitization, no information leakage
4. **Performance** - Efficient queries, minimal allocations, no obvious bottlenecks
5. **Maintainability** - Clear naming, appropriate comments, follows conventions

---

## 9. Getting Help

### Documentation Resources

- **Architecture:** [`docs/architecture/ADR.md`](docs/architecture/ADR.md) - Technical decisions and rationale
- **Codebase structure:** [`.planning/codebase/STRUCTURE.md`](.planning/codebase/STRUCTURE.md) - Directory layout and organization
- **Patterns:** [`.planning/codebase/ARCHITECTURE.md`](.planning/codebase/ARCHITECTURE.md) - Architecture patterns and data flow
- **Conventions:** [`.planning/codebase/CONVENTIONS.md`](.planning/codebase/CONVENTIONS.md) - Coding standards and style

### Getting Unstuck

- **Check existing code** - Look for similar implementations in the codebase
- **Read the tests** - Tests document expected behavior
- **Review ADR** - Understand the rationale behind decisions
- **Open an issue** - Ask questions or request clarification

### Communication

- **Be respectful** - We're all here to build something great
- **Be clear** - Provide context and examples when asking questions
- **Be patient** - Maintainers may have limited availability

---

## Welcome!

We're excited to have you contribute to Organizr! Whether you're fixing a bug, adding a feature, improving documentation, or helping with code review, your contributions make this project better for everyone.

If you have questions or need help, don't hesitate to open an issue. Happy coding! 🎉
