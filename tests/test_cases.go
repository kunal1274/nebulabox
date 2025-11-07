package tests

// Test cases for NebulaBox

// GetContainerTestCases returns test cases for container operations
func GetContainerTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "container_001",
			Name:        "List Containers",
			Description: "Test listing all containers",
			Category:    "containers",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_001",
					Description: "Call GET /api/containers",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "GET",
						"url":    "/api/containers",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"containers": "array",
				},
			},
			Status: "pending",
		},
		{
			ID:          "container_002",
			Name:        "Create Container",
			Description: "Test creating a new container",
			Category:    "containers",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_002",
					Description: "Call POST /api/containers/run",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/run",
						"body": map[string]interface{}{
							"image": "nginx:latest",
							"name":  "test-container",
						},
					},
					Expected: "201",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 201,
				Response: map[string]interface{}{
					"id":     "string",
					"status": "running",
				},
			},
			Status: "pending",
		},
		{
			ID:          "container_003",
			Name:        "Stop Container",
			Description: "Test stopping a container",
			Category:    "containers",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_003",
					Description: "Call POST /api/containers/{id}/stop",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/stop",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"status": "stopped",
				},
			},
			Status: "pending",
		},
		{
			ID:          "container_004",
			Name:        "Get Container Logs",
			Description: "Test getting container logs",
			Category:    "containers",
			Priority:    "medium",
			Steps: []TestStep{
				{
					ID:          "step_004",
					Description: "Call GET /api/containers/{id}/logs",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "GET",
						"url":    "/api/containers/mock-001/logs",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"logs": "array",
				},
			},
			Status: "pending",
		},
	}
}

// GetExecTestCases returns test cases for exec operations
func GetExecTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "exec_001",
			Name:        "Execute Command",
			Description: "Test executing a command in a container",
			Category:    "exec",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_005",
					Description: "Call POST /api/containers/{id}/exec",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/exec",
						"body": map[string]interface{}{
							"command": []string{"ls", "-la"},
						},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"exitCode": 0,
					"output":   "string",
				},
			},
			Status: "pending",
		},
		{
			ID:          "exec_002",
			Name:        "Execute Command with Streaming",
			Description: "Test executing a command with streaming output",
			Category:    "exec",
			Priority:    "medium",
			Steps: []TestStep{
				{
					ID:          "step_006",
					Description: "Call POST /api/containers/{id}/exec/stream",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/exec/stream",
						"body": map[string]interface{}{
							"command": []string{"env"},
						},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"streaming": true,
				},
			},
			Status: "pending",
		},
		{
			ID:          "exec_003",
			Name:        "Detect Shell",
			Description: "Test detecting available shell in container",
			Category:    "exec",
			Priority:    "medium",
			Steps: []TestStep{
				{
					ID:          "step_007",
					Description: "Call POST /api/containers/{id}/exec/shell",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/exec/shell",
						"body":   map[string]interface{}{},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"shell":     "string",
					"available": true,
				},
			},
			Status: "pending",
		},
	}
}

// GetFileTestCases returns test cases for file operations
func GetFileTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "file_001",
			Name:        "Upload File",
			Description: "Test uploading a file to a container",
			Category:    "files",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_008",
					Description: "Call POST /api/containers/{id}/files/upload",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/files/upload",
						"body": map[string]interface{}{
							"destinationPath": "/app/test.txt",
						},
						"file": "test.txt",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"success":          true,
					"bytesTransferred": "number",
				},
			},
			Status: "pending",
		},
		{
			ID:          "file_002",
			Name:        "Download File",
			Description: "Test downloading a file from a container",
			Category:    "files",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_009",
					Description: "Call POST /api/containers/{id}/files/download",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/files/download",
						"body": map[string]interface{}{
							"sourcePath": "/app/test.txt",
						},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"success":          true,
					"bytesTransferred": "number",
				},
			},
			Status: "pending",
		},
		{
			ID:          "file_003",
			Name:        "List Files",
			Description: "Test listing files in a container directory",
			Category:    "files",
			Priority:    "medium",
			Steps: []TestStep{
				{
					ID:          "step_010",
					Description: "Call POST /api/containers/{id}/files/list",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/files/list",
						"body": map[string]interface{}{
							"path": "/app",
						},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"files": "array",
					"count": "number",
				},
			},
			Status: "pending",
		},
		{
			ID:          "file_004",
			Name:        "Upload File with Streaming",
			Description: "Test uploading a file with streaming progress",
			Category:    "files",
			Priority:    "low",
			Steps: []TestStep{
				{
					ID:          "step_011",
					Description: "Call POST /api/containers/{id}/files/upload/stream",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/containers/mock-001/files/upload/stream",
						"body": map[string]interface{}{
							"destinationPath": "/app/large-file.txt",
						},
						"file": "large-file.txt",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"streaming": true,
				},
			},
			Status: "pending",
		},
	}
}

// GetImageTestCases returns test cases for image operations
func GetImageTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "image_001",
			Name:        "List Images",
			Description: "Test listing all images",
			Category:    "images",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_012",
					Description: "Call GET /api/images",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "GET",
						"url":    "/api/images",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"images": "array",
				},
			},
			Status: "pending",
		},
		{
			ID:          "image_002",
			Name:        "Pull Image",
			Description: "Test pulling an image from registry",
			Category:    "images",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_013",
					Description: "Call POST /api/images/pull",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "POST",
						"url":    "/api/images/pull",
						"body": map[string]interface{}{
							"image": "nginx:latest",
						},
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"success": true,
					"image":   "string",
				},
			},
			Status: "pending",
		},
	}
}

// GetSystemTestCases returns test cases for system operations
func GetSystemTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "system_001",
			Name:        "Get System Stats",
			Description: "Test getting system statistics",
			Category:    "system",
			Priority:    "medium",
			Steps: []TestStep{
				{
					ID:          "step_014",
					Description: "Call GET /api/system/stats",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "GET",
						"url":    "/api/system/stats",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"cpu":    "number",
					"memory": "number",
					"disk":   "number",
				},
			},
			Status: "pending",
		},
		{
			ID:          "system_002",
			Name:        "Health Check",
			Description: "Test API health check endpoint",
			Category:    "system",
			Priority:    "high",
			Steps: []TestStep{
				{
					ID:          "step_015",
					Description: "Call GET /api/health",
					Action:      "api_call",
					Input: map[string]interface{}{
						"method": "GET",
						"url":    "/api/health",
					},
					Expected: "200",
					Status:   "pending",
				},
			},
			Expected: ExpectedResult{
				StatusCode: 200,
				Response: map[string]interface{}{
					"status": "healthy",
				},
			},
			Status: "pending",
		},
	}
}

// GetAllTestCases returns all test cases
func GetAllTestCases() []TestCase {
	var allCases []TestCase
	
	allCases = append(allCases, GetContainerTestCases()...)
	allCases = append(allCases, GetExecTestCases()...)
	allCases = append(allCases, GetFileTestCases()...)
	allCases = append(allCases, GetImageTestCases()...)
	allCases = append(allCases, GetSystemTestCases()...)
	
	return allCases
}
