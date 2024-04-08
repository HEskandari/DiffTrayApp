ifneq ($(words $(MAKECMDGOALS)),1)
.DEFAULT_GOAL = all
%:
	@$(MAKE) $@ --no-print-directory -rRf $(firstword $(MAKEFILE_LIST))
else
ifndef ECHO
T := $(shell $(MAKE) $(MAKECMDGOALS) --no-print-directory \
      -nrRf $(firstword $(MAKEFILE_LIST)) \
      ECHO="COUNTTHIS" | grep -c "COUNTTHIS")

N := x
C = $(words $N)$(eval N := x $N)
ECHO = echo "`expr " [\`expr $C '*' 100 / $T\`" : '.*\(....\)$$'`%]"
endif

.PHONY: all clean

OUT_DIR = ./out
APP_MACOS = "./Diff Tray App.app/"
APP_WIN = "./Diff Tray App.exe"

all: target
	@$(ECHO) All done

run-build:
	@$(ECHO) "Building..."
	@go build ./...
	@$(ECHO) "✓ Done - Building"

run-cleanup:
	@$(ECHO) "Tidying mod file..."
	@go mod tidy
	@$(ECHO) "✓ Done - Tidying"

run-vet:
	@$(ECHO) "Vetting..."
	@go vet ./...
	@$(ECHO) "✓ Done - Vetting"

test-unit:
	@$(ECHO) "Run unit tests..."
	@go test -v ./...
	@$(ECHO) "✓ Done - Unit Tests"

run-tests: test-unit

run-integration-test: export RUN_INTEGRATION_TESTS=True
run-integration-test:
	@$(ECHO) Should run integration tests: $RUN_INTEGRATION_TESTS
	@$(ECHO) "Running integration tests..."
	@go test -v ./... -run Integration
	@$(ECHO) "✓ Done - Integration Tests"

create-package-macos:
	@$(ECHO) "Packing for MacOS..."
	@mkdir -p $(OUT_DIR)
	@fyne package -os darwin -icon icon.png
	@-mv $(APP_MACOS) $(OUT_DIR)
	@$(ECHO) "✓ Done - Packaging MacOS"

create-package-windows:
	@$(ECHO) "Packing for Windows..."
	@mkdir -p $(OUT_DIR)
	@fyne package -os windows -icon icon.png
	@-mv $(APP_WIN) $(OUT_DIR)
	@$(ECHO) "✓ Done - Packaging Windows"

package-all: create-package-macos create-package-windows

target: run-build run-vet run-tests create-package-macos

endif