package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewAPICommand creates the API server command
func NewAPICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the NebulaBox API server",
		Long: `Start the NebulaBox API server to serve the web dashboard.
This command starts an HTTP server that provides REST API endpoints
for container and image management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startAPIServer()
		},
	}

	return cmd
}

func startAPIServer() error {
	fmt.Println("🚀 Starting NebulaBox API server...")
	fmt.Println("📡 API will be available at: http://localhost:8080")
	fmt.Println("🌐 Dashboard should connect to: http://localhost:8080/api")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println("")

	// Import and start the API server
	// This will be implemented by calling the API package
	fmt.Println("💡 API server functionality will be implemented next")
	fmt.Println("   Run 'go run ./cmd/api' to start the API server")
	
	return nil
}
