package aiops

import (
	"context"
	"strings"
	"time"
)

// ChatOpsEngine provides natural language interface for operations
type ChatOpsEngine struct {
	analytics    *AnalyticsEngine
	scaling      *ScalingAdvisor
	commandHistory []ChatCommand
}

// NewChatOpsEngine creates a new ChatOps engine
func NewChatOpsEngine(analytics *AnalyticsEngine, scaling *ScalingAdvisor) *ChatOpsEngine {
	return &ChatOpsEngine{
		analytics:     analytics,
		scaling:       scaling,
		commandHistory: make([]ChatCommand, 0),
	}
}

// ProcessCommand processes a natural language command
func (ce *ChatOpsEngine) ProcessCommand(ctx context.Context, command string) (*ChatResponse, error) {
	cmd := ChatCommand{
		Command:     command,
		Timestamp:   time.Now(),
		Processed:   true,
	}

	// Parse and execute command
	response := ce.parseAndExecute(ctx, command)
	cmd.Response = response.Message
	cmd.Success = response.Success

	ce.commandHistory = append(ce.commandHistory, cmd)
	
	// Keep only last 100 commands
	if len(ce.commandHistory) > 100 {
		ce.commandHistory = ce.commandHistory[len(ce.commandHistory)-100:]
	}

	return response, nil
}

// parseAndExecute parses natural language and executes commands
func (ce *ChatOpsEngine) parseAndExecute(ctx context.Context, command string) *ChatResponse {
	command = strings.ToLower(strings.TrimSpace(command))

	// Check for scaling commands
	if strings.Contains(command, "scale") || strings.Contains(command, "replicas") {
		return ce.handleScalingCommand(ctx, command)
	}

	// Check for prediction commands
	if strings.Contains(command, "predict") || strings.Contains(command, "forecast") {
		return ce.handlePredictionCommand(ctx, command)
	}

	// Check for status/info commands
	if strings.Contains(command, "status") || strings.Contains(command, "info") || strings.Contains(command, "show") {
		return ce.handleStatusCommand(ctx, command)
	}

	// Check for recommendations
	if strings.Contains(command, "recommend") || strings.Contains(command, "suggest") {
		return ce.handleRecommendationCommand(ctx, command)
	}

	// Default response
	return &ChatResponse{
		Success: false,
		Message: "I didn't understand that command. Try asking about scaling, predictions, status, or recommendations.",
		Suggestions: []string{
			"Show scaling recommendations",
			"Predict resource usage",
			"Show container status",
			"Get optimization suggestions",
		},
	}
}

// handleScalingCommand handles scaling-related commands
func (ce *ChatOpsEngine) handleScalingCommand(ctx context.Context, command string) *ChatResponse {
	// Extract container/deployment ID if mentioned
	// For now, use a simple pattern
	if strings.Contains(command, "up") || strings.Contains(command, "increase") {
		return &ChatResponse{
			Success: true,
			Message: "To scale up, use the scaling recommendations in the dashboard. Would you like me to show current recommendations?",
			Suggestions: []string{"Show scaling recommendations", "Show container status"},
		}
	}

	if strings.Contains(command, "down") || strings.Contains(command, "decrease") {
		return &ChatResponse{
			Success: true,
			Message: "To scale down, use the scaling recommendations in the dashboard. Would you like me to show current recommendations?",
			Suggestions: []string{"Show scaling recommendations"},
		}
	}

	return &ChatResponse{
		Success: true,
		Message: "I can help with scaling. Try: 'show scaling recommendations' or 'predict resource usage for container X'",
		Suggestions: []string{"Show scaling recommendations", "Set scaling policy"},
	}
}

// handlePredictionCommand handles prediction commands
func (ce *ChatOpsEngine) handlePredictionCommand(ctx context.Context, command string) *ChatResponse {
	return &ChatResponse{
		Success: true,
		Message: "Use the Analytics dashboard to view resource usage predictions. I can provide predictions for CPU, memory, and network usage.",
		Suggestions: []string{"Show analytics dashboard", "Get resource predictions"},
	}
}

// handleStatusCommand handles status/info commands
func (ce *ChatOpsEngine) handleStatusCommand(ctx context.Context, command string) *ChatResponse {
	return &ChatResponse{
		Success: true,
		Message: "Check the dashboard for container status, resource usage, and system health. I can provide: container status, resource usage, scaling recommendations.",
		Suggestions: []string{"Show container status", "Show system metrics"},
	}
}

// handleRecommendationCommand handles recommendation requests
func (ce *ChatOpsEngine) handleRecommendationCommand(ctx context.Context, command string) *ChatResponse {
	return &ChatResponse{
		Success: true,
		Message: "Based on current metrics, I recommend checking the Analytics dashboard for optimization suggestions. These include scaling recommendations, resource limit adjustments, and anomaly alerts.",
		Suggestions: []string{"Show recommendations", "Show anomalies", "Show scaling suggestions"},
	}
}

// GetCommandHistory returns recent command history
func (ce *ChatOpsEngine) GetCommandHistory(limit int) []ChatCommand {
	if limit <= 0 || limit > len(ce.commandHistory) {
		limit = len(ce.commandHistory)
	}
	
	start := len(ce.commandHistory) - limit
	if start < 0 {
		start = 0
	}
	
	return ce.commandHistory[start:]
}

// ChatCommand represents a processed chat command
type ChatCommand struct {
	Command   string
	Timestamp time.Time
	Processed bool
	Success   bool
	Response  string
}

// ChatResponse represents a response to a chat command
type ChatResponse struct {
	Success    bool
	Message    string
	Data       map[string]interface{}
	Suggestions []string
}

