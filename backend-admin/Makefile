# 编译相关
FLAGS=-tags=jsoniter
DEV_FLAGS=-gcflags="-l -N" -race
GOOS=linux
GOARCH=amd64

# 项目相关
NAME=golang-common-base
# main.go 路径
MAIN_PATH=./app/cmd/
PORT=9088
COMMIT=$(shell git log -1 --pretty=format:%h)
DEV_NAME=$(NAME)-$(COMMIT)
RELEASE_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
RELEASE_NAME=$(NAME)-$(RELEASE_VERSION)
ifdef VERSION
RELEASE_VERSION=$(VERSION)
endif

# docker相关
DOCKER_REGISTRY=harbor.5qipa.com:6443/common
DOCKER_TARGET=$(DOCKER_REGISTRY)/$(NAME):$(RELEASE_VERSION)


.PHONY: all
all: build-dev

build-dev:
	go build -o $(NAME) $(FLAGS) $(DEV_FLAGS) $(MAIN_PATH)
	@echo "$(NAME) build okay"

build-release:
	 CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(NAME) $(FLAGS) $(MAIN_PATH)
	@echo "$(NAME) build okay"

start: build-dev
	./$(NAME)

pack-dev: build-dev
	@mkdir -p _dev
	@cp  $(NAME) _dev
	cd _dev && tar czf $(DEV_NAME).tar.gz *
	@echo "pack $(DEV_NAME).tar.gz okay"


pack-release: build-release
	@mkdir -p _release
	@cp $(NAME) _release
	@cp -r ./conf _release
	cd _release && tar czf $(RELEASE_NAME).tar.gz *
	@echo "pack $(RELEASE_NAME).tar.gz okay"

clean:
	@rm -rf $(NAME)
	@rm -rf _dev
	@rm -rf _release
	@rm -rf _docker
	@rm -rf log
	@echo "clean okay"

docker-build: clean pack-release
	@mkdir -p _docker
	@cp -f Dockerfile _docker
	@cp _release/$(RELEASE_NAME).tar.gz _docker/
	@cp docker/Shanghai _docker/
	cd _docker && docker buildx build --platform linux/amd64 --no-cache -t $(DOCKER_TARGET) --build-arg modeenv=$(ENV_SERVER_MODE) --build-arg exposeport=$(PORT) --build-arg procname=$(NAME) --build-arg packagefile=$(RELEASE_NAME).tar.gz .
	@echo "docker-build okay"

docker-clean:
	docker rmi $(DOCKER_TARGET)

docker-push:
	docker push $(DOCKER_TARGET)

docker-all: docker-build docker-push clean
