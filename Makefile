# List building
ALL_LIST = crittercism_telemetry_agent.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -race
GOFMT=gofmt -w

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)
RUN_LIST = $(foreach int, $(ALL_LIST), $(int)_run)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(BUILD_LIST)

all: build
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
test: $(TEST_LIST)
fmt: $(FMT_TEST)
run: $(RUN_LIST)

$(BUILD_LIST): %_build: %_fmt
	@if [ -f ./prebuild ]; then \
		echo "Running prebuild script in release mode..." ; \
		./prebuild --release ; \
	fi
	@echo "Building Linux AMD64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/linux-amd64/$*
	@echo "Building Darwin AMD64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/darwin-amd64/$*
	@echo "Building Windows AMD64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/windows-amd64/$*
	@if [ -f ./prebuild ]; then \
		echo "Running prebuild script in release mode..." ; \
		./prebuild --debug ; \
	fi

$(TEST_LIST): %_test:
	@echo "Running go test..."
	@$(GOTEST) ./...

$(FMT_TEST): %_fmt:
	@echo "Running go fmt..."
	@$(GOFMT) crittercism_telemetry_agent.go crittercism
