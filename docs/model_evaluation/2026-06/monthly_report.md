# Monthly Model Comparison Report

## Window

- Month: `2026-06`
- Task Pool: `benchmark_task_pool.md`
- Projects: `benchmark workspaces to be created per task`
- Snapshot Rule: `same repo snapshot, same handoff, same gate, same timeout/retry`

## Compared Models

1. `qwen3.7-max`
2. `glm-5.1 enterprise`
3. `glm-5.1`

## Executive Summary

- Primary Coder Recommendation: `qwen3.7-max (provisional — sole B001 sample, 86.6 adjusted score)`
- First Backup Recommendation: `TBD (glm-5.1 / glm-5.1 enterprise runs pending)`
- Main Reason: `batch_01 B001 completed for qwen3.7-max — Rust bugfix success with evidence gap; glm models still pending valid runs`
- Confidence Level: `low (single sample, evidence gap on diff)`

## Batch Execution Log

### batch_01: B001 Rust Service Bugfix

- Runs: `qwen3.7-max` (completed & reviewed); `glm-5.1 enterprise` (runtime_blocked x2 & reviewed — second attempt fell through to xfyun/astron-code-latest); `glm-5.1` (runtime_blocked & reviewed)
- Outcome: `qwen3.7-max` 已形成 1 个有效成功样本；两条 `glm` 线都被 runtime/provider blocker 截断，暂未形成可计分样本
- Current State: 三套 workspace 已在统一已知失败态上对齐过；`qwen3.7-max` 完成修复并通过全量验证；`glm-5.1 enterprise` 两次尝试均被 enterprise 路由 400 阻断（第二次 fallback 到 xfyun/astron-code-latest 完成任务但不可归因）；`glm-5.1` 因 provider 429 quota 阻断
- Classification: `qwen3.7-max = success (with evidence gap)`；`glm-5.1 enterprise = runtime_provider_error + model_mismatch (enterprise 400, fallback non-attributable)`；`glm-5.1 = runtime_provider_error / provider_quota_exhausted`
- qwen3.7-max behavior: 正向能力信号成立，但 `diff.patch` 仍是旧阶段产物，结论带 evidence gap
- Action Required: 修复 `glm-5.1 enterprise` model code 注册、恢复 `glm-5.1` provider quota，并补强 runner logging 后重跑两条 glm 样本

| Model | B001 Status | Files Changed | Tests Pass | Classification | Reviewed |
|-------|------------|---------------|------------|----------------|----------|
| qwen3.7-max | completed | 1 | 740/740 | success (with evidence gap) | yes — second review verified updated evidence chain |
| glm-5.1 enterprise | runtime_blocked x2 | 1 (by fallback model) | 30/30 (by xfyun/astron-code-latest) | runtime_provider_error + model_mismatch (enterprise 400, fallback silently executed task) | yes — second review confirmed fix correctness but data not attributable to glm-5.1 enterprise |
| glm-5.1 | runtime_blocked | 0 | not reached | runtime_provider_error (Z.AI 429) | yes — confirmed provider blocker, no model signal |

#### batch_01 Review Summary (2026-05-27)

Reviewer: `hermes-agent (qwen3.7-max)`

**Evidence verification:**
- pre_fix_failure_evidence.txt: 2 个目标失败 (`detects_bearer`, `extracts_full_context`) 与注入 bug 逻辑一致
- test_output_full.txt: 740 passed / 0 failed — independently verified by summing 56 test-result lines across 15 crates
- build_verification.txt: exit 0 on cargo test / fmt / clippy
- clippy_summary.txt: 0 errors，warnings 为 pre-existing
- work_report.md、audit_report.md、single_run_record 三者口径一致

**Problem attribution:**
- Primary result: `success (with evidence gap)` — qwen3.7-max 在统一失败态上完成了真实修复
- Root cause fixed: `auth.rs` 中 `service_id.is_some() || has_bearer` 让 `ExternalUser` 分支不可达
- Evidence gap: `diff.patch` 是旧阶段“无改动”产物，workspace 无 git 历史，无法独立回放最终 diff
- Score impact: counted with penalty — 当前作为唯一有效 B001 样本计入，但置信度仍低

**Qualitative signals observed (qwen3.7-max):**
- Scope discipline: 5/5 — 仅改 1 个条件、1 个文件
- Diagnostic rigor: 先复现失败，再定位根因，再做全量验证
- Write-back quality: 证据链完整，但 diff 证据仍需补强
- Toolchain competence: Rust 1.95.0, 15-crate workspace, cargo test/fmt/clippy 闭环完成

#### batch_01 glm-5.1 Enterprise Review Summary (2026-05-27)

Reviewer: `hermes-agent (qwen3.7-max)`

**Classification**: `runtime_provider_error` / `provider_route_not_validated`

**Evidence chain:**
- Runner log `batch01-B001-glm-ent.log`: latest rerun now captures provider 400 Unknown Model details after shell logger fix
- Hermes review log `review-batch01-B001-glm-ent.log`: shows EPERM on agent.log, then preparation tool calls only
- Workspace bootstrap: 2026-05-27T15:04:31.404Z — never progressed past initialization
- Git status: no commits, all files untracked — confirms no code execution
- `reports/` dir: residual artifacts from injected snapshot (iter 65-88, dated 2026-05-21) — NOT glm-5.1 enterprise output
- Direct route probes after review: `open.bigmodel.cn/api/coding/paas/v4` with `BIGMODEL_GLM_51_API_KEY` returned 401; `open.bigmodel.cn/api/paas/v4` also returned 401; `proxy_glm_51` returned 502 upstream failure
- Historical evidence: archived remote smoke sample already recorded repeated 401 `令牌已过期或验证不正确` on `https://open.bigmodel.cn/api/coding/paas/v4`

**Workspace integrity:**
- Bug injection verified: `auth.rs:139` has `if service_id.is_some() || has_bearer` (buggy — `ExternalUser` at line 141 is unreachable dead code)
- Workspace is intact and ready for rerun without re-injection

**Problem attribution (4-class breakdown):**
- Model issue: **none** — model never received the task
- Runtime issue: **primary** — configured OpenClaw route returned HTTP 400 "Unknown Model" for `zai/glm-5.1-enterprise`
- Orchestration issue: **secondary** — initial attempt had empty runner log, but latest rerun confirms shell logger fix; remaining issue is no validated enterprise route wired into OpenClaw
- External blocker: **root** — no confirmed working provider path for enterprise lane (configured route 400, bigmodel candidates 401, proxy candidate 502), plus repo-internal route convention drift

**Key difference from glm-5.1 (non-enterprise):** 
- glm-5.1 was blocked by Z.AI provider quota exhaustion (HTTP 429) — a resource/billing issue
- glm-5.1 enterprise was first blocked by provider rejecting the configured model code (HTTP 400), and follow-up candidate routes were also unavailable (bigmodel 401, proxy 502)
- Both are runtime blockers but different sub-categories

**Scoring**: Excluded from all quantitative dimensions. Zero capability data.

**Second attempt outcome (run ID batch01-B001-glm-ent, reviewed 2026-05-27):**
- Enterprise route still returned HTTP 400 "Unknown Model" — unchanged from first attempt
- OpenClaw silently fell through to default model `xfyun/astron-code-latest`, which completed the task successfully
- Independent verification: cargo test 30/30 pass, fmt check pass, clippy 0 new warnings — fix is correct
- However, this run provides ZERO capability data for glm-5.1 enterprise because the enterprise model was never reached
- Classification: `runtime_provider_error` + `model_mismatch`
- Orchestration concern: OpenClaw did not surface model mismatch in run metadata; first-pass human review scores were incorrectly attributed to glm-5.1 enterprise

**Action items:** 
1. **[CRITICAL]** Restore one confirmed-working enterprise provider path first (valid credentials + accepted endpoint/model), then point OpenClaw at that route
   Current conflict: `llm_wiki` says glm-5 / glm-5.1 should route via `api.z.ai`, but runtime registry / rotation config still point BIGMODEL credentials to `open.bigmodel.cn`
2. **[IMPORTANT]** Configure OpenClaw to surface model mismatch when fallback occurs (instead of silently using default) — current silent fallback produces misleading "success" signals in benchmark runs
3. Runner logging 已修复；继续保留当前脚本版本以捕获 provider error details
4. Rerun batch01-B001-glm-ent with confirmed glm-5.1 enterprise route against same workspace (bug injection preserved)

#### batch_01 glm-5.1 Review Summary (2026-05-27)

Reviewer: `hermes-agent (qwen3.7-max)`

**Classification**: `runtime_provider_error` / `provider_quota_exhausted`

**Evidence chain:**
- Runner log `batch01-B001-glm.log`: latest rerun now captures provider 429 / failover details after shell logger fix
- Workspace bootstrap: 2026-05-27T15:05:11.081Z — never progressed past initialization
- Git status: no commits, all files untracked — confirms no code execution
- `reports/` dir: residual artifacts from injected snapshot (iter 65-88, dated 2026-05-21) — NOT glm-5.1 output

**Workspace integrity:**
- Bug injection verified: `auth.rs:139` has `if service_id.is_some() || has_bearer` (buggy — `ExternalUser` at line 141 is unreachable dead code)
- Workspace is intact and ready for rerun without re-injection

**Problem attribution (4-class breakdown):**
- Model issue: **none** — model never received the task
- Runtime issue: **primary** — Z.AI provider returned 429 "Insufficient balance or no resource package"
- Orchestration issue: **secondary** — initial attempt had empty runner log, but latest rerun confirms shell logger fix; remaining blocker is provider quota
- External blocker: **root** — Z.AI account quota exhaustion

**Scoring**: Excluded from all quantitative dimensions. Zero capability data.

**Action item**: Rerun after Z.AI quota restoration. Current shell runner logging is already fixed and should retain provider error details.

## This Month's Mission

`2026-06` 的目标不是复述已有印象，而是通过统一 benchmark task pool，把以下问题做成可量化结论：

- 哪个模型更适合做默认主编码师
- 哪个模型在 Rust / Go / C++ 更强
- 哪个模型在 Web / App / Backend / Agent 更强
- 哪个模型在电商 / 游戏 / 多模态 / 数据平台等不同产品类型更强
- 哪个模型在 UI / 前端 / 后端 / 高并发 / 测试CI / 故障恢复等不同开发环节更强

## Benchmark Coverage Plan

### Core Tasks

- `B001` Rust Service Bugfix
- `B002` Go Backend Incremental Feature
- `B003` C++ CLI / Engine Mini Task
- `B004` Web Dashboard Feature
- `B005` App Closed-Loop Task
- `B006` Pixel-Accurate UI Task
- `B007` High-QPS Queue / Cache Task
- `B008` Deployment Rollback Drill

### Extended Tasks

- `E001` Python Data Pipeline Task
- `E002` Android-Side Service Task
- `E003` iOS Integration Task
- `E004` Turn-Based Game Logic Task
- `E005` RAG Workflow App Task
- `E006` Media Generation API Integration
- `E007` E-Commerce Cart + Checkout Flow

## Run Matrix

最低建议样本量：

- `8 core tasks * 3 models = 24 runs`

可选扩展：

- `+ 2~4 extended tasks * 3 models`

## Quantitative Summary

以下字段待 benchmark 跑完后补齐。

| Metric | qwen3.7-max | glm-5.1 enterprise | glm-5.1 | Notes |
| --- | ---: | ---: | ---: | --- |
| Tasks Run | 1 (success sample) | 2 (runtime_blocked x2) | 1 (runtime_blocked) | B001 only |
| Success Rate | 1/1 (with evidence gap) | 0/2 (runtime_provider_error + model_mismatch) | N/A (runtime_provider_error) | only qwen has a countable sample; glm-ent fallback success not attributable |
| Work Report Write Rate | 1/1 | 0/2 (fallback wrote, but not glm-5.1 enterprise) | 0/1 | glm-ent fallback model wrote report, not enterprise model |
| Real Code Change Rate | 1/1 | 0/2 (blocked before enterprise exec; fallback exec non-attributable) | 0/1 (blocked before exec) | qwen changed only `auth.rs` |
| Avg First Write Iteration | 1 | TBD | N/A | qwen fixed in first implementation pass |
| Avg Audit Pass Rate | 1/1 (success with evidence gap) | N/A (blocked before exec) | N/A |  |
| Avg Blocked Rate | 0/1 | 2/2 (provider route unavailable x2) | 1/1 (provider 429 quota) | blocker types differ |
| Avg Read-Only Rounds | ~10 | 0 | 0 | both glm never started |
| Avg Report-Only Count | 0 | 0 | 0 |  |
| Avg Wall Clock Minutes | TBD | TBD | TBD |  |
| Avg Token Estimate | TBD | TBD | TBD |  |
| Avg Cost | TBD | TBD | TBD |  |

## Capability Matrix Link

详见：

- `capability_matrix.md`

## Open Questions

1. `glm-5.1 enterprise` 是否能在统一 benchmark 中摆脱 OneLink 中观测到的 read-only storm 倾向？
2. `qwen3.7-max` 在 App、C++、高并发、游戏、多模态这些当前缺样本维度中是否依然领先？
3. 纯 `glm-5.1` 在剥离 enterprise/runtime 干扰后，能力是否会高于当前的保守评分？
4. **[NEW — from batch_01]** B001 如何注入可复现的 bug 以确保所有模型的公平评测？建议使用 pre-fix commit 或 bug-introducing patch。

## Action Items

1. **[CRITICAL]** 先恢复一条真实可用的 `glm-5.1 enterprise` provider 路线：当前配置路由 400、bigmodel 候选 401、proxy 候选 502；确认可用后再设置 `OPENCLAW_MODEL_GLM_ENTERPRISE` / OpenClaw provider 映射。
2. **[CRITICAL]** 恢复 `glm-5.1` 的 provider quota / resource package，然后在现有失败态 workspace 上重跑 `B001`。
3. **[IMPORTANT]** 保持当前 runner logging 修复版本：provider HTTP 错误现已能记录到 `runner_logs/openclaw/*.log`。
4. 建立 benchmark 工作区并冻结统一快照。
5. 为每次 run 填写 `single_run_records/`。
6. 每个模型每个 core task 至少跑 1 次。
7. 跑完后更新 `scorecards/`、`capability_matrix.md` 与本月月报。
