#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_step() {
  echo -e "${CYAN}==> $1${NC}"
}
print_ok() {
  echo -e "${GREEN}✔ $1${NC}"
}
print_fail() {
  echo -e "${RED}✖ $1${NC}"
}
print_warn() {
  echo -e "${YELLOW}! $1${NC}"
}

print_step "Testing MCP Server in Go (all tools)..."

# Check for built server
if [ ! -f "./mcp-project-helper" ]; then
    print_warn "Server not found. Building..."
    make build-local
fi

print_step "Running unit tests..."
if go test -v; then
    print_ok "Unit tests passed."
else
    print_fail "Unit tests failed!"
    exit 1
fi

# Create temporary test directory
TESTDIR="./test_tmp_dir_$$"
mkdir -p "$TESTDIR"
cd "$TESTDIR"

print_step "Creating test files..."
mkdir subdir

echo "hello world" > file1.txt
echo "foo bar" > subdir/file2.txt

print_step "Generating MCP requests for all tools..."
cat > test_input.jsonrpc << EOF
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "tool-generator", "arguments": {"query": "генерировать инструмент для хайку"}}}
EOF

print_step "Running server with test MCP data..."

cat test_input.jsonrpc | ../mcp-project-helper -transport stdio "$PWD" | tee test_output.log

print_step "Checking results..."

# Check for successful responses
check_ok() {
  local id="$1"
  local desc="$2"
  if grep -q '"id":'$id',"result"' test_output.log; then
    print_ok "$desc"
  else
    print_fail "$desc"
  fi
}

check_ok 2 "tool-generator (генерировать инструмент для хайку)"

# Дополнительные проверки содержимого ответа для новых инструментов
if grep -q '"id":2,"result"' test_output.log && grep -q 'prompt' test_output.log; then
  print_ok "tool-generator вернул prompt"
else
  print_fail "tool-generator не вернул prompt"
fi

cd ..
rm -rf "$TESTDIR"

print_step "Testing completed."
