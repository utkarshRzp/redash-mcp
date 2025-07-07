#!/bin/bash

# Test script for Redash MCP Server - Payment Data Query
# This script demonstrates how to test the server for payment data with limit 10

echo "=== Redash MCP Server - Payment Data Test ==="
echo

# Check if environment variables are set
if [ -z "$REDASH_BASE_URL" ]; then
    echo "❌ REDASH_BASE_URL environment variable is not set"
    echo "   Please set it to your Redash instance URL (e.g., https://redash.example.com)"
    exit 1
fi

if [ -z "$REDASH_API_KEY" ]; then
    echo "❌ REDASH_API_KEY environment variable is not set"
    echo "   Please set it to your Redash API key"
    exit 1
fi

echo "✅ Environment variables are set:"
echo "   REDASH_BASE_URL: $REDASH_BASE_URL"
echo "   REDASH_API_KEY: ${REDASH_API_KEY:0:10}..."
echo

# Build the server if not already built
if [ ! -f "./redash-mcp" ]; then
    echo "🔨 Building Redash MCP server..."
    go build -o redash-mcp ./cmd/main
    if [ $? -ne 0 ]; then
        echo "❌ Failed to build server"
        exit 1
    fi
    echo "✅ Server built successfully"
fi

echo "🚀 Starting Redash MCP Server..."
echo "   Use the following tools to query payment data:"
echo
echo "Available tools for payment data:"
echo "1. list_queries - Find queries related to payments"
echo "2. get_query - Get a specific payment query by ID"
echo "3. execute_query - Execute a payment query with limit 10"
echo "4. list_dashboards - Find payment dashboards"
echo
echo "Example MCP requests you can send:"
echo
echo "1. List all queries (to find payment-related queries):"
echo '   {"method": "tools/call", "params": {"name": "list_queries", "arguments": {"page": 1, "page_size": 10, "search": "payment"}}}'
echo
echo "2. Execute a query with limit (replace QUERY_ID with actual payment query ID):"
echo '   {"method": "tools/call", "params": {"name": "execute_query", "arguments": {"queryId": QUERY_ID, "parameters": {"limit": 10}}}}'
echo
echo "3. Get payment dashboards:"
echo '   {"method": "tools/call", "params": {"name": "list_dashboards", "arguments": {"page": 1, "page_size": 10}}}'
echo
echo "Starting server in stdio mode..."
echo "Send MCP requests via stdin to interact with the server."
echo

# Start the server
./redash-mcp stdio 