package shareruntime

import (
	"fmt"
	"sync"
	"time"
)

// CRDTType represents the type of CRDT operation
type CRDTType string

const (
	CRDTTypeORSet  CRDTType = "orset"  // Observe-Remove Set (for member lists)
	CRDTTypeLWWReg CRDTType = "lwwreg" // Last-Write-Wins Register (for settings)
	CRDTTypeCounter CRDTType = "counter" // Increment-only counter
	CRDTTypeMap    CRDTType = "map"    // CRDT Map (for nested structures)
)

// CRDTOperation represents a CRDT operation
type CRDTOperation struct {
	ID          string                 `json:"id"`
	Type        CRDTType               `json:"type"`
	WorkspaceID string                 `json:"workspaceId"`
	ResourceID  string                 `json:"resourceId"`
	ResourceType string                `json:"resourceType"` // member, setting, container, etc.
	Operation   string                 `json:"operation"`   // add, remove, update, increment
	Key         string                 `json:"key,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"userId"`
	VectorClock map[string]int64       `json:"vectorClock"` // Vector clock for ordering
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// VectorClock manages vector clocks for CRDT operations
type VectorClock struct {
	clocks map[string]int64 // nodeID -> logical clock
	mu     sync.RWMutex
}

// NewVectorClock creates a new vector clock
func NewVectorClock() *VectorClock {
	return &VectorClock{
		clocks: make(map[string]int64),
	}
}

// Tick increments the clock for a node
func (vc *VectorClock) Tick(nodeID string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.clocks[nodeID]++
}

// Get returns the clock value for a node
func (vc *VectorClock) Get(nodeID string) int64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.clocks[nodeID]
}

// GetClocks returns a copy of all clocks
func (vc *VectorClock) GetClocks() map[string]int64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	result := make(map[string]int64)
	for k, v := range vc.clocks {
		result[k] = v
	}
	return result
}

// Merge merges another vector clock into this one (takes maximum)
func (vc *VectorClock) Merge(other map[string]int64) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	for nodeID, clock := range other {
		if current, exists := vc.clocks[nodeID]; !exists || clock > current {
			vc.clocks[nodeID] = clock
		}
	}
}

// HappensBefore checks if this clock happens before another
func (vc *VectorClock) HappensBefore(other map[string]int64) bool {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	
	strictlyLess := false
	for nodeID, clock := range other {
		if vc.clocks[nodeID] > clock {
			return false
		}
		if vc.clocks[nodeID] < clock {
			strictlyLess = true
		}
	}
	return strictlyLess
}

// ORSet is an Observe-Remove Set CRDT
type ORSet struct {
	addSet    map[string]map[string]bool // element -> nodeID -> true (add tags)
	removeSet map[string]map[string]bool // element -> nodeID -> true (remove tags)
	mu        sync.RWMutex
}

// NewORSet creates a new ORSet
func NewORSet() *ORSet {
	return &ORSet{
		addSet:    make(map[string]map[string]bool),
		removeSet: make(map[string]map[string]bool),
	}
}

// Add adds an element to the set
func (s *ORSet) Add(element string, nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.addSet[element] == nil {
		s.addSet[element] = make(map[string]bool)
	}
	s.addSet[element][nodeID] = true
}

// Remove removes an element from the set
func (s *ORSet) Remove(element string, nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.removeSet[element] == nil {
		s.removeSet[element] = make(map[string]bool)
	}
	s.removeSet[element][nodeID] = true
}

// Contains checks if an element is in the set
func (s *ORSet) Contains(element string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	addTags := s.addSet[element]
	if len(addTags) == 0 {
		return false
	}
	
	removeTags := s.removeSet[element]
	for tag := range addTags {
		if !removeTags[tag] {
			return true
		}
	}
	return false
}

// Elements returns all elements in the set
func (s *ORSet) Elements() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := []string{}
	for element := range s.addSet {
		if s.Contains(element) {
			result = append(result, element)
		}
	}
	return result
}

// Merge merges another ORSet into this one
func (s *ORSet) Merge(other *ORSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	other.mu.RLock()
	defer other.mu.RUnlock()
	
	// Merge add sets
	for element, tags := range other.addSet {
		if s.addSet[element] == nil {
			s.addSet[element] = make(map[string]bool)
		}
		for tag := range tags {
			s.addSet[element][tag] = true
		}
	}
	
	// Merge remove sets
	for element, tags := range other.removeSet {
		if s.removeSet[element] == nil {
			s.removeSet[element] = make(map[string]bool)
		}
		for tag := range tags {
			s.removeSet[element][tag] = true
		}
	}
}

// LWWRegister is a Last-Write-Wins Register CRDT
type LWWRegister struct {
	value     interface{}
	timestamp time.Time
	nodeID    string
	mu        sync.RWMutex
}

// NewLWWRegister creates a new LWW Register
func NewLWWRegister() *LWWRegister {
	return &LWWRegister{
		timestamp: time.Time{},
	}
}

// Set sets the value with a timestamp and node ID
func (r *LWWRegister) Set(value interface{}, timestamp time.Time, nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Last-write-wins: compare timestamp, then node ID if equal
	if timestamp.After(r.timestamp) || (timestamp.Equal(r.timestamp) && nodeID > r.nodeID) {
		r.value = value
		r.timestamp = timestamp
		r.nodeID = nodeID
	}
}

// Get returns the current value
func (r *LWWRegister) Get() interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.value
}

// Merge merges another LWW Register into this one
func (r *LWWRegister) Merge(other *LWWRegister) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	other.mu.RLock()
	defer other.mu.RUnlock()
	
	r.Set(other.value, other.timestamp, other.nodeID)
}

// ConflictResolutionManager manages CRDT operations and conflict resolution
type ConflictResolutionManager struct {
	operations    map[string]*CRDTOperation // operationID -> operation
	vectors       map[string]*VectorClock    // workspaceID -> vector clock
	orsets        map[string]*ORSet          // resourceID -> ORSet
	registers     map[string]*LWWRegister    // resourceID -> LWWRegister
	nodeID        string                     // This node's ID
	mu            sync.RWMutex
}

// NewConflictResolutionManager creates a new conflict resolution manager
func NewConflictResolutionManager(nodeID string) *ConflictResolutionManager {
	return &ConflictResolutionManager{
		operations: make(map[string]*CRDTOperation),
		vectors:    make(map[string]*VectorClock),
		orsets:     make(map[string]*ORSet),
		registers:  make(map[string]*LWWRegister),
		nodeID:     nodeID,
	}
}

// GetVectorClock gets or creates a vector clock for a workspace
func (crm *ConflictResolutionManager) GetVectorClock(workspaceID string) *VectorClock {
	crm.mu.Lock()
	defer crm.mu.Unlock()
	
	if crm.vectors[workspaceID] == nil {
		crm.vectors[workspaceID] = NewVectorClock()
	}
	return crm.vectors[workspaceID]
}

// RecordOperation records a CRDT operation
func (crm *ConflictResolutionManager) RecordOperation(op *CRDTOperation) error {
	crm.mu.Lock()
	defer crm.mu.Unlock()
	
	// Generate operation ID if not provided
	if op.ID == "" {
		op.ID = fmt.Sprintf("op-%s-%d", crm.nodeID, time.Now().UnixNano())
	}
	
	// Initialize vector clock if not provided
	if op.VectorClock == nil {
		vc := crm.GetVectorClock(op.WorkspaceID)
		vc.Tick(crm.nodeID)
		op.VectorClock = vc.GetClocks()
	}
	
	op.Timestamp = time.Now()
	crm.operations[op.ID] = op
	
	// Apply operation to CRDTs
	switch op.Type {
	case CRDTTypeORSet:
		crm.applyORSetOperation(op)
	case CRDTTypeLWWReg:
		crm.applyLWWRegOperation(op)
	}
	
	return nil
}

// applyORSetOperation applies an ORSet operation
func (crm *ConflictResolutionManager) applyORSetOperation(op *CRDTOperation) {
	if crm.orsets[op.ResourceID] == nil {
		crm.orsets[op.ResourceID] = NewORSet()
	}
	
	switch op.Operation {
	case "add":
		if val, ok := op.Value.(string); ok {
			crm.orsets[op.ResourceID].Add(val, crm.nodeID)
		}
	case "remove":
		if val, ok := op.Value.(string); ok {
			crm.orsets[op.ResourceID].Remove(val, crm.nodeID)
		}
	}
}

// applyLWWRegOperation applies an LWW Register operation
func (crm *ConflictResolutionManager) applyLWWRegOperation(op *CRDTOperation) {
	if crm.registers[op.ResourceID] == nil {
		crm.registers[op.ResourceID] = NewLWWRegister()
	}
	
	if op.Operation == "update" {
		crm.registers[op.ResourceID].Set(op.Value, op.Timestamp, crm.nodeID)
	}
}

// GetOperationsSince gets operations since a given timestamp or operation ID
func (crm *ConflictResolutionManager) GetOperationsSince(workspaceID string, since time.Time) []*CRDTOperation {
	crm.mu.RLock()
	defer crm.mu.RUnlock()
	
	var result []*CRDTOperation
	for _, op := range crm.operations {
		if op.WorkspaceID == workspaceID && op.Timestamp.After(since) {
			result = append(result, op)
		}
	}
	return result
}

// DetectConflicts detects conflicts between operations
func (crm *ConflictResolutionManager) DetectConflicts(ops []*CRDTOperation) []Conflict {
	crm.mu.RLock()
	defer crm.mu.RUnlock()
	
	var conflicts []Conflict
	
	// Check for concurrent modifications to the same resource
	resourceOps := make(map[string][]*CRDTOperation)
	for _, op := range ops {
		key := fmt.Sprintf("%s:%s", op.ResourceID, op.ResourceType)
		resourceOps[key] = append(resourceOps[key], op)
	}
	
	for _, ops := range resourceOps {
		if len(ops) < 2 {
			continue
		}
		
		// Check if operations are concurrent (not causally ordered)
		for i := 0; i < len(ops); i++ {
			for j := i + 1; j < len(ops); j++ {
				op1 := ops[i]
				op2 := ops[j]
				
				// Check if operations are concurrent
				if !crm.isCausallyOrdered(op1, op2) && !crm.isCausallyOrdered(op2, op1) {
					conflicts = append(conflicts, Conflict{
						ResourceID:   op1.ResourceID,
						ResourceType: op1.ResourceType,
						Operations:   []*CRDTOperation{op1, op2},
						Type:         ConflictTypeConcurrent,
					})
				}
			}
		}
	}
	
	return conflicts
}

// isCausallyOrdered checks if op1 causally precedes op2
func (crm *ConflictResolutionManager) isCausallyOrdered(op1, op2 *CRDTOperation) bool {
	vc1 := NewVectorClock()
	vc1.clocks = op1.VectorClock
	
	return vc1.HappensBefore(op2.VectorClock)
}

// ResolveConflict resolves a conflict using CRDT semantics
func (crm *ConflictResolutionManager) ResolveConflict(conflict Conflict) (*CRDTOperation, error) {
	crm.mu.Lock()
	defer crm.mu.Unlock()
	
	if len(conflict.Operations) < 2 {
		return nil, fmt.Errorf("conflict must have at least 2 operations")
	}
	
	op1 := conflict.Operations[0]
	op2 := conflict.Operations[1]
	
	// Apply CRDT merge semantics based on type
	switch op1.Type {
	case CRDTTypeORSet:
		return crm.resolveORSetConflict(op1, op2)
	case CRDTTypeLWWReg:
		return crm.resolveLWWRegConflict(op1, op2)
	default:
		// Default: last-write-wins
		if op1.Timestamp.After(op2.Timestamp) {
			return op1, nil
		}
		return op2, nil
	}
}

// resolveORSetConflict resolves an ORSet conflict
func (crm *ConflictResolutionManager) resolveORSetConflict(op1, op2 *CRDTOperation) (*CRDTOperation, error) {
	// ORSet operations commute, so both can be applied
	// Return a merge operation
	return &CRDTOperation{
		ID:          fmt.Sprintf("merge-%d", time.Now().UnixNano()),
		Type:        CRDTTypeORSet,
		WorkspaceID: op1.WorkspaceID,
		ResourceID:  op1.ResourceID,
		ResourceType: op1.ResourceType,
		Operation:   "merge",
		Timestamp:   time.Now(),
		UserID:      crm.nodeID,
		VectorClock: crm.mergeVectorClocks(op1.VectorClock, op2.VectorClock),
	}, nil
}

// resolveLWWRegConflict resolves an LWW Register conflict
func (crm *ConflictResolutionManager) resolveLWWRegConflict(op1, op2 *CRDTOperation) (*CRDTOperation, error) {
	// Last-write-wins: choose operation with later timestamp
	// If timestamps equal, use node ID for deterministic ordering
	if op1.Timestamp.After(op2.Timestamp) {
		return op1, nil
	} else if op2.Timestamp.After(op1.Timestamp) {
		return op2, nil
	} else {
		// Timestamps equal, use node ID
		if op1.UserID > op2.UserID {
			return op1, nil
		}
		return op2, nil
	}
}

// mergeVectorClocks merges two vector clocks
func (crm *ConflictResolutionManager) mergeVectorClocks(vc1, vc2 map[string]int64) map[string]int64 {
	result := make(map[string]int64)
	
	// Add all from vc1
	for k, v := range vc1 {
		result[k] = v
	}
	
	// Merge vc2 (take max)
	for k, v := range vc2 {
		if current, exists := result[k]; !exists || v > current {
			result[k] = v
		}
	}
	
	return result
}

// Conflict represents a detected conflict
type Conflict struct {
	ResourceID   string            `json:"resourceId"`
	ResourceType string            `json:"resourceType"`
	Operations   []*CRDTOperation  `json:"operations"`
	Type         ConflictType      `json:"type"`
	Resolved     bool              `json:"resolved"`
	Resolution   *CRDTOperation    `json:"resolution,omitempty"`
}

// ConflictType represents the type of conflict
type ConflictType string

const (
	ConflictTypeConcurrent ConflictType = "concurrent" // Concurrent modifications
	ConflictTypeDivergent  ConflictType = "divergent" // Divergent state
	ConflictTypeLostUpdate ConflictType = "lost_update" // Lost update problem
)

