# Contributing to Jone

Thank you for your interest in contributing to **Jone**! We welcome contributions from the community to help make this Go database migration tool even better.

The primary goal for this project is to maintain a stable and reliable tool. As the author puts it: **"Everything should work without breaking what currently works."**

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) 1.21 or higher.
- A SQL database (PostgreSQL for now, as it's the most stable).

### Setup
1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/jone.git
   cd jone
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

---

## 🛠️ Development Workflow

1. **Keep your fork up to date**: Before starting new work, ensure your local `main` branch is synced with the upstream repository.
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```
2. **Create a Branch**: Always create a new branch for your changes from the latest `main`.
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Write Code**: Follow standard Go conventions. Use `gofmt` to format your code.
4. **Write Tests**: This is critical! If you add a feature or fix a bug, you **must** add tests to verify it and ensure no regressions.
5. **Run Tests**: Ensure all tests pass before submitting your PR.
   ```bash
   go test ./...
   ```
6. **Linting**: Run `go vet ./...` to check for common mistakes.

---

## ✅ Contribution Guidelines

### 1. Bug Reports
If you find a bug, please open an issue and include:
- A clear description of the problem.
- Steps to reproduce the bug.
- The expected vs. actual behavior.

### 2. Feature Requests
We are open to new features (like more database dialects or the upcoming Query Builder!). Please open an issue first to discuss the design before starting work.

### 3. Pull Requests
- Keep PRs focused. One PR should ideally solve one problem or add one feature.
- Ensure your PR includes updated documentation if applicable.
- **Tests are mandatory**. PRs without tests or with failing tests will not be merged.
- Link the PR to any relevant issues.

---

## 📜 Code Style
- Use standard Go formatting (`gofmt`).
- Use descriptive variable and function names.
- Document exported functions and types using Go doc comments.

---

## ⚖️ License
By contributing to Jone, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
