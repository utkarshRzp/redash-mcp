#!/bin/bash

echo "=== Redash MCP Server - Payment Data Test Example ==="
echo
echo "To test for payment data with limit 10, follow these steps:"
echo
echo "1. Set your environment variables:"
echo "   export REDASH_BASE_URL=\"https://your-redash-instance.com\""
echo "   export REDASH_API_KEY=\"your-api-key\""
echo
echo "2. Start the server:"
echo "   ./redash-mcp stdio"
echo
echo "3. Send MCP requests for payment data (examples):"
echo
echo "   a) Find payment queries:"
echo "   {\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"list_queries\",\"arguments\":{\"page\":1,\"page_size\":10,\"search\":\"payment\"}}}"
echo
echo "   b) Execute payment query with limit 10:"
echo "   {\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"execute_query\",\"arguments\":{\"queryId\":123,\"parameters\":{\"limit\":10}}}}"
echo
echo "   c) List all queries (first 10):"
echo "   {\"jsonrpc\":\"2.0\",\"id\":3,\"method\":\"tools/call\",\"params\":{\"name\":\"list_queries\",\"arguments\":{\"page\":1,\"page_size\":10}}}"
echo
echo "4. The server will respond with JSON containing payment data limited to 10 records"
echo
echo "Available tools for payment data:"
echo "   - list_queries: Find payment-related queries"
echo "   - get_query: Get specific payment query details"
echo "   - execute_query: Execute queries with limit parameter"
echo "   - list_dashboards: Find payment dashboards"
echo "   - list_data_sources: List available data sources"
echo 