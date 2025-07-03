# Contributing to Graph

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

> And if you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues

## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Your First Code Contribution](#your-first-code-contribution)


## I Have a Question

> If you want to ask a question, we assume that you have read the available documentation.

Before you ask a question, it is best to search for existing [Issues](https://github.com/lexkrstn/go-graph/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Issue](https://github.com/lexkrstn/go-graph/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions, depending on what seems relevant.

We will then take care of the issue as soon as possible.



## I Want To Contribute

> ### Legal Notice
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

### Reporting Bugs

#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the documentation. If you are looking for support, you might want to check [this section](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/lexkrstn/go-graphissues?q=label%3Abug).
- Also make sure to search the internet (including Stack Overflow) to see if users outside of the GitHub community have discussed the issue.
- Collect information about the bug:
- Stack trace (Traceback)
- OS, Platform and Version (Windows, Linux, macOS, x86, ARM)
- Version of the interpreter, compiler, SDK, runtime environment, package manager, depending on what seems relevant.
- Possibly your input and the output
- Can you reliably reproduce the issue? And can you also reproduce it with older versions?

#### How Do I Submit a Good Bug Report?

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/lexkrstn/go-graph/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
- Provide the information you collected in the previous section.

Once it's filed:

- The project team will label the issue accordingly.
- A team member will try to reproduce the issue with your provided steps. If there are no reproduction steps or no obvious way to reproduce the issue, the team will ask you for those steps and mark the issue as `needs-repro`. Bugs with the `needs-repro` tag will not be addressed until they are reproduced.
- If the team is able to reproduce the issue, it will be marked `needs-fix`, as well as possibly other tags (such as `critical`), and the issue will be left to be [implemented by someone](#your-first-code-contribution).


### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Graph, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the documentation carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/lexkrstn/go-graph/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
- Find out whether your idea fits with the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset. If you're just targeting a minority of users, consider writing an add-on/plugin library.

#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/lexkrstn/go-graph/issues).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- You may want to **include screenshots and animated GIFs** which help you demonstrate the steps or point out the part which the suggestion is related to. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
- **Explain why this enhancement would be useful** to most Graph users. You may also want to point out the other projects that solved it better and which could serve as inspiration.


### Your First Code Contribution

Thank you for your interest in contributing to the Graph library! This section will guide you through setting up your development environment and making your first contribution.

#### Prerequisites

Before you start contributing, make sure you have the following installed:

- **Go 1.19 or later** - Download from [golang.org](https://golang.org/dl/)
- **Git** - For version control
- **A code editor** - We recommend:
  - [GoLand](https://www.jetbrains.com/go/) (commercial, excellent Go support)
  - [VS Code](https://code.visualstudio.com/) with Go extension
  - [Vim/Neovim](https://neovim.io/) with Go plugins
  - Any editor you're comfortable with

#### Setting Up Your Development Environment

1. **Fork the repository**
   ```bash
   # Go to https://github.com/lexkrstn/go-graph and click "Fork"
   # Then clone your fork locally
   git clone https://github.com/YOUR_USERNAME/go-graph.git
   cd go-graph
   ```

2. **Set up the upstream remote**
   ```bash
   git remote add upstream https://github.com/lexkrstn/go-graph.git
   ```

3. **Verify your Go installation**
   ```bash
   go version
   # Should show Go 1.19 or later
   ```

4. **Install dependencies and run tests**
   ```bash
   go mod download
   go test ./...
   ```

5. **Run benchmarks to understand performance**
   ```bash
   go test -bench=. ./...
   ```

#### IDE Setup

**VS Code (Recommended for beginners):**
1. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
2. Install Go tools when prompted: `Ctrl+Shift+P` → "Go: Install/Update Tools"
3. Enable format on save: `Ctrl+Shift+P` → "Preferences: Open Settings (JSON)" and add:
   ```json
   {
     "go.formatTool": "goimports",
     "editor.formatOnSave": true,
     "[go]": {
       "editor.formatOnSave": true
     }
   }
   ```

**GoLand:**
1. Open the project
2. GoLand will automatically detect it as a Go module
3. Enable "Optimize imports on the fly" in Settings → Editor → General → Auto Import

#### Making Your First Contribution

1. **Choose an issue to work on**
   - Look for issues labeled `good first issue` or `help wanted`
   - Comment on the issue to let maintainers know you're working on it
   - If no suitable issues exist, check the TODO comments

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or for bug fixes:
   git checkout -b fix/issue-description
   ```

3. **Make your changes**
   - Follow the [style guide](#styleguides) below
   - Write tests for new functionality
   - Update documentation if needed

4. **Run tests and benchmarks**
   ```bash
   # Run all tests
   go test ./...
   
   # Run tests with coverage
   go test -cover ./...
   
   # Run benchmarks
   go test -bench=. ./...
   
   # Run race condition detection
   go test -race ./...
   ```

5. **Format and lint your code**
   ```bash
   # Format code
   go fmt ./...
   
   # Run linter (if you have golangci-lint installed)
   golangci-lint run
   ```

6. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new graph algorithm"
   # Follow the [conventional commit message guidelines](https://www.conventionalcommits.org/en/v1.0.0/).
   ```

7. **Push and create a pull request**
   ```bash
   git push origin feature/your-feature-name
   # Go to GitHub and create a pull request
   ```

#### Development Workflow

**Before starting work:**
```bash
# Always sync with upstream
git fetch upstream
git checkout main
git merge upstream/main
```

**During development:**
- Write tests first (TDD approach)
- Keep commits small and focused
- Use descriptive commit messages
- Test with different Go versions if possible

**Before submitting:**
- Ensure all tests pass
- Run benchmarks to check for performance regressions
- Update documentation if needed
- Squash commits if you have many small ones

#### Common Development Tasks

**Adding a new algorithm:**
1. Create a new file in the appropriate package
2. Add tests in a corresponding `_test.go` file
3. Add benchmarks to measure performance
4. Update documentation and examples

**Fixing a bug:**
1. Create a test that reproduces the bug
2. Fix the bug
3. Ensure the test passes
4. Add additional tests to prevent regression

**Improving performance:**
1. Use benchmarks to measure current performance
2. Make your improvements
3. Verify performance improvements with benchmarks
4. Document any trade-offs made

#### Getting Help

If you get stuck:
1. Check existing issues and pull requests
2. Ask questions in the issue you're working on
3. Join discussions in existing issues
4. Create a draft pull request early to get feedback

#### Code Review Process

1. **Self-review**: Review your own code before submitting
2. **Automated checks**: Ensure CI passes
3. **Maintainer review**: Address feedback from maintainers
4. **Merge**: Once approved, your code will be merged

Remember: Every contribution, no matter how small, is valuable! Don't hesitate to start with documentation improvements or simple bug fixes.
