# Espresso

A macOS menu bar app that wraps `caffeinate` to prevent the Mac from sleeping.
Written in Go. No dock icon (`LSUIElement = true`).

## Project structure

| File | Purpose |
|---|---|
| `caffeinate.go` | Process manager — spawns/kills `caffeinate`, signals expiry |
| `caffeinate_test.go` | Unit tests for the process manager |
| `tray.go` | Menu bar UI, animation ticker, event loop |
| `login.go` + `login.m` | CGO bridge to SMAppService for Login Item registration |
| `main.go` | Entry point — wires Manager → systray, registers login item |
| `Info.plist` | App bundle metadata |
| `Makefile` | Build targets |
| `docs/index.html` | Landing page template — contains `{{VERSION}}` placeholder |
| `.github/workflows/release.yml` | Release pipeline |

## Build

```bash
make build          # produces Espresso.app in the project root
make clean          # removes binary and Espresso.app
```

Requires macOS and Go 1.21+. The CGO login item code links against
ServiceManagement and Foundation frameworks — no extra setup needed.

## Run tests

```bash
go test ./... -v
```

## Release workflow

To ship a new version:

```bash
git tag v1.x.x
git push --tags
```

GitHub Actions takes it from there:
1. Builds `Espresso.app` on a macOS runner
2. Packages `Espresso.dmg`
3. Injects the version into the landing page and deploys it to GitHub Pages
4. Creates a GitHub Release and attaches `Espresso.dmg`

**Important:** `docs/index.html` contains a `{{VERSION}}` placeholder that
Actions replaces at deploy time. Do not manually edit the version string in
that file — it will be overwritten on the next release.

## GitHub Pages

The landing page is served from GitHub Pages.
Source must be set to **"GitHub Actions"** in repo Settings → Pages.
