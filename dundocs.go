package dundocs

import "github.com/Sinar/go-dundocs/internal/hansard"

type DUNDocs struct {
	Conf Configuration
}

// CommandMode specifies the operation being executed.
type CommandMode int

// The available commands.
const (
	PLAN CommandMode = iota
	SPLIT
	RESET
)

// Configuration of a Context.
type Configuration struct {
	// DUN Session Label
	DUNSession string

	// Hansard Type
	HansardType hansard.HansardType

	// ./raw + ./data folders are assumed to be relative to this dir
	WorkingDir string

	// Source PDF can be anywhere; maybe make it a Reader to be read direct from S3?
	SourcePDFPath string

	// Command being executed.
	Cmd CommandMode
}
