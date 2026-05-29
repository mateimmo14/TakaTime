# Installation & Development Setup

## Prerequisites

| Component | Version             |
| --------- | ------------------- |
| Go        | >= 1.25.3           |
| MongoDB   | >= 6.0 recommended  |
| Node.js   | >= 18 (recommended) |
| VS Code   | >= 1.80.0           |
| JDK       | 17 or 21            |
| Neovim    | >= 0.9 recommended  |

---

# Core Go Backend

## Install Dependencies

```bash
go mod download
```

## Build

```bash
go build -trimpath -ldflags="-s -w"
```

## Run

```bash
go run .
```

---

# VS Code Extension

## Install Dependencies

```bash
npm install
```

## Lint

```bash
npm run lint
```

## Run Tests

```bash
npm test
```

## Launch Extension Development Host

```bash
code .
```

Press `F5` inside VS Code to launch the Extension Development Host.

---

# Neovim Plugin

## Requirements

* Neovim >= 0.9 recommended

## Example Setup

```lua
require("taka-time").setup({
  debounce_seconds = 2,
  debug = false,
})
```

## Initial Configuration

```vim
:TakaInit
```

Run `:TakaInit` inside Neovim to configure the MongoDB connection URI.

## Available Commands

```vim
:TakaInit
:TakaDash
:TakaStatus
:TakaIgnore
:TakaTrack
```


---

# JetBrains Plugin

## Requirements

* JDK 17 or JDK 21
* Gradle

## Run IDE Sandbox

```bash
./gradlew runIde
```

## Build Plugin

```bash
./gradlew buildPlugin
```
