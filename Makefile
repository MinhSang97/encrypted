.PHONY: build 
# include Makefile-deploy-webhook-lambda.mk
# include Makefile-deploy-activity-log-lambda.mk
# include Makefile-deploy-conclusion-esignature-granted-lambda.mk
# include Makefile-deploy-conclusion-contract-concluded-lambda.mk
# include Makefile-deploy-conclusion-contract-failed-lambda.mk

SRC_PATH:= ${PWD}

# JAWSDB_URL is ENV variable auto generated of heroku deployment
# mysql://e7w86uyn0hyzvnr5:cecvj9al6zbdz02e@phtfaw4p6a970uc0.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/xog7a44zta9of1st
DB_CONFIG_PROFILE := $(if ${JAWSDB_URL},heroku,mysql_env)
ifeq (${JAWSDB_URL},)
else
ConStr=${JAWSDB_URL}
ConStr:=$(subst mysql://,,${ConStr})
TCP_PREFIX:=@tcp(
PORT_SUFFIX:=)/
ConStr:=$(subst @,${TCP_PREFIX},${ConStr})
ConStr:=$(subst /,${PORT_SUFFIX},${ConStr})
JAWSDB_URL:=${ConStr}?parseTime=true&charset=utf8mb4
endif

## Developing jobs
prepare:
#	@go install github.com/cosmtrek/air@latest
	@go install github.com/jstemmer/go-junit-report@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/daixiang0/gci@latest
	@go install github.com/rubenv/sql-migrate/...@latest
	@go install github.com/rillig/gobco@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${shell go env GOPATH}/bin v1.55.0

mod:
	@go mod tidy
	@go mod vendor

gen:
	## Go generate
	@go generate ./...
	## Swagger generate
	@swag init -g app/external/framework/route.go -o app/interface/api/docs  --exclude pkg,db,deployment,scripts,vendor
	@./scripts/gci.sh

fmt: ## gofmt and goimports all go files
	@find . -name '*.go' -not -wholename './vendor/*' -not -wholename '*_gen.go' -not -wholename '*/mock_*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

build:
	@go build -tags=jsoniter -o eboost-api-partner ${SRC_PATH}/cmd/srv/...

cron:
	@go build -o eboost-api-partner-cron ${SRC_PATH}/cmd/cron/...

up:
	@go run ${SRC_PATH}/cmd/srv/...	

run:
	air --build.cmd "go build -o tmp/runner-build ${SRC_PATH}/cmd/srv/..." --build.bin "./tmp/runner-build" --build.include_ext "go,tpl,html,tmpl" --build.exclude_dir "cronjobs,db,assets,scripts,test-results,tools,vendor"

debug:
	@go run ${SRC_PATH}/cmd/srv/... 2>&1 | tee running.log

lint: ## Run lint
	@./scripts/lint.sh

lint-sonar:
	@./scripts/lint-sonar.sh

test: ## Run unit testing
	@./scripts/test.sh

test-coverage:
	@go tool cover -html=test-results/.testCoverage.txt

test-sonar:
	@./scripts/test-sonar.sh

branch-test-sonar: ## Run branch coverage testing for Sonarqube
	@./scripts/branch-test-sonar.sh

branch-test-all: ## Run branch coverage testing for whole source
	@./scripts/branch-test-all.sh

sec:
	@gosec -exclude-dir=mock ./...

gci:
	@./scripts/gci.sh

## setup command
check: mod lint sec test build

init: prepare mod up-db migrate-up

migrate-create:
	@$(eval NAME := $(shell read -p "Enter new file name: " v && echo $$v))
	$(eval CMD:= $*)
	cd db;\
	sql-migrate new ${NAME}

# [up,down]
migrate-%:
	$(eval CMD:= $*)
	cd db;\
	sql-migrate $(CMD) -config=dbconfig.yml;

migrate:
ifeq (${JAWSDB_URL},)
ifeq (${MYSQL_USER},)
	$(error env MYSQL_USER is empty)
endif
ifeq ($(MYSQL_PASSWORD),)
	$(error env MYSQL_PASSWORD is empty)
endif
ifeq ($(MYSQL_HOST),)
	$(error env MYSQL_HOST is empty)
endif
ifeq ($(MYSQL_PORT),)
	$(error env MYSQL_PORT is empty)
endif
ifeq ($(MYSQL_DB),)
	$(error env MYSQL_DB is empty)
endif
endif

	cd db;\
	sql-migrate up -env=${DB_CONFIG_PROFILE} -config=dbconfig.yml

## Docker/Container jobs
up-db:
	@docker-compose up -d

PPROF_TYPE := "profile"
pprof:
	@go tool pprof http://localhost:30000/debug/pprof/$(PPROF_TYPE)

sonar-scan:
	@./scripts/local-test-sonar.sh $(test)
	@./scripts/local-scan-sonar.sh
