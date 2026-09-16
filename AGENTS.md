# Xanders AI Video Platform — Agent Harness

This project uses a tailored subset of the SAFe Agentic Workflow (SAW) harness.
It is a workflow and quality system; it does not replace application security
controls or human approval for product, billing, or deployment decisions.

## Project map

- `backend/`: Go, Fiber v2, GORM, PostgreSQL, JWT authentication.
- `web/`: Next.js 16, React 19, TypeScript, TanStack Query.
- `backend/migrations/`: ordered Goose SQL migrations.
- `TASKS.md`: project implementation backlog and context.

## Standard delivery loop

1. Define acceptance criteria before implementation.
2. Inspect existing patterns before introducing a new one.
3. Implement the smallest coherent change.
4. Run focused validation, then the affected backend and frontend checks.
5. Document evidence and unresolved risk in the handoff.

## Security and credit invariants

- Authentication and authorization are verified by backend handlers; UI state is never authoritative.
- A new account receives one 15-credit signup allowance. Login never replenishes credits.
- A video costs exactly 5 credits unless the product owner changes the backend constant and matching UI copy together.
- Credit deduction and job creation remain one database transaction with a row lock.
- A user with insufficient credits receives a 402 response and is directed to `alfanfaturahman10@gmail.com`.
- Do not create a public endpoint that adds credits. Any manual grant must be an audited administrative action.
- Preserve existing balances unless an explicit product decision authorizes a migration to change them.

## Validation commands

```powershell
cd backend; go test ./...; go vet ./...; go build ./...
cd web; npm run lint; npm run build
```

## Quality gates

- API changes: verify auth, ownership isolation, validation, and error responses.
- Database changes: include a forward and rollback migration where feasible; do not alter existing migration history alone.
- UI changes: verify loading, success, insufficient-credit, and unauthenticated states.
- Before deployment: run the `security-audit` skill/checklist and never expose secrets in code or client bundles.

## Harness skills installed

Use the relevant project skill for spec creation, safe workflow, API work,
frontend work, migrations, testing, security audits, pattern discovery, and
release readiness. The skills are located in `.agents/skills/`.
