Project brief — MCF TUI (for pc-style/MCF)

Context
MCF (My Claude Flow ) is a sophisticated development automation platform built around Claude Code with Serena semantic analysis integration (see repository: pc-style/MCF). The goal is to design and implement a high-quality, production-ready Terminal UI (TUI) that exposes MCF’s operational surface, accelerates developer workflows, and leverages Serena and MCF’s existing agents, hooks, and commands.  [oai_citation:1‡GitHub](https://github.com/pc-style/MCF)

Primary goals
1. Fast, keyboard-first, Vim-like navigation for everyday MCF operator tasks.
2. Deep integration with MCF features: Claude Code commands, Serena semantic analysis, the custom command system, and event hooks.
3. A polished UX with fuzzy search, a command bar, dialogs, compositing (notifications/popups), and optional mouse support.
4. Distributable as a single Go binary (no runtime deps).

Tech stack & libraries
- Language: Go (single binary distribution)
- TUI framework: BubbleTea (charmbracelet/bubbletea)
- Theming: bubbletint
- Components: bubbles + selective 3rd-party components (bubble-table, teacup, bubbleup, stickers) as needed
- Persistent storage: use .claude/settings.json for settings
- project managment

- Integration: Use existing MCF CLI/commands and Serena endpoints (/serena:* and .claude commands) exposed by the repository for semantic operations.  [oai_citation:2‡GitHub](https://github.com/pc-style/MCF)

High-level feature list (user-facing)
1. Dashboard / Cluster Overview
   - At-a-glance MCF health: version, Serena status.
   - Quick actions: run health checks, trigger /serena:status.
   - Compact, colorized status line + summary widgets.

3. System logs (streaming)
   - On-demand tailing of MCF logs (searchable, follow/stop).
   - Filter logs by subsystem (serena, claude, hooks).



MCP server manager



Command bar & fuzzy search
    - Command palette (colon commands like `:secrets`, `:policies`, `:systemlogs`).
    - Integrated fuzzy search (`/` to search current view) with live filtering; configurable threshold to disable fuzzy for very large lists.

Navigation model
    - Root pages, sub-pages, breadcrumbs.
    - Escape to pop navigation stack.
    - Each page can expose page-scoped commands.
    - Optional tree navigator with keyboard shortcuts.

Security & permissions
- The TUI must respect MCF auth permissions. If LIST is forbidden but READ permitted, provide a path input dialog so users can fetch known paths.
- Sensitive data masking by default; explicit user action required to reveal secrets.
- Audit logging of tool actions (optionally to the configured audit device).
- Command whitelisting for any shell plugin execution.

Integration points (developer-facing)
- Use existing MCF command surface from `.claude/commands/` to implement backend actions (e.g., call `/serena:find` or `gh:push` equivalents via MCF’s command runner).
- Serena operations: implement dedicated Serena command adapters for symbol find/analyze/refs/status to power semantic navigation and editing workflows.  [oai_citation:3‡GitHub](https://github.com/pc-style/MCF)
- Hooks: surface hook-driven suggestions in the TUI (e.g., when viewing a function, prompt `/serena:analyze` via a non-blocking suggestion).

Non-functional requirements
- Fast: <200ms interaction latency for local operations on typical repos.
- Small binary footprint: compile with Go 1.20+; avoid heavy dependencies.
- Cross-platform: Linux, macOS, Windows (WSL).
- Configurable: theming, keybindings, fuzzy-search behavior, persistence location (XDG).
- Testable: unit tests for core view logic and integration tests that can run against an instrumented MCF dev instance.


Notes & references
- Use existing MCF layout and commands (.claude/, .serena/) to avoid duplicating logic; adapter layer should translate TUI actions into MCF commands and Serena requests.  [oai_citation:4‡GitHub](https://github.com/pc-style/MCF)
- Look to k9s for UX inspiration (resource context plugins) and charmbracelet examples (bubbletea examples).

If you want, next I can:
- produce terminal wireframes for the main pages (dashboard, secrets, policy editor),
- scaffold a Go module with BubbleTea main loop and example page,
- or write the adapter interface that maps TUI actions to MCF commands (serena, gh, hooks).

