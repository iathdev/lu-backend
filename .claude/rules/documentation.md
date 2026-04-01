# Documentation

All design documents, requirement specs, API contracts, and technical documentation MUST be placed in `.claude/docs/`, NOT inside `internal/<module>/docs/`.

## Directory Structure

```
.claude/docs/
└── 01-backend/
    └── phase-1-vocabulary/
        ├── auth/
        ├── ocr/
        └── vocabulary/
```

## Rules

- New module documentation goes under the appropriate phase and module directory in `.claude/docs/`.
- Do NOT create `docs/` directories inside `internal/<module>/`.
- Domain model diagrams, API contracts, requirements, and technical challenges all belong in `.claude/docs/`.
