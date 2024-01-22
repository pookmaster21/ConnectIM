TARGET_EXEC=ConnectIM

BUILD_DIR=./bin

default: build

build:
	@echo "Building"
	@go build -o $(BUILD_DIR)/$(TARGET_EXEC)
	@echo "Finished building"

run: build
	@echo "Running"
	@$(BUILD_DIR)/$(TARGET_EXEC)
	@echo "Finished running"

clean:
	@echo "Removing"
	@rm -rf $(BUILD_DIR)
	@echo "Finished removing"
