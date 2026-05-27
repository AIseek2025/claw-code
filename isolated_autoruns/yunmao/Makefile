# yunmao monorepo Makefile
# 顶层 make 入口，负责编排 rust/ + go/ + deploy/

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

# 兼容未把 ~/.local/go/bin 加入 PATH 的开发机
export PATH := $(HOME)/.local/go/bin:$(HOME)/.cargo/bin:$(HOME)/go/bin:$(PATH)

.DEFAULT_GOAL := help

.PHONY: help fmt lint test \
	rust-fmt rust-lint rust-test rust-build \
	go-fmt go-vet go-test go-build go-tidy \
	dev-up dev-down dev-logs poc-feed bench-ws gen-proto \
	buf-lint buf-breaking buf-generate \
	migrate-up migrate-down \
	app-up app-down web-demo e2e integration \
	web-dev web-build web-test admin-dev admin-build admin-test android-build ios-build \
	openapi-lint openapi-gen openapi-test

help:
	@echo "yunmao monorepo targets:"
	@echo "  make fmt              - 格式化 rust + go"
	@echo "  make lint             - clippy + go vet/golangci-lint(可选)"
	@echo "  make test             - 跑 rust 与 go 全部单测"
	@echo "  make rust-{fmt,lint,test,build}"
	@echo "  make go-{fmt,vet,test,build,tidy}"
	@echo "  make dev-up | dev-down | dev-logs"
	@echo "  make app-up | app-down   - 起 / 停 yunmao 应用层 (rust+go) docker compose"
	@echo "  make web-demo            - 启动 clients/web-demo 静态服务 (http://localhost:5173)"
	@echo "  make e2e                 - 跑 scripts/e2e.sh 端到端 smoke（要先 dev-up + app-up）"
	@echo "  make migrate-up | migrate-down - 用 psql 跑 / 回滚 go/migrations"
	@echo "  make buf-{lint,breaking,generate}"
	@echo "  make poc-feed         - 触发一次端到端投喂事件链路 PoC"
	@echo "  make bench-ws         - WebSocket 网关基准压测"
	@echo "  make gen-proto        - 通过 buf 生成 proto stub"

fmt: rust-fmt go-fmt
lint: rust-lint go-vet
test: rust-test go-test

rust-fmt:
	cd rust && cargo fmt --all

rust-lint:
	cd rust && cargo clippy --workspace --all-targets -- -D warnings

rust-test:
	cd rust && cargo test --workspace

rust-build:
	cd rust && cargo build --workspace --all-targets

GO_MODULES := pkg/yunmao proto services/user-svc services/room-svc services/feeding-svc services/device-svc services/billing-svc services/admin-svc services/chat-svc

go-fmt:
	cd go && gofmt -w .

go-vet:
	@cd go && for m in $(GO_MODULES); do echo "vet $$m"; (cd $$m && go vet ./...); done

go-test:
	@cd go && for m in $(GO_MODULES); do echo "test $$m"; (cd $$m && go test ./...); done

go-build:
	@cd go && for m in $(GO_MODULES); do echo "build $$m"; (cd $$m && go build ./...); done

go-tidy:
	@cd go && for m in $(GO_MODULES); do echo "tidy $$m"; (cd $$m && go mod tidy); done

dev-up:
	docker compose -f deploy/docker-compose.dev.yml up -d

dev-down:
	docker compose -f deploy/docker-compose.dev.yml down -v

dev-logs:
	docker compose -f deploy/docker-compose.dev.yml logs -f --tail=200

poc-feed:
	bash scripts/poc-feed.sh

bench-ws:
	bash scripts/bench-ws.sh

gen-proto:
	bash scripts/gen-proto.sh

buf-lint:
	cd proto && buf lint

buf-breaking:
	@echo "buf breaking against origin/main..."
	cd proto && buf breaking --against ".git#branch=origin/main,subdir=proto" || \
	  echo "(skipped: no origin/main baseline available)"

buf-generate:
	cd proto && buf generate

# psql 必须在宿主可用；docker-compose 启动后 5432 暴露在 localhost。
PG_URL ?= postgres://yunmao:yunmao@localhost:5432/yunmao?sslmode=disable

migrate-up:
	@for f in $$(ls go/migrations/*.sql | sort); do \
	  echo "applying $$f"; \
	  psql "$(PG_URL)" -v ON_ERROR_STOP=1 -f $$f; \
	done

migrate-down:
	@echo "no automatic down-migration yet; use psql + DROP TABLE manually"

# Docker Compose 应用层
app-up:
	docker compose -f deploy/docker-compose.dev.yml -f deploy/docker-compose.app.yml up -d --build

app-down:
	docker compose -f deploy/docker-compose.dev.yml -f deploy/docker-compose.app.yml down

# 静态 web-demo（最小回归页；DEPRECATED，正式 web 见 clients/web/）
web-demo:
	@echo "serving clients/web-demo at http://localhost:5173 (Ctrl-C to stop)"
	@cd clients/web-demo && python3 -m http.server 5173

# 第八轮（G）：clients/web Next.js 15 正式工程
web-dev:
	cd clients/web && pnpm install && pnpm dev

web-build:
	cd clients/web && pnpm install --frozen-lockfile && pnpm build

web-test:
	cd clients/web && pnpm install && pnpm test:run

# 第八轮（H）：clients/admin 运营后台
admin-dev:
	cd clients/admin && pnpm install && pnpm dev

admin-build:
	cd clients/admin && pnpm install --frozen-lockfile && pnpm build

admin-test:
	cd clients/admin && pnpm install && pnpm test:run

# 第八轮（F / E）：客户端骨架编译入口（本机若无对应工具链会失败，CI 跑）
android-build:
	cd clients/android && ./gradlew assembleDebug

ios-build:
	@if ! command -v xcodebuild >/dev/null 2>&1; then echo "xcodebuild not found"; exit 0; fi
	cd clients/ios/YunmaoApp && xcodebuild -scheme YunmaoApp -destination 'platform=iOS Simulator,name=iPhone 15' build

# 端到端 smoke：需要先 `make dev-up && make app-up`。
# 自动跑登录→订阅 token→投喂→详情查询→LL-HLS playlist→各服务 metrics。
e2e:
	@bash scripts/e2e.sh

# testcontainers-go 集成测试：需要本机 docker 可用，CI 走 ubuntu runner。
# INTEGRATION=1 启用，否则 t.Skip。
# 第七轮（D）：--results-out 把 JSON/TXT 落到 reports/integration/。
integration:
	@mkdir -p reports/integration
	@ts=$$(date +%Y%m%d-%H%M%S); \
	  out="reports/integration/integration-$$ts"; \
	  echo "[integration] writing $$out.log + $$out.json"; \
	  cd go/test/e2e && INTEGRATION=1 go test -tags=integration ./... \
	    -v -timeout 15m -json 2>&1 \
	    | tee "../../../$$out.log" \
	    | awk '/^{/ {print > "/dev/stderr"; print}' \
	    > "../../../$$out.json" || true; \
	  echo "[integration] artifacts: $$out.log $$out.json"

# integration-up：起 dev infra + 跑 integration 测试 + 归档报告。
# 报告：reports/integration/<date>.log
.PHONY: integration-up perf-baseline
integration-up:
	@mkdir -p reports/integration
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "docker not found; skipping integration-up"; exit 0; \
	fi
	@if ! docker info >/dev/null 2>&1; then \
		echo "docker daemon not available; skipping"; exit 0; \
	fi
	@ts=$$(date +%Y%m%d-%H%M%S); \
	  log="reports/integration/$$ts.log"; \
	  echo "[integration-up] log=$$log"; \
	  $(MAKE) dev-up >>$$log 2>&1; \
	  cd go/test/e2e && INTEGRATION=1 go test -tags=integration ./... -v -timeout 15m 2>&1 | tee -a "../../../$$log"

# perf-baseline：一键拉 docker compose + 跑 bench_ws + 写 reports/perf/ws-baseline-<date>.md
perf-baseline:
	@mkdir -p reports/perf
	bash scripts/perf/ws-baseline-all.sh $${YUNMAO_BENCH_CONNS:-10000} $${YUNMAO_BENCH_DURATION_SECS:-60}

# 第七轮（D / I）：integration-up-remote + chat-baseline 工件。
# integration-up-remote：在 macOS 上用 DOCKER_HOST 指向远端 docker（Colima / 远端 Linux）。
integration-up-remote:
	@if [ -z "$$DOCKER_HOST" ]; then \
		echo "DOCKER_HOST is empty; set to a remote docker, e.g. tcp://192.168.1.10:2375"; exit 1; \
	fi
	@$(MAKE) integration-up

# chat-baseline：起 chat-svc + gateway + Redis 模拟器，跑 N 用户 / M 观众压测。
chat-baseline:
	@mkdir -p reports/perf
	bash scripts/perf/chat-baseline-up.sh $${YUNMAO_CHAT_USERS:-1000} $${YUNMAO_CHAT_VIEWERS:-5000}

# 第九轮 Phase 1：OpenAPI 共享契约输出 + 客户端消费链路。
# openapi-lint：校验 v3.json 结构完整性（JSON 格式 + 必需字段）。
openapi-lint:
	@echo "[openapi-lint] validating go/pkg/yunmao/openapi/v3.json..."
	@cd go/pkg/yunmao && go test ./openapi/... -v -count=1

# openapi-test：跑 openapi_test.go（等价于 openapi-lint，CI 入口）。
openapi-test:
	@cd go/pkg/yunmao && go test ./openapi/... -v -count=1

# openapi-gen：从 v3.json 生成客户端 TypeScript 类型。
# 目标：clients/web/src/lib/generated-api.ts + clients/admin/src/lib/generated-api.ts。
openapi-gen:
	@echo "[openapi-gen] generating TypeScript types from v3.json..."
	@cd clients/web && pnpm run openapi-gen
	@echo "[openapi-gen] done: clients/web/src/lib/generated-api.ts"
	@cd clients/admin && pnpm run openapi-gen
	@echo "[openapi-gen] done: clients/admin/src/lib/generated-api.ts"

