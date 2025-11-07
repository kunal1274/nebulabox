package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	server := &Server{
		teams:           make(map[string]*Team),
		teamMembers:     make(map[string]map[string]*TeamMember),
		userTenants:     make(map[string]string),
		users:           map[string]string{"testuser": "pass"},
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "testuser") })
	router.POST("/teams", server.authMiddleware(), server.createTeam)
	
	body := createTeamRequest{Name: "MyTeam", Description: "Test team"}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/teams", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var team Team
	err := json.Unmarshal(w.Body.Bytes(), &team)
	assert.NoError(t, err)
	assert.Equal(t, "MyTeam", team.Name)
	assert.Equal(t, "testuser", team.CreatedBy)
	
	// Verify creator is admin member
	assert.NotNil(t, server.teamMembers[team.ID])
	assert.Equal(t, "admin", server.teamMembers[team.ID]["testuser"].Role)
}

func TestInviteMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	teamID := "team-1"
	server := &Server{
		teams: map[string]*Team{
			teamID: {ID: teamID, Name: "TestTeam"},
		},
		teamMembers: map[string]map[string]*TeamMember{
			teamID: {"admin": {Username: "admin", Role: "admin"}},
		},
		users: map[string]string{"newuser": "pass"},
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "admin") })
	router.POST("/teams/:id/invite", server.authMiddleware(), server.inviteMember)
	
	body := inviteMemberRequest{Username: "newuser", Role: "editor"}
	jsonBody, _ := json.Marshal(body)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/teams/team-1/invite", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, server.teamMembers[teamID]["newuser"])
	assert.Equal(t, "editor", server.teamMembers[teamID]["newuser"].Role)
}

func TestListTeams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	teamID := "team-1"
	server := &Server{
		teams: map[string]*Team{
			teamID: {ID: teamID, Name: "TestTeam"},
		},
		teamMembers: map[string]map[string]*TeamMember{
			teamID: {"user1": {Username: "user1", Role: "admin"}},
		},
	}
	
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("username", "user1") })
	router.GET("/teams", server.authMiddleware(), server.listTeams)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/teams", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var teams []Team
	err := json.Unmarshal(w.Body.Bytes(), &teams)
	assert.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.Equal(t, "TestTeam", teams[0].Name)
}

