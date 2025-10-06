package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gopher-lab/gopher-mcp-server/data/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var GOPHER_API = os.Getenv("GOPHER_API")

type Input struct {
	Query string `json:"query" jsonschema:"the query to search for"`
}

type Output struct {
	Tweets []client.TwitterPost `json:"tweets" jsonschema:"the tweets of the search"`
	Error  string               `json:"error,omitempty" jsonschema:"the error of the search"`
}

var maxResults = 15

func init() {
	var err error
	maxResults, err = strconv.Atoi(os.Getenv("MAX_RESULTS"))
	if err != nil {
		maxResults = 15
	}
}

func SearchTwitter(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	twitterAPI := client.NewGopherAIClient(GOPHER_API)

	resp, err := twitterAPI.SearchTwitter(input.Query, maxResults)
	if err != nil {
		return nil, Output{Tweets: []client.TwitterPost{}, Error: err.Error()}, nil
	}

	return nil, Output{Tweets: resp.Data}, nil
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "gopher-lab/gopher-mcp-server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "search_twitter", Description: "search the Gopher AI twitter API"}, SearchTwitter)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
