# Model Context Protocol (MCP)

A monorepo containing implementations of the Model Context Protocol across multiple programming languages.

## Repository Structure

```
root/
├── .github/             # GitHub workflows, issue templates, etc.
├── docs/                # Documentation for the entire repo
├── scripts/             # Shared scripts (build, deployment, etc.)
├── coralogix/           # Go project
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── go.sum
├── querybook/           # Another Go project
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── go.sum
├── devstack/            # Python project
│   ├── src/
│   ├── tests/
│   ├── requirements.txt
│   └── setup.py
├── something_else/      # TypeScript project
│   ├── src/
│   ├── tests/
│   ├── package.json
│   └── tsconfig.json
└── shared/              # Shared libraries or utilities
    ├── go/
    ├── python/
    └── typescript/
```

> For detailed setup and usage instructions for the Coralogix MCP Server, please refer to the [coralogix/README.md](coralogix/README.md).

## Contributing

Please see the CODEOWNERS file for maintainers of each component.


```
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
        "REDASH_API_KEY": "your-key-here"
      }
    }
    
```