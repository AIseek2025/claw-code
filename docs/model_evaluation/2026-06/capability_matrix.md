# Capability Matrix

## Scope

本文用于记录 `2026-06` benchmark 执行后的多维能力判断。

在 benchmark 尚未跑完之前，本文件只保留：

- 当前待验证维度
- 每个维度的 benchmark 对应任务
- 当前结论状态

状态说明：

- `TBD`: 已纳入 benchmark，但尚未形成结论
- `no conclusion`: 当前无样本或样本仍不足
- `provisional`: 已有少量样本，但结论不稳定
- `validated`: 已有足够 benchmark 样本支撑

## Language Matrix

| Language | qwen3.7-max | glm-5.1 enterprise | glm-5.1 | Current Winner | Status | Benchmark Task |
| --- | --- | --- | --- | --- | --- | --- |
| Rust | provisional (success, 1 sample) | no conclusion (runtime_blocked x2 — provider 400 + model_mismatch on fallback, reviewed 2026-05-27) | no conclusion (runtime_blocked — provider 429, reviewed 2026-05-27) | qwen3.7-max (sole sample) | provisional | `B001` — qwen3.7-max: success (evidence gap on diff); glm-5.1 enterprise: runtime_blocked x2 (model code 400, fallback to xfyun/astron-code-latest produced non-attributable success); glm-5.1: runtime_blocked (Z.AI quota) |
| Go | TBD | TBD | TBD | TBD | TBD | `B002` |
| C++ | TBD | TBD | TBD | TBD | TBD | `B003` |
| TypeScript / JavaScript | TBD | TBD | TBD | TBD | TBD | `B004`, `B006` |
| Python | TBD | TBD | TBD | TBD | no conclusion | `E001`, `E005`, `E006` |
| Java / Kotlin | TBD | TBD | TBD | TBD | no conclusion | `E002` |
| Swift / Objective-C | TBD | TBD | TBD | TBD | no conclusion | `E003` |

## Product Surface Matrix

| Surface | qwen3.7-max | glm-5.1 enterprise | glm-5.1 | Current Winner | Status | Benchmark Task |
| --- | --- | --- | --- | --- | --- | --- |
| Web | TBD | TBD | TBD | TBD | TBD | `B004`, `B006`, `E007` |
| App | TBD | TBD | TBD | TBD | TBD | `B005`, `E002`, `E003` |
| Backend / API | provisional (success, 1 sample) | no conclusion (runtime_blocked x2 — provider 400 + model_mismatch) | no conclusion (runtime_blocked — provider 429) | qwen3.7-max (sole sample) | provisional | `B001` — qwen3.7-max success; glm-5.1 enterprise runtime_blocked x2 (model code 400, fallback non-attributable), glm-5.1 runtime_blocked (Z.AI 429), `B002`, `B007` |
| CLI / Tooling | TBD | TBD | TBD | TBD | TBD | `B003` |
| Agent / Workflow | TBD | TBD | TBD | TBD | TBD | `E005` |

## Product Domain Matrix

| Domain | qwen3.7-max | glm-5.1 enterprise | glm-5.1 | Current Winner | Status | Benchmark Task |
| --- | --- | --- | --- | --- | --- | --- |
| 电商 | TBD | TBD | TBD | TBD | no conclusion | `E007` |
| 游戏 | TBD | TBD | TBD | TBD | no conclusion | `E004` |
| 智能体产品 | TBD | TBD | TBD | TBD | no conclusion | `E005` |
| 垂直行业大模型应用 | TBD | TBD | TBD | TBD | no conclusion | `E005` |
| 图片 / 视频生成 | TBD | TBD | TBD | TBD | no conclusion | `E006` |
| 大数据 / 数据平台 | TBD | TBD | TBD | TBD | no conclusion | `E001` |
| SaaS / 企业软件 | provisional (success, 1 sample) | no conclusion (runtime_blocked x2 — provider 400 + model_mismatch) | no conclusion (runtime_blocked — provider 429) | qwen3.7-max (sole sample) | provisional | `B001` — qwen3.7-max success; glm-5.1 enterprise runtime_blocked x2 (model code 400, fallback non-attributable), glm-5.1 runtime_blocked (Z.AI 429), `B002`, `B004`, `B008` |

## Development Stage Matrix

| Stage | qwen3.7-max | glm-5.1 enterprise | glm-5.1 | Current Winner | Status | Benchmark Task |
| --- | --- | --- | --- | --- | --- | --- |
| UI / 视觉实现 | TBD | TBD | TBD | TBD | TBD | `B006` |
| 前端交互 | TBD | TBD | TBD | TBD | TBD | `B004`, `B005` |
| 后端业务逻辑 | provisional (success, 1 sample) | no conclusion (runtime_blocked x2 — provider 400 + model_mismatch) | no conclusion (runtime_blocked — provider 429) | qwen3.7-max (sole sample) | provisional | `B001` — qwen3.7-max success; glm-5.1 enterprise runtime_blocked x2 (model code 400, fallback non-attributable), glm-5.1 runtime_blocked (Z.AI 429), `B002` |
| 数据库 / 存储 | no conclusion (blocked) | no conclusion (blocked) | no conclusion (blocked) | none | no conclusion | `B001` (glm blocked; qwen3.7-max B001 was auth bugfix, not storage-focused), `E001` |
| 高并发架构 | TBD | TBD | TBD | TBD | TBD | `B007` |
| 契约 / API 设计 | TBD | TBD | TBD | TBD | TBD | `B002`, `E007` |
| 测试 / CI | provisional (test closure verified) | no conclusion (runtime_blocked x2 — provider 400 + model_mismatch) | no conclusion (runtime_blocked — provider 429) | qwen3.7-max (sole sample) | provisional | `B001` — qwen3.7-max: full test/fmt/clippy closure; glm-5.1 enterprise runtime_blocked x2 (model code 400, fallback non-attributable), glm-5.1 runtime_blocked (Z.AI 429), `B002`, `B004` |
| 运维 / 交付 / 故障恢复 | TBD | TBD | TBD | TBD | TBD | `B008` |

## Update Rule

每完成一轮 benchmark，至少更新：

1. 各模型该 task 的单轮记录
2. 对应维度的 `TBD -> provisional / validated`
3. `Current Winner`
4. 支撑结论的 evidence 路径
