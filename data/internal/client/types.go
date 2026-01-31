package client

// TwitterRequest represents the request structure for Twitter API (job type "twitter")
type TwitterRequest struct {
	Type      string            `json:"type"`
	Arguments TwitterArguments  `json:"arguments"`
}

// TwitterArguments contains the common parameters for Twitter API operations.
// See docs: type (required), query, max_results (default 1000), count, next_cursor, start_time, end_time.
type TwitterArguments struct {
	Type       string `json:"type"`                   // Operation: searchbyquery, getbyid, getreplies, getretweeters, gettweets, getmedia, searchbyprofile, getprofilebyid, getfollowers, getfollowing, gettrends, getspace, searchbyfullarchive
	Query      string `json:"query,omitempty"`         // Username (@user), tweet ID, search terms, or user ID depending on operation
	MaxResults int    `json:"max_results,omitempty"`   // Max 1000, default 1000
	Count      int    `json:"count,omitempty"`         // Number of results per request (max 1000)
	NextCursor string `json:"next_cursor,omitempty"`   // Pagination cursor
	StartTime  string `json:"start_time,omitempty"`    // ISO 8601 timestamp
	EndTime    string `json:"end_time,omitempty"`      // ISO 8601 timestamp
}

// TwitterSearchRequest is an alias for TwitterRequest (backward compatible).
type TwitterSearchRequest = TwitterRequest

// TwitterSearchArguments is an alias for TwitterArguments (backward compatible).
type TwitterSearchArguments = TwitterArguments

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
