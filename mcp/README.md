# MCP Server

## Purpose
This MCP Server provides tools to communicate with an instance of a service dependency API.

## Tools Available
### Hello World
Provides a basic hello world example to understand how to use the mcp server.


## Prompts
### Hello World
Points a user to use the hello world tool.

## Testing

To test the mcp server locally, use the `ModelContextProtocolInspector`.
This tool requires node and can be run with the following command from the root of the project:

```bash
npx @modelcontextprotocol/inspector
```

Once the mcp inspector is running, navigate to the localhost url that is displayed in the terminal.
You will need to input the run command as `go run` and the optional arguments as `./mcp/cmd/service-dependency-mcp`.

**Note: You will need to set the `API_URL` environment variable to the url of the service dependency api.**

To facilitate testing, a docker compose file is provided in the root of the project. You can run the docker compose file with the following command:

```bash
docker compose -f service.compose.yml up -d
```

Then you can run the `Sample Service Map` folder in the Bruno collection. This will populate the docker compose database with sample data via the api.

## References
- [BruceTraining MCP Example](https://github.com/BruceTraining/MCP-Oreilly-2)
- [MCP Resources](https://modelcontextprotocol.io/specification/2025-06-18/server/resources)
- [Go MCP SDK Example](https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/everything/main.go)