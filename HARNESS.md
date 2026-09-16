# Xanders harness

This project adopts a focused Codex/SAW harness based on
[`bybren-llc/safe-agentic-workflow`](https://github.com/bybren-llc/safe-agentic-workflow).

Included components:

- `.codex/`: Codex configuration and SAFe role profiles.
- `.agents/skills/`: delivery, API, frontend, migration, testing, security,
  discovery, and release skills.
- `AGENTS.md`: project-specific architecture, credit policy, and quality gates.
- `.harness-manifest.yml`: upstream-sync metadata.

The upstream project is MIT licensed. This adoption is intentionally limited to
the capabilities used by Xanders; it does not configure Linear, Confluence,
Stripe, or remote autonomous agent infrastructure.
