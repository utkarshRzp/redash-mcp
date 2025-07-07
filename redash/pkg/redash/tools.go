package redash

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/razorpay/mcp/redash/pkg/mcpgo"
	"github.com/razorpay/mcp/redash/pkg/toolsets"
)

// NewRedashToolset creates a new Redash toolset
func NewRedashToolset(log *slog.Logger, readOnly bool) *toolsets.Toolset {
	toolset := toolsets.NewToolset("redash", "Redash query and dashboard management tools")

	// Create Redash client
	client, err := NewRedashClient()
	if err != nil {
		log.Error("failed to create redash client", "error", err)
		return toolset
	}

	// Read tools
	toolset.AddReadTools(
		createGetQueryTool(client, log),
		createListQueriesTools(client, log),
		createListDataSourcesTool(client, log),
		createExecuteQueryFreshTool(client, log),
		createListDashboardsTool(client, log),
		createGetDashboardTool(client, log),
		createGetVisualizationTool(client, log),
	)

	// Write tools (only if not read-only)
	if !readOnly {
		toolset.AddWriteTools(
			createCreateQueryTool(client, log),
			createUpdateQueryTool(client, log),
			createArchiveQueryTool(client, log),
		)
	}

	return toolset
}

// createGetQueryTool creates the get_query tool
func createGetQueryTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"get_query",
		"Get a specific query by ID",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("queryId", mcpgo.Required(), mcpgo.Description("The ID of the query to retrieve")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			queryIdFloat, ok := request.Arguments["queryId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("queryId must be a number"), nil
			}
			queryId := int(queryIdFloat)

			query, err := client.GetQuery(queryId)
			if err != nil {
				log.Error("error getting query", "queryId", queryId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error getting query %d: %s", queryId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(query)
		},
	)
}

// createListQueriesTools creates the list_queries tool
func createListQueriesTools(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"list_queries",
		"List queries with optional pagination and search",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("page", mcpgo.DefaultValue(1), mcpgo.Description("Page number (default: 1)")),
			mcpgo.WithNumber("pageSize", mcpgo.DefaultValue(25), mcpgo.Description("Number of queries per page (default: 25)")),
			mcpgo.WithString("search", mcpgo.Description("Search term to filter queries")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			page := 1
			pageSize := 25
			var search string

			if p, ok := request.Arguments["page"].(float64); ok {
				page = int(p)
			}
			if ps, ok := request.Arguments["pageSize"].(float64); ok {
				pageSize = int(ps)
			}
			if s, ok := request.Arguments["search"].(string); ok {
				search = s
			}

			result, err := client.GetQueries(page, pageSize, search)
			if err != nil {
				log.Error("error listing queries", "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error listing queries: %s", err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}

// createListDataSourcesTool creates the list_data_sources tool
func createListDataSourcesTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"list_data_sources",
		"List all available data sources",
		[]mcpgo.ToolParameter{},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			dataSources, err := client.GetDataSources()
			if err != nil {
				log.Error("error listing data sources", "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error listing data sources: %s", err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(dataSources)
		},
	)
}

// createExecuteQueryFreshTool creates the execute_query_fresh tool
func createExecuteQueryFreshTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"execute_query_fresh",
		"Execute a query with fresh data (no cache)",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("queryId", mcpgo.Required(), mcpgo.Description("The ID of the query to execute")),
			mcpgo.WithObject("parameters", mcpgo.Description("Parameters for the query")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			queryIdFloat, ok := request.Arguments["queryId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("queryId must be a number"), nil
			}
			queryId := int(queryIdFloat)

			var parameters map[string]interface{}
			if p, ok := request.Arguments["parameters"].(map[string]interface{}); ok {
				parameters = p
			}

			result, err := client.ExecuteQueryFresh(queryId, parameters)
			if err != nil {
				log.Error("error executing query fresh", "queryId", queryId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error executing query %d fresh: %s", queryId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}

// createListDashboardsTool creates the list_dashboards tool
func createListDashboardsTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"list_dashboards",
		"List dashboards with optional pagination",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("page", mcpgo.DefaultValue(1), mcpgo.Description("Page number (default: 1)")),
			mcpgo.WithNumber("pageSize", mcpgo.DefaultValue(25), mcpgo.Description("Number of dashboards per page (default: 25)")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			page := 1
			pageSize := 25

			if p, ok := request.Arguments["page"].(float64); ok {
				page = int(p)
			}
			if ps, ok := request.Arguments["pageSize"].(float64); ok {
				pageSize = int(ps)
			}

			result, err := client.GetDashboards(page, pageSize)
			if err != nil {
				log.Error("error listing dashboards", "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error listing dashboards: %s", err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}

// createGetDashboardTool creates the get_dashboard tool
func createGetDashboardTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"get_dashboard",
		"Get a specific dashboard by ID",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("dashboardId", mcpgo.Required(), mcpgo.Description("The ID of the dashboard to retrieve")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			dashboardIdFloat, ok := request.Arguments["dashboardId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("dashboardId must be a number"), nil
			}
			dashboardId := int(dashboardIdFloat)

			dashboard, err := client.GetDashboard(dashboardId)
			if err != nil {
				log.Error("error getting dashboard", "dashboardId", dashboardId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error getting dashboard %d: %s", dashboardId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(dashboard)
		},
	)
}

// createGetVisualizationTool creates the get_visualization tool
func createGetVisualizationTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"get_visualization",
		"Get a specific visualization by ID",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("visualizationId", mcpgo.Required(), mcpgo.Description("The ID of the visualization to retrieve")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			visualizationIdFloat, ok := request.Arguments["visualizationId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("visualizationId must be a number"), nil
			}
			visualizationId := int(visualizationIdFloat)

			visualization, err := client.GetVisualization(visualizationId)
			if err != nil {
				log.Error("error getting visualization", "visualizationId", visualizationId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error getting visualization %d: %s", visualizationId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(visualization)
		},
	)
}

// createCreateQueryTool creates the create_query tool
func createCreateQueryTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"create_query",
		"Create a new query",
		[]mcpgo.ToolParameter{
			mcpgo.WithString("name", mcpgo.Required(), mcpgo.Description("Name of the query")),
			mcpgo.WithNumber("data_source_id", mcpgo.Required(), mcpgo.Description("ID of the data source")),
			mcpgo.WithString("query", mcpgo.Required(), mcpgo.Description("SQL query text")),
			mcpgo.WithString("description", mcpgo.Description("Description of the query")),
			mcpgo.WithObject("options", mcpgo.Description("Query options")),
			mcpgo.WithObject("schedule", mcpgo.Description("Query schedule")),
			mcpgo.WithArray("tags", mcpgo.Description("Query tags")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			name, ok := request.Arguments["name"].(string)
			if !ok {
				return mcpgo.NewToolResultError("name must be a string"), nil
			}

			dataSourceIdFloat, ok := request.Arguments["data_source_id"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("data_source_id must be a number"), nil
			}
			dataSourceId := int(dataSourceIdFloat)

			query, ok := request.Arguments["query"].(string)
			if !ok {
				return mcpgo.NewToolResultError("query must be a string"), nil
			}

			queryData := CreateQueryRequest{
				Name:         name,
				DataSourceId: dataSourceId,
				Query:        query,
			}

			if desc, ok := request.Arguments["description"].(string); ok {
				queryData.Description = &desc
			}

			if options, ok := request.Arguments["options"].(map[string]interface{}); ok {
				queryData.Options = options
			}

			if schedule, ok := request.Arguments["schedule"].(map[string]interface{}); ok {
				queryData.Schedule = schedule
			}

			if tagsInterface, ok := request.Arguments["tags"].([]interface{}); ok {
				tags := make([]string, len(tagsInterface))
				for i, tag := range tagsInterface {
					if tagStr, ok := tag.(string); ok {
						tags[i] = tagStr
					}
				}
				queryData.Tags = tags
			}

			result, err := client.CreateQuery(queryData)
			if err != nil {
				log.Error("error creating query", "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error creating query: %s", err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}

// createUpdateQueryTool creates the update_query tool
func createUpdateQueryTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"update_query",
		"Update an existing query",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("queryId", mcpgo.Required(), mcpgo.Description("ID of the query to update")),
			mcpgo.WithString("name", mcpgo.Description("Name of the query")),
			mcpgo.WithNumber("data_source_id", mcpgo.Description("ID of the data source")),
			mcpgo.WithString("query", mcpgo.Description("SQL query text")),
			mcpgo.WithString("description", mcpgo.Description("Description of the query")),
			mcpgo.WithObject("options", mcpgo.Description("Query options")),
			mcpgo.WithObject("schedule", mcpgo.Description("Query schedule")),
			mcpgo.WithArray("tags", mcpgo.Description("Query tags")),
			mcpgo.WithBoolean("is_archived", mcpgo.Description("Whether the query is archived")),
			mcpgo.WithBoolean("is_draft", mcpgo.Description("Whether the query is a draft")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			queryIdFloat, ok := request.Arguments["queryId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("queryId must be a number"), nil
			}
			queryId := int(queryIdFloat)

			updateData := UpdateQueryRequest{}

			if name, ok := request.Arguments["name"].(string); ok {
				updateData.Name = &name
			}

			if dataSourceIdFloat, ok := request.Arguments["data_source_id"].(float64); ok {
				dataSourceId := int(dataSourceIdFloat)
				updateData.DataSourceId = &dataSourceId
			}

			if query, ok := request.Arguments["query"].(string); ok {
				updateData.Query = &query
			}

			if desc, ok := request.Arguments["description"].(string); ok {
				updateData.Description = &desc
			}

			if options, ok := request.Arguments["options"].(map[string]interface{}); ok {
				updateData.Options = options
			}

			if schedule, ok := request.Arguments["schedule"].(map[string]interface{}); ok {
				updateData.Schedule = schedule
			}

			if tagsInterface, ok := request.Arguments["tags"].([]interface{}); ok {
				tags := make([]string, len(tagsInterface))
				for i, tag := range tagsInterface {
					if tagStr, ok := tag.(string); ok {
						tags[i] = tagStr
					}
				}
				updateData.Tags = tags
			}

			if isArchived, ok := request.Arguments["is_archived"].(bool); ok {
				updateData.IsArchived = &isArchived
			}

			if isDraft, ok := request.Arguments["is_draft"].(bool); ok {
				updateData.IsDraft = &isDraft
			}

			result, err := client.UpdateQuery(queryId, updateData)
			if err != nil {
				log.Error("error updating query", "queryId", queryId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error updating query %d: %s", queryId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}

// createArchiveQueryTool creates the archive_query tool
func createArchiveQueryTool(client *RedashClient, log *slog.Logger) mcpgo.Tool {
	return mcpgo.NewTool(
		"archive_query",
		"Archive a query",
		[]mcpgo.ToolParameter{
			mcpgo.WithNumber("queryId", mcpgo.Required(), mcpgo.Description("ID of the query to archive")),
		},
		func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
			queryIdFloat, ok := request.Arguments["queryId"].(float64)
			if !ok {
				return mcpgo.NewToolResultError("queryId must be a number"), nil
			}
			queryId := int(queryIdFloat)

			result, err := client.ArchiveQuery(queryId)
			if err != nil {
				log.Error("error archiving query", "queryId", queryId, "error", err)
				return mcpgo.NewToolResultError(fmt.Sprintf("Error archiving query %d: %s", queryId, err.Error())), nil
			}

			return mcpgo.NewToolResultJSON(result)
		},
	)
}
