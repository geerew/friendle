# OffCourse — Agent Guide

This file is the primary reference for AI assistants working in this repository. Follow it for Go style, project layout, and how to build and test locally. When the project changes in ways this guide should reflect—new packages, tools, workflows, or conventions—update **AGENTS.md** as part of the same work, where applicable.

## Project purpose

**OffCourse** is a local course management application. It scans course directories on the filesystem, organizes assets (videos, PDFs, markdown, text) into lessons, and exposes a web UI for browsing courses and tracking progress—without requiring an internet connection.

- **Backend**: Go REST API, SQLite (`data.db`, `logs.db`), CGO (SQLite + libwebp for course cards), FFmpeg and FFprobe (on-demand HLS video transcoding)
- **Frontend**: SvelteKit + TypeScript in `ui/`, embedded into the Go binary at build time
- **Data**: Default `oc_data/` (override with `--data-dir`); holds DBs, HLS transcodes, optimized card WebP files

See [README.md](README.md) for architecture, CLI commands, bootstrapping, and Docker.

## Project layout

| Path | Role |
|------|------|
| `main.go` | Binary entry; delegates to `cmd` |
| `cmd/` | CLI (Cobra): `serve`, `admin`, etc. |
| `api/` | HTTP handlers (Fiber), route wiring per resource |
| `dao/` | Database access |
| `models/` | Domain types |
| `database/` | DB setup and migrations wiring |
| `migrations/` | SQL migrations |
| `cron/` | Scheduled jobs (e.g. card cache warm) |
| `utils/` | Shared packages (`cardcache`, `media`, `logger`, …) |
| `ui/` | SvelteKit frontend |
| `app/` | App-level wiring (`Config`, bootstrap state, log components) |
| `docker/` | Container build and compose docs |

Packages under `utils/` (and similar modules) should use **`base.go`** as the module’s starting file and expose **`New`** as the constructor—not names like `NewCardCache`.

## Tools and prerequisites

- **Go** ≥ 1.22.4, **CGO** enabled (SQLite + libwebp)
- **Node.js** ≥ 22.12, **pnpm** ≥ 8 (frontend)
- **Make** — `make build` produces `offcourse` (build UI first; `ui/build` is embedded)
- **[air](https://github.com/air-verse/air)** — backend live reload in dev (`.air.toml`)
- **FFmpeg / FFprobe** — video HLS only (not required for card-only work)
- Native libs: SQLite dev headers, libwebp dev headers (see README Prerequisites)

## Development (two terminals)

Run backend and frontend separately; both reload on change. Traffic goes through the Go app—ignore the Vite URL for browsing.

**Terminal 1 — frontend**

```bash
cd ui
pnpm install   # first time
pnpm run dev
```

**Terminal 2 — backend**

```bash
go mod download   # first time
air               # optional: air -- --http 0.0.0.0:8080
```

Open the URL from `Bootstrap required: ...` (first run) or `Server started at ...` (already bootstrapped).

## Testing

From the repository root:

```bash
go test -tags dev ./...
```

- **`-tags dev`** — stub UI embed so tests do not require `ui/build` first
- **CGO** — same SQLite/libwebp requirements as a normal build
- Verbose: add `-v`

If `ui/build` already exists from a frontend build, `go test ./...` without `-tags dev` also works.

When adding or changing Go code, add tests **where they make sense**. Do not chase 100% coverage unless that is easily attainable. See [§9 Tests](#9-tests) below.

## AI-generated code policy

AI-assisted changes are welcome. Before anything is committed, a **human must review and understand** the diff—what changed, why, and that it fits the codebase. Refactor or clean up as needed; do not commit code you have not read and could not explain.

**Keep this guide current.** If your changes alter how the repo is structured, built, tested, or how Go code should be written, update **AGENTS.md** in the same change (or call it out for the human to do so). Do not let the guide drift from reality.

---

## Go coding standards

These rules apply to all `.go` files in this project.

### 1. Section separators

Separate logical sections with a blank line and this comment on its own line:

```go
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
```

Use it between types, vars, constructors, and major function groups—see
`app/base.go` for the canonical pattern.

**Always separate with a divider, even when items are closely related:**

- Each **type** definition from the next (including consecutive request/response
  structs in `api/types.go`)
- A **type** and its **helper function** (e.g. `userResponse` and
  `userResponseHelper`)—the helper gets its own section after the type
- Each **function** from the next, including helpers that share a purpose (e.g.
  `userGroupSummaryResponsesFromRows` and `userGroupSummariesByUserID`)

**Example:**

```go
// userResponse is a user returned by the API
type userResponse struct {
	// ...
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userResponseHelper maps users to API responses
func userResponseHelper(users []*models.User) []*userResponse {
	// ...
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminUserResponse is a user row enriched for the site admin list
type adminUserResponse struct {
	// ...
}
```

### 2. Spacious control flow

Prefer a blank line before a trailing `return` that is not part of the same block as the preceding logic.

Also put a blank line between **consecutive `if` statements** at the same block level—each
`if` is its own step, not a run-on chain.

**Exceptions** (no blank line):

- Assignment and its `if err != nil` check—see [§3](#3-error-handling-stays-tight)
- `else` / `else if` attached to the preceding `if`
- An `if` nested **inside** another `if`'s body (only **sibling** `if` blocks need spacing)

**Avoid:**

```go
if a == b {
	return this
}
return that
```

```go
if req.Name != nil {
	g.Name = strings.TrimSpace(*req.Name)
}
```

**Prefer:**

```go
if a == b {
	return this
}

return that
```

```go
if req.Name != nil {
	g.Name = strings.TrimSpace(*req.Name)
}
```

### 3. Error handling stays tight

Keep the assignment and its error check together—no blank line between them.

```go
a, err := doSomething()
if err != nil {
	return xxx
}
```

### 4. Function order

Within a file:

1. **Public** functions and methods at the **top**
2. **Private** (unexported) functions at the **bottom**

### 5. Functions and visibility

**Size and helpers**

- Avoid tiny functions (a few lines) that are only called **once**—inline them when readability is just as good
- A small shared function used in **multiple** places is fine
- Judge length by **how much work** the function does, not raw line count. Many lines of error handling or multi-line log calls can still be a small, focused function
- Extract helpers when they take on a **meaningful chunk** of work—not one-liner wrappers used in a single caller

**Signatures**

Keep the **function definition on one line**—`func` name, parameters, and return
types together. Do not break parameters or returns across lines for length; a long
signature is fine.

If a function’s **body** grows hard to follow, shorten the body (extract helpers,
simplify logic)—not by splitting the signature.

**Avoid:**

```go
func groupSearchResponsesFromRows(
	rows []*models.GroupSearchRow,
	memberGroupIDs map[string]struct{},
	pendingGroupIDs map[string]struct{},
) []*groupSearchResponse {
```

**Prefer:**

```go
func groupSearchResponsesFromRows(rows []*models.GroupSearchRow, memberGroupIDs map[string]struct{}, pendingGroupIDs map[string]struct{}) []*groupSearchResponse {
```

**One-line function bodies**

Do not put the entire function on one line—definition and body together—even for trivial
methods. Give each function its own multi-line body and separate consecutive functions with
a section divider (§1).

**Avoid:**

```go
func (a *App) Close() error { return nil }

func (a *App) IsBootstrapped() bool { return a.bootstrapped.Load() == 1 }
func (a *App) SetBootstrapped()     { a.bootstrapped.Store(1) }
```

**Prefer:**

```go
// Close releases application resources
func (a *App) Close() error {
	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsBootstrapped reports whether the application has completed first-run setup
func (a *App) IsBootstrapped() bool {
	return a.bootstrapped.Load() == 1
}
```

**Visibility**

- Default to **unexported** functions, methods, and struct fields unless something outside the package genuinely needs them
- Do not export for convenience or “just in case”

**Globals**

- Avoid package-level app singletons and mutable globals when dependency injection or passing values through `New` / constructors is reasonable

### 6. Module entrypoints

For a cohesive package (especially under `utils/`):

- **`base.go`** — types, config, and `New` (or primary setup)
- Constructor name: **`New`**, not `NewCardCache`, `NewFooService`, etc.
- Split other concerns into additional files in the same package (`serve.go`, `http.go`, …)

### 7. Comments

Every function and method needs a comment, including unexported ones. Keep comments **short and concise**. Start with the function name (Go doc convention). Wrap comment lines around **85–90** columns—do not let comments run too long on one line.

**Punctuation**

- A **single sentence** (one comment line, or the final sentence of a block) does **not** end with a full stop
- When **two or more sentences** appear on the **same** comment line, put a full stop **between** those sentences—not after the last one

```go
// validateCourseID rejects path-unsafe course IDs
func validateCourseID(id string) error {

// warmServeIndex loads card serve metadata from the database and disk. It is called once when
// the cron scheduler starts
func warmServeIndex(ctx context.Context) error {
```

**Describe the present**

Comments explain how code works **now**. Do not document history, migrations, or “we used to X but now Y”—especially when changing behaviour. Remove or rewrite stale comments; do not leave changelog-style notes in source.

**Avoid:**

```go
// Previously we returned nil here; now we return an error when the path is missing
```

**Prefer:** update the comment to match current behaviour, or delete it if the code is clear without it.

**Separate blocks**

Related detail can stay on one line or consecutive `//` lines. For several **standalone** important points, split them and put a blank commented line (`//`) between blocks.

**Avoid:**

```go
// Some comment about the function
// extra note: this panics on nil input
func xxx(xxx) {
```

**Prefer:**

```go
// Some comment about the function
//
// extra note: this panics on nil input
func xxx(xxx) {
```

### 8. Review before commit

AI may implement or edit Go code in this repo. The author who commits is responsible for **reading and understanding** those changes first—not shipping blind copy-paste. See [AI-generated code policy](#ai-generated-code-policy) above.

### 9. Tests

Add tests when they add real value—behaviour worth guarding, non-trivial branches, regressions. Skip trivial or redundant tests. Do not target 100% coverage unless it falls out naturally.

Use **`github.com/stretchr/testify/require`** only—always **`require.xxx`** (fail fast), not `assert`.

**One test function per production function.** Group scenarios with `t.Run`. Do not split the same function across multiple top-level `TestXxx` / `TestYyy` helpers.

**`t.Run` names** — keep them short and specific (`nil config`, `empty string`, `stat error`), not long prose.

**Test comments** — short line above each scenario: `// Test successfully …` or `// Test error due to …`. Same for a test function with **no** `t.Run`: one comment above the function body.

```go
func TestSomething(t *testing.T) {
	// Test successfully processing all courses in a batch
	t.Run("update", func(t *testing.T) {
		// ...
	})

	// Test error due to filesystem stat failure
	t.Run("stat error", func(t *testing.T) {
		// ...
	})
}

// Test successfully formatting an ETag
func TestFormatETag(t *testing.T) {
	require.Equal(t, `"abc"`, FormatETag("abc"))
}
```

See `dao/user_test.go` and `utils/types/date_time_test.go`.

**Table-driven tests** when exercising many inputs. Prefer separate tables (or `t.Run` blocks) for **success** and **error** paths—not one mixed table unless cases truly share the same setup.

```go
func TestNewSiteRole(t *testing.T) {
	// Test successfully parsing known site role strings
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			input    string
			expected SiteRole
		}{
			{"site_admin", SiteRoleAdmin},
			{"admin", SiteRoleAdmin},
			{"site_user", SiteRoleUser},
		}

		for _, tt := range tests {
			require.Equal(t, tt.expected, NewSiteRole(tt.input))
		}
	})

	// Test unknown role defaults to site user
	t.Run("default", func(t *testing.T) {
		require.Equal(t, SiteRoleUser, NewSiteRole("unknown"))
	})
}
```
