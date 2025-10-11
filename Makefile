SHELL := /bin/bash
MAKEFLAGS += --warn-undefined-variables
.DEFAULT_GOAL := help

# Toolchain configuration
CC ?= clang
CXX ?= clang++
AR ?= llvm-ar
ARFLAGS ?= rcs
PKG_CONFIG ?= pkg-config

# Build layout
BUILD_TYPE ?= debug
SRC_DIR ?= src
TEST_DIR ?= tests
BUILD_DIR ?= build
OBJ_DIR := $(BUILD_DIR)/obj
TEST_OBJ_DIR := $(BUILD_DIR)/tests/obj
TEST_BIN := $(BUILD_DIR)/tests/all_tests

# Compiler flags
CFLAGS_COMMON := -std=c23 -Wall -Wextra -Wpedantic -Wstrict-overflow -Wundef -Winline -Wimplicit-fallthrough -Wformat=2 -Wvla -march=native
CFLAGS_DEBUG := -g --save-temps
CFLAGS_RELEASE := -O2 -DNDEBUG -Werror #-fsanitize=address -fsanitize=bounds

ifeq ($(BUILD_TYPE),release)
CFLAGS := $(CFLAGS_COMMON) $(CFLAGS_RELEASE)
else
CFLAGS := $(CFLAGS_COMMON) $(CFLAGS_DEBUG)
endif

CPPFLAGS ?= -I$(SRC_DIR)
CXXFLAGS ?= -std=c++17 -Wall -Wextra -Wpedantic -Wstrict-overflow -Wundef -Winline -Wimplicit-fallthrough -Wformat=2 -Wvla -march=native -Werror
LDFLAGS ?=
LDLIBS ?=

# GoogleTest discovery (falls back to common linker flags if pkg-config is unavailable)
GTEST_AVAILABLE := $(shell $(PKG_CONFIG) --exists gtest && echo yes || echo no)
ifeq ($(GTEST_AVAILABLE),yes)
GTEST_CFLAGS := $(shell $(PKG_CONFIG) --cflags gtest)
GTEST_LIBS := $(shell $(PKG_CONFIG) --libs gtest_main)
else
GTEST_CFLAGS :=
GTEST_LIBS := -lgtest -lgtest_main -pthread
endif

# Source discovery
C_SRCS := $(shell test -d $(SRC_DIR) && fd . $(SRC_DIR) -t f -e c 2>/dev/null)
C_OBJS := $(addsuffix .o,$(basename $(patsubst $(SRC_DIR)/%, $(OBJ_DIR)/%, $(C_SRCS))))

TEST_SRCS := $(shell test -d $(TEST_DIR) && fd -g "*.test.{c,cpp}" $(TEST_DIR) -t f 2>/dev/null)
TEST_OBJS := $(addsuffix .o,$(basename $(patsubst $(TEST_DIR)/%, $(TEST_OBJ_DIR)/%, $(TEST_SRCS))))

ALL_OBJS := $(C_OBJS) $(TEST_OBJS)
DEPS := $(ALL_OBJS:.o=.d)

.PHONY: help test build debug release clean dirs format

help: ## Display available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: dirs $(C_OBJS) ## Build all source files

# Toggle build types on demand
debug: BUILD_TYPE := debug
debug: clean build

release: BUILD_TYPE := release
release: clean build

# Test orchestration
ifneq ($(strip $(TEST_SRCS)),)
test: $(TEST_BIN) ## Build and run all tests
	"$(TEST_BIN)"
else
test:
	@echo "error: no GoogleTest sources (*.c/*.cpp) found under '$(TEST_DIR)'" >&2
	@exit 1
endif

ifneq ($(strip $(TEST_SRCS)),)
$(TEST_BIN): $(C_OBJS) $(TEST_OBJS)
	@mkdir -p $(dir $@)
	$(CXX) $(LDFLAGS) $^ $(GTEST_LIBS) $(LDLIBS) -o $@
endif

# Object compilation rules
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c
	@mkdir -p $(dir $@)
	$(CC) $(CPPFLAGS) $(CFLAGS) -MMD -MP -c $< -o $@

$(TEST_OBJ_DIR)/%.o: $(TEST_DIR)/%.c
	@mkdir -p $(dir $@)
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) $(GTEST_CFLAGS) -MMD -MP -c $< -o $@

$(TEST_OBJ_DIR)/%.o: $(TEST_DIR)/%.cpp
	@mkdir -p $(dir $@)
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) $(GTEST_CFLAGS) -MMD -MP -c $< -o $@

# Housekeeping targets
clean: ## Remove all build artefacts
	rm -rf "$(BUILD_DIR)"

dirs:
	@mkdir -p "$(OBJ_DIR)" "$(TEST_OBJ_DIR)"

format: ## Format all source files using clang-format
	@if ! command -v clang-format &> /dev/null; then \
		echo "error: clang-format not found, please install it to use this target" >&2; \
		exit 1; \
	fi
	@fd . $(SRC_DIR) -t f -e c -e h -e cpp 2>/dev/null | xargs clang-format -i
	@fd -g "*.test.{c,cpp}" $(TEST_DIR) -t f 2>/dev/null | xargs clang-format -i


ifneq ($(MAKECMDGOALS),clean)
-include $(DEPS)
endif