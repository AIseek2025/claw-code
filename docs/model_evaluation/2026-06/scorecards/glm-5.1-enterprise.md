# Model Scorecard

## Basic Info

- Evaluation Window: `2026-06`
- Model: `glm-5.1 enterprise`
- Role In This Window: `benchmark candidate`
- Sample Count: `0 valid (2 attempted, B001 runtime_blocked both times — excluded)`
- Task Coverage: `B001 attempted x2 (runtime_blocked + model_mismatch, reviewed 2026-05-27)`

## Weighted Score

| Dimension | Weight | Raw Score (0-100) | Weighted Score | Notes |
| --- | ---: | ---: | ---: | --- |
| Delivery Success | 30 | TBD | TBD |  |
| Write-Back Reliability | 20 | TBD | TBD |  |
| Validation Closure | 15 | TBD | TBD |  |
| Blocker Convergence | 15 | TBD | TBD |  |
| First-Iteration Effectiveness | 10 | TBD | TBD |  |
| Long-Run Stability | 5 | TBD | TBD |  |
| Cost Efficiency | 5 | TBD | TBD |  |
| **Total** | **100** |  | TBD |  |

## Run Log

| # | Task | Classification | Score Impact | Notes |
|---|------|---------------|-------------|-------|
| 1 | B001 Rust Bugfix (attempt 1) | `runtime_blocked` (no validated provider route) | excluded | Configured Z.AI enterprise id returned 400; direct bigmodel route returned 401; proxy route returned 502 |
| 2 | B001 Rust Bugfix (attempt 2) | `runtime_blocked` + `model_mismatch` | excluded | Enterprise route still HTTP 400; OpenClaw silently fell through to `xfyun/astron-code-latest` which completed the task successfully — but this is NOT glm-5.1 enterprise data |

## Coverage Focus

- Core tasks to prioritize:
  - `B001`
  - `B005`
  - `B007`
  - `B008`
- High-value unknowns to validate:
  - 是否能摆脱长链路 write-back 脆弱性
  - `App`
  - `高并发架构`
  - `运维 / 故障恢复`

## Observations (from batch_01 B001)

### B001 Run Status (attempt 1 — first review)
- **Classification**: `runtime_provider_error` — current enterprise lane has no validated provider route
- **Sample validity**: Run is **not countable** — model never received the benchmark task
- **Workspace state**: Correctly injected (auth.rs:139 has unreachable `ExternalUser` branch); workspace is intact and ready for rerun
- **Runner log status**: latest rerun now captures provider 400 Unknown Model details after shell logger fix
- **Model capability signal**: `none` — this sample contains configuration/runtime information only

### B001 Run Status (attempt 2 — second review, 2026-05-27)
- **Classification**: `runtime_provider_error` + `model_mismatch`
- **Run ID**: `batch01-B001-glm-ent`
- **Nominal model**: `glm-5.1 enterprise` (configured)
- **Actual executing model**: `xfyun/astron-code-latest` (OpenClaw fallback default)
- **Enterprise route**: still HTTP 400 "Unknown Model" for `zai/glm-5.1-enterprise` — unchanged from attempt 1
- **Task outcome**: Completed successfully by fallback model — fix correct, tests pass, regression test added
- **Independently verified**: cargo test 30/30 pass, fmt check pass, clippy 0 new warnings
- **Capability data for glm-5.1 enterprise**: ZERO — cannot attribute this work to glm-5.1 enterprise
- **Orchestration concern**: OpenClaw silently fell through to default model without surfacing model mismatch in run metadata; first-pass human review scores were incorrectly attributed to glm-5.1 enterprise

### What we know about glm-5.1 enterprise from batch_01
1. **Zero capability data** — across 2 attempts, no enterprise route has reached task dispatch
2. **Infrastructure**: Workspace binding and run_id selection were correct
3. **Thinking mode compatibility**: Runner now uses `thinking="on"` for glm-family runs
4. **Configured route failure**: `zai/glm-5.1-enterprise` fails with 400 on both attempts
5. **Alternate route failure**: direct `open.bigmodel.cn/api/coding/paas/v4` and `open.bigmodel.cn/api/paas/v4` probes both return 401; `proxy_glm_51` returns 502 upstream failure
6. **Route-convention drift**: `llm_wiki` preset says glm-5 / glm-5.1 should route via `api.z.ai`, but runtime registry still points BIGMODEL_GLM_51_API_KEY to `open.bigmodel.cn`
7. **Silent fallback risk**: OpenClaw fallback to default model without mismatch detection can produce misleading "success" signals

### What we still don't know (all dimensions remain TBD)
- Rust bugfix capability (B001 rerun with actual enterprise route needed)
- Scope discipline
- Diagnostic methodology
- Write-back reliability
- Toolchain competence
- All other benchmark dimensions

## Current Notes

- 本文件用于剥离 `2026-05` 自然样本中的 runtime/orchestrator 干扰。
- `2026-06` 的目标是看 enterprise 版本在统一 benchmark 中能否稳定转化成真实交付。
- 两次 B001 尝试均被 enterprise provider 路由问题阻断，未产生任何 glm-5.1 enterprise 能力数据。
- 第二次尝试（batch01-B001-glm-ent）虽然任务被执行完成，但实际执行模型是 `xfyun/astron-code-latest`，不能计入 glm-5.1 enterprise 样本。
- 当前没有一条已验证可用的 enterprise 路线：配置路由 400、bigmodel 候选 401、proxy 候选 502。
- 仓库内部还存在 provider 约定漂移：UI preset 认为 glm-5.1 应走 `api.z.ai`，但运行时 registry 仍把 BIGMODEL 线当作 glm-5.1 候选。
- 继续推进前需要先恢复一条真实可用的 enterprise provider path，再设置 `OPENCLAW_MODEL_GLM_ENTERPRISE` / OpenClaw provider 映射。
- 建议 OpenClaw 增加 model mismatch detection，避免 fallback 静默替换导致评估信号污染。
