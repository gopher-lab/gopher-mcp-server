package client

// TwitterSearchRequest represents the request structure for Twitter search API
type TwitterSearchRequest struct {
	Type      string                 `json:"type"`
	Arguments TwitterSearchArguments `json:"arguments"`
}

// TwitterSearchArguments contains the specific arguments for Twitter search
type TwitterSearchArguments struct {
	Type       string `json:"type"`
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
}

// TwitterSearchResponse represents the response from the Twitter search API
type TwitterSearchResponse struct {
	Status  string        `json:"status"`
	Data    []TwitterPost `json:"data"`
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// TwitterSearchInitResponse represents the initial response with UUID
type TwitterSearchInitResponse struct {
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

// TwitterPost represents a single Twitter post in the response
type TwitterPost struct {
	ID       string      `json:"ID"`
	Content  string      `json:"Content"`
	Metadata interface{} `json:"Metadata"`
	Score    float64     `json:"Score"`
}

// APIError represents an error response from the API
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
