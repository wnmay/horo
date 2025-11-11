#!/usr/bin/env bash
set -e

# Load environment variables from .env file
ENV_FILE="$(dirname "$0")/.env"

if [[ -f "$ENV_FILE" ]]; then
  echo -e "\033[0;32mLoading environment variables from .env file...\033[0m"

  while IFS='=' read -r name value; do
    # Skip comments and empty lines
    [[ "$name" =~ ^#.*$ || -z "$name" ]] && continue

    # Trim whitespace
    name=$(echo "$name" | xargs)
    value=$(echo "$value" | xargs)

    # Remove quotes if present
    value="${value%\"}"
    value="${value#\"}"
    value="${value%\'}"
    value="${value#\'}"

    export "$name"="$value"
    echo -e "\033[0;36m  Loaded: $name\033[0m"
  done < "$ENV_FILE"

  echo ""
else
  echo -e "\033[0;33mWarning: .env file not found at $ENV_FILE\033[0m"
  echo -e "\033[0;33mCreate a .env file with AUTH_TOKEN=your_token\033[0m"
  echo ""
fi

# Run k6 with the test file passed as argument
TEST_FILE="$1"

if [[ -z "$TEST_FILE" ]]; then
  echo -e "\033[0;31mUsage: ./run-test.sh <test-file>\033[0m"
  echo -e "\033[0;33mExample: ./run-test.sh full-flow-test.js\033[0m"
  exit 1
fi

echo -e "\033[0;32mRunning k6 test: $TEST_FILE\033[0m"
echo -e "\033[0;32m==========================================\033[0m\n"

k6 run "$TEST_FILE"
