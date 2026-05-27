# Code Excerpts - Iteration 6 Evidence

## 1. Local CI Gate Script Header + Job Logic (scripts/local-ci/openapi-contract.sh)

```bash
# scripts/local-ci/openapi-contract.sh — local CI gate for yunmao shared contract.
#
# Purpose:
#   Runs the same 4-job pipeline as .github/workflows/openapi-contract.yml
#   but locally, on this machine (the "project-required environment").
#
# Exit 0 = gate PASS (all 4 jobs green).
# Exit 1 = gate FAIL (any job red).
# Output directory: reports/local-ci-runs/<timestamp>/

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
YUNMAO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RUN_ID="$(date +%Y%m%dT%H%M%S)"
RUN_DIR="$YUNMAO_ROOT/reports/local-ci-runs/$RUN_ID"

mkdir -p "$RUN_DIR"
cd "$YUNMAO_ROOT"

LOG="$RUN_DIR/run.log"
JOB_SUMMARY="$RUN_DIR/jobs.json"
# ... (full script at ./local-ci-scripts/openapi-contract.sh, 4898 bytes)
```

The script runs 4 jobs matching `.github/workflows/openapi-contract.yml`:
- **Job 1 (spec-lint)**: `pushd go/pkg/yunmao && go test ./openapi/... -v -count=1`
- **Job 2 (gen-typescript-web)**: `pushd clients/web && npm run openapi-gen && npx tsc --noEmit && npm run test:run`
- **Job 3 (gen-typescript-admin)**: same for `clients/admin`
- **Job 4 (contract-consistency)**: SHA256 diff before/after re-gen to detect drift

## 2. Fresh Run Output (20260526T224212, captured iter 6)

```
[22:42:12] === yunmao local-ci openapi-contract ===
[22:42:12] Run ID: 20260526T224212
[22:42:12] Repo:   /Users/brando/Documents/trae_projects/CodeMaster/isolated_autoruns/yunmao
[22:42:12] User:   brando
[22:42:12] Host:   brandos-MacBook-Pro-2.local
[22:42:12]
[22:42:12] = START JOB: spec-lint
[22:42:12] = END JOB: spec-lint (PASS)
[22:42:12] = START JOB: gen-typescript-web
[22:42:14] = END JOB: gen-typescript-web (PASS)      Duration  848ms
[22:42:14] = START JOB: gen-typescript-admin
[22:42:17] = END JOB: gen-typescript-admin (PASS)    Duration  482ms
[22:42:17] = START JOB: contract-consistency
[22:42:18] = END JOB: contract-consistency (PASS)
[22:42:18]
[22:42:18] === SUMMARY ===
[22:42:18] Passed: 4  (spec-lint gen-typescript-web gen-typescript-admin contract-consistency)
[22:42:18] Failed: 0  (none)
[22:42:18] === OVERALL: PASS ===
```

(Full log at `./run-20260526T224212/run.log`)

## 3. jobs.json (20260526T224212, captured iter 6)

```json
{
  "run_id": "20260526T224212",
  "timestamp": "2026-05-27T05:42:18Z",
  "user": "brando",
  "host": "brandos-MacBook-Pro-2.local",
  "repo": "/Users/brando/Documents/trae_projects/CodeMaster/isolated_autoruns/yunmao",
  "jobs_passed": ["spec-lint","gen-typescript-web","gen-typescript-admin","contract-consistency"],
  "jobs_failed": [],
  "overall": "PASS",
  "log_dir": "/Users/brando/Documents/trae_projects/CodeMaster/isolated_autoruns/yunmao/reports/local-ci-runs/20260526T224212"
}
```

(Full JSON at `./run-20260526T224212/jobs.json`)

## 4. Crontab Registration (captured iter 6)

```
=== crontab -l (captured at 2026-05-27T05:42:25Z) ===
0 * * * * /Users/brando/Documents/trae_projects/CodeMaster/isolated_autoruns/yunmao/scripts/local-ci/openapi-contract.sh >> /Users/brando/Documents/trae_projects/CodeMaster/isolated_autoruns/yunmao/reports/local-ci-runs/cron.log 2>&1
```

(Full content at `./host-state-crontab.log`)

## 5. spec-lint.log excerpt (idempotent across 3 runs - same SHA256: 77c0052b2d7b5df8...)

```
=== RUN   TestSpecIsValidJSON
--- PASS: TestSpecIsValidJSON (0.00s)
=== RUN   TestSpecHasRequiredFields
--- PASS: TestSpecHasRequiredFields (0.00s)
=== RUN   TestSpecHasSchemas
--- PASS: TestSpecHasSchemas (0.00s)
=== RUN   TestSpecPathsCoverMainServices
--- PASS: TestSpecPathsCoverMainServices (0.00s)
=== RUN   TestSpecOperationsHaveSchemas
--- PASS: TestSpecOperationsHaveSchemas (0.00s)
PASS
ok      yunmao.live/pkg/yunmao/openapi  0.394s
```

## 6. contract-consistency.log (idempotent across 3 runs - same SHA256: b7e186f...)

```
Web:   pre=857368d5c88e7510... post=857368d5c88e7510...  MATCH
Admin: pre=857368d5c88e7510... post=857368d5c88e7510...  MATCH
```

(SHA256 of both clients' generated-api.ts matches pre- and post-regeneration)

## 7. Historical Run Summary

| Run ID | Timestamp UTC | Exit | Jobs Passed | Evidence Dir |
| --- | --- | --- | --- | --- |
| 20260526T223702 | 2026-05-27T05:37:11Z | 0 | 4/4 | `./run-20260526T223702/` |
| 20260526T223916 | 2026-05-27T05:39:22Z | 0 | 4/4 | `./run-20260526T223916/` |
| 20260526T224212 | 2026-05-27T05:42:18Z | 0 | 4/4 | `./run-20260526T224212/` |

All 12 job runs (3 gate invocations × 4 jobs each) returned PASS.

## 8. Artifact Inventory (24 files, SHA256 checksummed)

All files listed in `./artifact_inventory.txt`. Auditor may verify any hash via:
```
cd reports/iteration_6_evidence && sha256sum -c artifact_inventory.txt
```
