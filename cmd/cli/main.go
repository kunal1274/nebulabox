package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	defaultAPIURL = "http://localhost:8080"
	version       = "0.1.0-alpha"
)

var (
	apiURL  string
	apiToken string
	verbose bool
)

func main() {
	flag.StringVar(&apiURL, "api", defaultAPIURL, "NebulaBox API URL")
	flag.StringVar(&apiToken, "token", "", "Authentication token")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		showUsage()
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	var err error
	switch command {
	case "workspace":
		err = handleWorkspace(commandArgs)
	case "invite":
		err = handleInvite(commandArgs)
	case "join":
		err = handleJoin(commandArgs)
	case "version":
		fmt.Println("NebulaBox CLI", version)
		return
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("NebulaBox CLI - Manage workspaces from the command line")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  nebulabox <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  workspace list                    List all workspaces")
	fmt.Println("  workspace create <name>          Create a new workspace")
	fmt.Println("  workspace get <id>                Get workspace details")
	fmt.Println("  workspace share <id>             Get shareable invite link")
	fmt.Println("  invite create <workspace-id>     Create an invite for a workspace")
	fmt.Println("  invite list <workspace-id>       List invites for a workspace")
	fmt.Println("  invite accept <token>            Accept an invite using token")
	fmt.Println("  join <token>                     Join a workspace using invite token")
	fmt.Println("  version                          Show version")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -api <url>       API server URL (default: http://localhost:8080)")
	fmt.Println("  -token <token>   Authentication token")
	fmt.Println("  -verbose         Verbose output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  nebulabox workspace list")
	fmt.Println("  nebulabox workspace share ws-123")
	fmt.Println("  nebulabox join abc123def456")
}

func handleWorkspace(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("workspace command requires a subcommand")
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "list":
		return listWorkspaces()
	case "create":
		if len(subArgs) < 1 {
			return fmt.Errorf("workspace create requires a name")
		}
		return createWorkspace(subArgs[0])
	case "get":
		if len(subArgs) < 1 {
			return fmt.Errorf("workspace get requires an ID")
		}
		return getWorkspace(subArgs[0])
	case "share":
		if len(subArgs) < 1 {
			return fmt.Errorf("workspace share requires an ID")
		}
		return shareWorkspace(subArgs[0])
	default:
		return fmt.Errorf("unknown workspace command: %s", subcommand)
	}
}

func handleInvite(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("invite command requires a subcommand")
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "create":
		if len(subArgs) < 1 {
			return fmt.Errorf("invite create requires a workspace ID")
		}
		return createInvite(subArgs[0], subArgs[1:]...)
	case "list":
		if len(subArgs) < 1 {
			return fmt.Errorf("invite list requires a workspace ID")
		}
		return listInvites(subArgs[0])
	default:
		return fmt.Errorf("unknown invite command: %s", subcommand)
	}
}

func handleJoin(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("join requires an invite token")
	}
	return joinWorkspace(args[0])
}

func apiRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", apiURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	}
	req.Header.Set("Content-Type", "application/json")

	if verbose {
		fmt.Printf("→ %s %s\n", method, url)
	}

	return http.DefaultClient.Do(req)
}

func listWorkspaces() error {
	resp, err := apiRequest("GET", "/api/shareruntime/workspaces", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to list workspaces: %s", string(body))
	}

	var result struct {
		Workspaces []Workspace `json:"workspaces"`
		Count      int         `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Printf("Found %d workspace(s):\n\n", result.Count)
	for _, ws := range result.Workspaces {
		fmt.Printf("  ID:          %s\n", ws.ID)
		fmt.Printf("  Name:        %s\n", ws.Name)
		fmt.Printf("  Status:      %s\n", ws.Status)
		fmt.Printf("  Members:     %d\n", len(ws.Members))
		fmt.Printf("  Created:     %s\n", ws.CreatedAt)
		fmt.Println()
	}

	return nil
}

func createWorkspace(name string) error {
	// For now, we'll need a container ID - in a real implementation, this would be more interactive
	containerID := "container-placeholder"
	if containerIDEnv := os.Getenv("NEBULABOX_CONTAINER_ID"); containerIDEnv != "" {
		containerID = containerIDEnv
	}

	payload := map[string]interface{}{
		"name":        name,
		"description": fmt.Sprintf("Workspace created via CLI"),
		"containerId": containerID,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := apiRequest("POST", "/api/shareruntime/workspaces", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create workspace: %s", string(body))
	}

	var ws Workspace
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return err
	}

	fmt.Println("✓ Workspace created successfully")
	fmt.Printf("  ID:   %s\n", ws.ID)
	fmt.Printf("  Name: %s\n", ws.Name)
	fmt.Printf("  Status: %s\n", ws.Status)

	return nil
}

func getWorkspace(id string) error {
	resp, err := apiRequest("GET", fmt.Sprintf("/api/shareruntime/workspaces/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get workspace: %s", string(body))
	}

	var ws Workspace
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return err
	}

	fmt.Printf("Workspace: %s\n\n", ws.Name)
	fmt.Printf("  ID:          %s\n", ws.ID)
	fmt.Printf("  Description: %s\n", ws.Description)
	fmt.Printf("  Status:      %s\n", ws.Status)
	fmt.Printf("  Container:   %s\n", ws.ContainerID)
	fmt.Printf("  Members:     %d\n", len(ws.Members))
	fmt.Printf("  Created:     %s\n", ws.CreatedAt)
	fmt.Printf("  Updated:     %s\n", ws.UpdatedAt)
	
	if len(ws.Members) > 0 {
		fmt.Println("\n  Members:")
		for _, member := range ws.Members {
			fmt.Printf("    - %s (%s)\n", member.Username, member.Role)
		}
	}

	return nil
}

func shareWorkspace(id string) error {
	// First try to get an existing invite link
	resp, err := apiRequest("GET", fmt.Sprintf("/api/shareruntime/workspaces/%s/invites", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			Invites []Invite `json:"invites"`
			Count   int      `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result.Count > 0 {
			invite := result.Invites[0]
			link, err := getInviteLink(id, invite.Token)
			if err == nil && link != "" {
				fmt.Println("✓ Shareable invite link:")
				fmt.Println()
				fmt.Printf("  %s\n", link)
				fmt.Println()
				fmt.Println("You can share this link with others to invite them to the workspace.")
				return nil
			}
		}
	}

	// Create a new invite if none exists
	payload := map[string]interface{}{
		"role":    "editor",
		"expires": 0, // No expiration
	}
	jsonData, _ := json.Marshal(payload)
	resp, err = apiRequest("POST", fmt.Sprintf("/api/shareruntime/workspaces/%s/invites", id), strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create invite: %s", string(body))
	}

	var invite Invite
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return err
	}

	link, err := getInviteLink(id, invite.Token)
	if err != nil {
		return err
	}

	fmt.Println("✓ Invite created successfully")
	fmt.Println()
	fmt.Printf("  Share link: %s\n", link)
	fmt.Printf("  Token:      %s\n", invite.Token)
	fmt.Println()
	fmt.Println("Share this link with others to invite them to the workspace.")

	return nil
}

func getInviteLink(workspaceID, token string) (string, error) {
	resp, err := apiRequest("GET", fmt.Sprintf("/api/shareruntime/workspaces/%s/invites/%s/link", workspaceID, token), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get invite link")
	}

	var result struct {
		Link string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Link, nil
}

func createInvite(workspaceID string, roleArgs ...string) error {
	role := "editor"
	if len(roleArgs) > 0 {
		role = roleArgs[0]
	}

	payload := map[string]interface{}{
		"role":    role,
		"expires": 0, // No expiration
	}
	jsonData, _ := json.Marshal(payload)
	resp, err := apiRequest("POST", fmt.Sprintf("/api/shareruntime/workspaces/%s/invites", workspaceID), strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create invite: %s", string(body))
	}

	var invite Invite
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return err
	}

	fmt.Println("✓ Invite created successfully")
	fmt.Printf("  Token:      %s\n", invite.Token)
	fmt.Printf("  Role:       %s\n", invite.Role)
	fmt.Printf("  Expires:    %s\n", invite.ExpiresAt)

	link, err := getInviteLink(workspaceID, invite.Token)
	if err == nil {
		fmt.Printf("  Share link: %s\n", link)
	}

	return nil
}

func listInvites(workspaceID string) error {
	resp, err := apiRequest("GET", fmt.Sprintf("/api/shareruntime/workspaces/%s/invites", workspaceID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to list invites: %s", string(body))
	}

	var result struct {
		Invites []Invite `json:"invites"`
		Count   int      `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Printf("Found %d invite(s) for workspace:\n\n", result.Count)
	for _, invite := range result.Invites {
		fmt.Printf("  Token:      %s\n", invite.Token)
		fmt.Printf("  Role:       %s\n", invite.Role)
		fmt.Printf("  Created:    %s\n", invite.CreatedAt)
		fmt.Printf("  Expires:    %s\n", invite.ExpiresAt)
		fmt.Printf("  Status:     %s\n", invite.Status)
		fmt.Println()
	}

	return nil
}

func joinWorkspace(token string) error {
	resp, err := apiRequest("POST", fmt.Sprintf("/api/shareruntime/invites/%s/accept", token), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to join workspace: %s", string(body))
	}

	var result struct {
		Workspace Workspace `json:"workspace"`
		Message   string    `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Println("✓ Successfully joined workspace")
	fmt.Printf("  Workspace: %s\n", result.Workspace.Name)
	fmt.Printf("  ID:        %s\n", result.Workspace.ID)
	fmt.Printf("  Status:    %s\n", result.Workspace.Status)

	return nil
}

type Workspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ContainerID string    `json:"containerId"`
	Members     []Member  `json:"members"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type Member struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Invite struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	WorkspaceID string `json:"workspaceId"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
	Status    string `json:"status"`
}

