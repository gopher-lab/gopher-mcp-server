# Gopher MCP Server

A Model Context Protocol (MCP) server that provides access to the Gopher AI API. This server allows MCP clients to search data through the Gopher AI platform.

## Features

- **Twitter Search**: Search Twitter posts using natural language queries
- **MCP Integration**: Fully compatible with Model Context Protocol clients
- **Configurable Results**: Adjustable maximum number of search results
- **Error Handling**: Comprehensive error handling and reporting
- **Async Processing**: Handles asynchronous API responses with polling

## Prerequisites

- Go 1.24.7 or later
- Gopher AI API key
- MCP-compatible client (e.g., Claude Desktop, Cline, etc.)

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/gopher-lab/gopher-mcp-server.git
cd gopher-mcp-server
```

2. Build the server:
```bash
go build -o gopher-mcp-server ./data
```

### Using Docker

```bash
docker build -t gopher-mcp-server .
```

## Configuration

### Environment Variables

- `GOPHER_API`: **Required** - Your Gopher AI API key
- `MAX_RESULTS`: **Optional** - Maximum number of search results to return (default: 15)

### Example Configuration

```bash
export GOPHER_API="your-gopher-api-key-here"
export MAX_RESULTS="20"
```

## Usage

### Running the Server

```bash
# Using the binary
./gopher-mcp-server

# Using Docker
docker run -e GOPHER_API="your-api-key" gopher-mcp-server
```

### MCP Client Configuration

Add this server to your MCP client configuration:

#### Claude Desktop (claude_desktop_config.json)

```json
{
  "mcpServers": {
    "gopher-twitter": {
      "command": "./gopher-mcp-server",
      "env": {
        "GOPHER_API": "your-gopher-api-key-here",
        "MAX_RESULTS": "15"
      }
    }
  }
}
```

#### Cline Configuration

```json
{
  "mcpServers": {
    "gopher-twitter": {
      "command": "docker",
      "args": ["run", "--rm", "-e", "GOPHER_API=your-api-key", "gopher-mcp-server"],
      "env": {
        "GOPHER_API": "your-gopher-api-key-here"
      }
    }
  }
}
```

## Available Tools

### search_twitter

Searches Twitter posts using the Gopher AI API.

**Parameters:**
- `query` (string): The search query to execute

**Returns:**
- `tweets` (array): Array of Twitter posts matching the search query
- `error` (string, optional): Error message if the search fails

**Example Usage:**

```json
{
  "name": "search_twitter",
  "arguments": {
    "query": "artificial intelligence trends 2024"
  }
}
```

**Response Example:**

```json
{
  "tweets": [
    {
      "ID": "tweet_123",
      "Content": "Exciting developments in AI this year...",
      "Metadata": {...},
      "Score": 0.95
    }
  ],
  "error": null
}
```

## API Details

### Twitter Post Structure

Each Twitter post in the response contains:

- `ID`: Unique identifier for the post
- `Content`: The text content of the tweet
- `Metadata`: Additional metadata about the post
- `Score`: Relevance score for the search query

### Error Handling

The server handles various error scenarios:

- **API Errors**: Invalid API keys, rate limiting, service unavailable
- **Network Errors**: Connection timeouts, network failures
- **Processing Errors**: Search timeout, invalid responses

Error messages are returned in the `error` field of the response.

## Development


### Building

```bash
# Build the server
make build

# Run tests
make test

# Build Docker image
make docker-build
```

### Dependencies

- `github.com/modelcontextprotocol/go-sdk`: MCP SDK for Go
- `github.com/tmc/langchaingo`: Language chain utilities

## Troubleshooting

### Common Issues

1. **API Key Not Set**
   - Ensure `GOPHER_API` environment variable is set
   - Verify the API key is valid and has proper permissions

2. **Search Timeout**
   - The server polls for results up to 30 times (60 seconds total)
   - Check if the Gopher AI service is experiencing delays

3. **No Results Returned**
   - Verify the search query is valid
   - Check if there are any posts matching your criteria
   - Ensure the API key has access to Twitter data

4. **MCP Client Connection Issues**
   - Verify the server binary is executable
   - Check that the MCP client configuration is correct
   - Ensure the server is running and accessible

### Debugging

Enable verbose logging by setting the log level:

```bash
export LOG_LEVEL=debug
./gopher-mcp-server
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

## Support

For issues and questions:

- Create an issue on GitHub
- Check the Gopher AI documentation for API-related questions
- Review MCP documentation for protocol-related issues
