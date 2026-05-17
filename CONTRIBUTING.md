# Contributing to TakaTime

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them.

## Development Setup

### Prerequisites
1.  **Node.js** (v16+) and **npm**
2.  **Go** (v1.20+) for the binary
3.  **VS Code**

### Setup Steps
1.  Fork the repository.
2.  Clone your fork: `git clone https://github.com/YOUR-USERNAME/TakaTime.git`
3.  Install dependencies:
    ```bash
    npm install
    ```
4.  Build the Go binary (if modifying the uploader):
    ```bash
    cd binary
    go build -o taka-uploader
    ```

## Reporting Bugs

Bugs are tracked as [GitHub issues](https://github.com/Rtarun3606k/TakaTime/issues).
Explain the problem and include additional details to help maintainers reproduce the problem:
* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible.
* **Provide specific examples to demonstrate the steps**.

## Pull Requests

1.  Fill in the required template
2.  Do not include issue numbers in the PR title
3.  Include screenshots and animated GIFs in your pull request whenever possible.

**WARNING: If this PR is just fixing a typo, adding your name to a list, or a trivial README change, DO NOT submit it. It will be marked as spam and closed.**
- [ ] I confirm this is a substantial contribution, not a trivial spam PR.