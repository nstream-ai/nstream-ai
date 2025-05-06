package utils

// ANSI color codes
const (
	RedColor     = "\033[31m"
	BoldColor    = "\033[1m"
	ResetColor   = "\033[0m"
	BlinkColor   = "\033[5m"
	ReverseColor = "\033[7m"
)

// Common error messages
const (
	ErrAuthRequired     = "authentication required"
	ErrNoClusters       = "no clusters available"
	ErrInvalidChoice    = "invalid choice"
	ErrConfigLoadFailed = "failed to load config"
)

// Common success messages
const (
	SuccessAuth = "Successfully authenticated!"
	SuccessInit = "Successfully initialized!"
)

// Common prompts
const (
	PromptSignIn        = "Please sign in first:"
	PromptSignUp        = "Please sign up first:"
	PromptCloudProvider = "Please select a cloud provider:"
)

// Common table headers
const (
	TableHeaderCluster = "ID\tRegion\tCloud\tBucket\tIdentity"
	TableHeaderBucket  = "ID\tName\tRegion\tCloud\tStatus"
)
