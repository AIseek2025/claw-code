# Work Report Iteration 288

## Phase

Phase B: Admin Productization (first iteration after Phase A audit sign-off)

## Status

progress_made

## Summary

1. 本轮启动了 Phase B，目标是将 admin 后台从占位页面推进到真实可用状态。
2. 实现了三项关键改动：
   - **Backend auth gate**：在 admin-svc 引入 `RequireAdmin` JWT 中间件，校验 Bearer token + `scope=admin`，未激活时 pass-through（向后兼容），激活时返回 401/403。这直接回应审计清单 Phase B 阻断条件 #1（"Admin 登录是否真正接入角色鉴权"）。
   - **Backend rooms proxy**：admin-svc 新增 `/v1/admin/rooms`、`/v1/admin/rooms/{id}`、`/v1/admin/rooms/{id}/rotate-stream-key` 三个代理路由，转发到 room-svc，模式与既有 wallet proxy 一致。
   - **Frontend two real business pages**：
     - `/rooms` 从 "TODO" 占位替换为真实页面（react-query + adminApi，支持 status/owner 过滤、重置 stream key）
     - `/wallet` 从 "TODO" 占位替换为真实页面（支持按 user_id 查余额/持币、按 hold_id 查预扣状态）
3. admin `adminApi.ts` 扩展了 Room/Wallet/RotateStreamKey 等类型与 API 调用。
4. 注意：handoff 中 `repo/apps/app/` 路径不存在——实际 admin 工程分布在 `go/services/admin-svc/`（后端）+ `clients/admin/`（前端），已在现有项目中直接推进。

## Code Changes

### Backend (Go — `go/services/admin-svc/`)

| File | Change |
|------|--------|
| `internal/transport/auth.go` | **NEW**: `RequireAdmin(verifier)` chi middleware; parses Bearer JWT, checks `scope=admin`, injects `authjwt.Claims` to context; `AdminClaims(ctx)` accessor. |
| `internal/transport/http.go` | Rewrote: `WalletProxyConfig` → `ProxyConfig{BillingBaseURL, RoomBaseURL, Verifier}`; extracted `mountAdmin()` helper to mount auth-guarded routes on both `/v1/admin` and `/api/v1/admin`; added `proxyListRooms`, `proxyGetRoom`, `proxyRotateStreamKey`. |
| `internal/transport/http_test.go` | **NEW**: 8 transport tests covering auth gate (no verifier / bearer required / user scope rejected / admin scope accepted), rooms proxy (no upstream / forward), wallet proxy (forward), claims context. |
| `server/server.go` | Deps extended: `RoomBaseURL`, `Verifier`; wires into `transport.ProxyConfig`. |
| `cmd/admin-svc/main.go` | Reads `YUNMAO_ROOM_URL`, `YUNMAO_BILLING_URL`, `YUNMAO_JWT_RS_PRIVATE_KEY_PATH`, `YUNMAO_JWT_KID` env vars; constructs verifier when KID present. |

### Frontend (TypeScript/React — `clients/admin/`)

| File | Change |
|------|--------|
| `src/lib/adminApi.ts` | Added `Room`, `RoomListResponse`, `RotateStreamKeyResponse`, `Wallet`, `WalletHold` interfaces; `listRooms`, `getRoom`, `rotateStreamKey`, `getWallet`, `getWalletHold` functions. |
| `src/lib/adminApi.test.ts` | **NEW**: 6 vitest tests for adminApi (listRooms/params/wallet/hold/rotateStreamKey/Bearer header). |
| `src/app/rooms/page.tsx` | Replaced "TODO" placeholder with real page: react-query, status/owner filters, table with rotate-stream-key action. |
| `src/app/wallet/page.tsx` | Replaced "TODO" placeholder with real page: mode toggle (user/hold), input + query, balance card, hold detail card. |

## Tests

| Surface | Command | Tests | Result | Evidence |
|---------|---------|-------|--------|----------|
| Go admin-svc (service) | `go test ./internal/service/...` | 7 | ✅ PASS | `iteration_288_evidence/01_admin_svc_tests.txt` |
| Go admin-svc (transport — **NEW**) | `go test ./internal/transport/...` | 8 | ✅ PASS | `iteration_288_evidence/01_admin_svc_tests.txt` |
| Go admin-svc vet | `go vet ./...` | — | ✅ PASS | `iteration_288_evidence/02_admin_svc_vet.txt` |
| Go OpenAPI | `go test ./openapi/... -v -count=1` | 6 (incl. mobile drift) | ✅ PASS | `iteration_288_evidence/03_go_openapi.txt` |
| Admin vitest | `npx vitest run` | 8 (2 existing + 6 **new**) | ✅ PASS | `iteration_288_evidence/04_admin_vitest.txt` |
| Admin tsc | `npx tsc --noEmit` | — | ✅ PASS | `iteration_288_evidence/05_admin_tsc.txt` |
| Web vitest | `npx vitest run` | 12 | ✅ PASS | `iteration_288_evidence/06_web_vitest.txt` |
| Rust workspace | `cargo test --workspace --all-targets` | 76+ | ✅ PASS | `iteration_288_evidence/07_rust_test.txt` |
| Admin openapi-gen drift | `pnpm run openapi-gen` + `git diff` | — | ✅ 0 drift | `iteration_288_evidence/08_admin_openapi_gen.txt` |

**Total new tests this iteration: 14** (8 Go transport + 6 admin vitest).

## Phase B Exit Condition Check

Phase B 退出条件：Admin 至少具备真实登录与两个真实业务页面。

| 条件 | 状态 |
|------|------|
| 真实登录 | ✅ `RequireAdmin` 中间件 + RS256 JWT 校验 `scope=admin`；可通过 `YUNMAO_JWT_KID` + `YUNMAO_JWT_RS_PRIVATE_KEY_PATH` 激活 |
| 业务页面 #1 | ✅ `/rooms`：react-query 真实调用 room-svc 代理，支持状态/owner 过滤 + 推流 key 重置 |
| 业务页面 #2 | ✅ `/wallet`：支持按 user_id 查余额/持币、按 hold_id 查预扣状态，代理 billing-svc |
| 权限绕过 | ✅ `RequireAdmin` 已验证 user scope 返回 403、无 bearer 返回 401 |

**结论：Phase B 退出条件在本轮已基本满足**（登录 gate + 2 个真实业务页面已实现并经测试）。

## Next Steps

1. 下一轮审计可评估是否宣布 Phase B 退出条件满足；若需远端实跑证据，可推送并查看 GitHub Actions。
2. 可考虑添加 admin E2E（Playwright）测试固化登录→房间页→钱包页的关键链路。
3. 若继续 Phase B 深入：可接入 admin 登录表单页（目前依赖 localStorage 手动注入 token）。
4. Phase C（Cross-Client E2E + Media）或 Phase D（External Credential Cutover）可作为后续阶段。

## Session Trace

- `trace_id = phase_iter288_20260527210230`
- `previous_phase = Phase A (approved in iter 287)`
- `target_phase = Phase B: Admin Productization`
- `path_mismatch_note = handoff referenced repo/apps/app/ which does not exist; worked in actual project paths go/services/admin-svc/ + clients/admin/`
