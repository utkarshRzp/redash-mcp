# Razorpay Redash MCP Server

The Redash MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server that provides seamless integration with Redash APIs, enabling advanced query execution, dashboard management, and data visualization capabilities for developers and AI tools.

## Available Tools

The Redash MCP Server provides comprehensive tools for interacting with your Redash instance:

### Read Tools

| Tool                    | Description                                           |
|-------------------------|-------------------------------------------------------|
| `list_queries`          | List queries with pagination and search capabilities |
| `get_query`             | Get specific query details by ID                     |
| `execute_query_fresh`   | Execute query with fresh data (bypasses cache)      |
| `list_data_sources`     | List all available data sources                      |
| `list_dashboards`       | List dashboards with pagination                      |
| `get_dashboard`         | Get specific dashboard details by ID                 |
| `get_visualization`     | Get specific visualization details by ID             |

### Write Tools (Available in write mode)

| Tool                    | Description                                           |
|-------------------------|-------------------------------------------------------|
| `create_query`          | Create a new query                                   |
| `update_query`          | Update an existing query                             |
| `archive_query`         | Archive a query                                      |

## Tool Parameters

### `list_queries`
- `page` (optional, number): Page number (default: 1)
- `pageSize` (optional, number): Number of queries per page (default: 25, max: 100)
- `search` (optional, string): Search term to filter queries

### `get_query`
- `queryId` (**required**, number): The ID of the query to retrieve

### `execute_query_fresh`
- `queryId` (**required**, number): The ID of the query to execute
- `parameters` (optional, object): Parameters for parameterized queries

### `list_data_sources`
No parameters required.

### `list_dashboards`
- `page` (optional, number): Page number (default: 1)
- `pageSize` (optional, number): Number of dashboards per page (default: 25, max: 100)

### `get_dashboard`
- `dashboardId` (**required**, number): The ID of the dashboard to retrieve

### `get_visualization`
- `visualizationId` (**required**, number): The ID of the visualization to retrieve

### `create_query` (Write mode only)
- `name` (**required**, string): Name of the query
- `data_source_id` (**required**, number): ID of the data source
- `query` (**required**, string): SQL query text
- `description` (optional, string): Description of the query
- `options` (optional, object): Query options
- `schedule` (optional, object): Query schedule configuration
- `tags` (optional, array): Query tags

### `update_query` (Write mode only)
- `queryId` (**required**, number): ID of the query to update
- `name` (optional, string): New name for the query
- `data_source_id` (optional, number): New data source ID
- `query` (optional, string): New SQL query text
- `description` (optional, string): New description
- `options` (optional, object): New query options
- `schedule` (optional, object): New schedule configuration
- `tags` (optional, array): New tags
- `is_archived` (optional, boolean): Archive status
- `is_draft` (optional, boolean): Draft status

### `archive_query` (Write mode only)
- `queryId` (**required**, number): ID of the query to archive

## Use Cases

- **Data Analysis**: Execute queries and analyze results through AI assistance
- **Dashboard Management**: Create, update, and manage Redash dashboards
- **Query Development**: Develop and test SQL queries with AI guidance
- **Data Visualization**: Create and manage visualizations for your data
- **Workflow Automation**: Automate your data workflow using the Redash MCP Server

## Setup

### Prerequisites
- Docker
- Go 1.23+ (for building from source)
- Access to a Redash instance
- Redash API key

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/razorpay/redash-mcp.git
cd redash-mcp/redash/

# Build the Docker image
docker build -t redash-mcp:latest .

# Run the server
docker run --rm -i \
  -e REDASH_URL="https://redash.razorpay.com" \
  -e REDASH_API_KEY="your-api-key" \
  redash-mcp:latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/razorpay/redash-mcp.git
cd mcp/redash-mcp/redash/

# Build the binary
go build -o redash-mcp ./cmd/main/

# Run the server
./redash-mcp stdio \
  --redash-url "https://redash.razorpay.com" \
  --redash-api-key "your-api-key"
```

## Usage with Claude Desktop or Cursor

Add the following to your `claude_desktop_config.json` or `mcp.json` file:

### Docker Configuration

```json
{
  "mcpServers": {
    "razorpay-redash-mcp-server": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e",
        "REDASH_URL",
        "-e",
        "REDASH_API_KEY",
        "redash-mcp:latest"
      ],
      "env": {
        "REDASH_URL": "https://redash.razorpay.com",
        "REDASH_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Configuration

### Required Configuration

- `REDASH_URL`: `https://redash.razorpay.com`
- `REDASH_API_KEY`: Your Redash API key
