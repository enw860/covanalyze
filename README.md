# covanalyze

## Project Summary

A deterministic, offline, local Go command-line tool that uses `coverage.out` as the single source of truth for coverage and optionally enriches that data with AST-derived semantics from Go source files. The tool should produce a structured JSON report that remains valid with `coverage.out` alone, while adding human-readable semantic detail such as function names, uncovered branch types, loop and return hints, error-path hints, and best-effort conditions when source code is available.