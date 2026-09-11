package handlers
import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)
const maxQueryBodyBytes = 1 <<20

type QueryRequest struct {
	SQL string `json:"sql"`
	Params map[string]interface{} `json:"params,omitempty"`
}
type QueryResponse struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	ExecutedAt time.Time `json:"executed_at"`
	LatencyMS int64 `json:"latency_ms"`
	ResultCount int `json:result_count,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
func Liveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func Metrics(w http.ResponseWriter, r *http.Request) {
	stats :=map[string]interface{}{
		"accepted_total": 14050,
		"rejected_total": 12,
		"active_channels": 4,
		"layer_health":    "optimal",
	}
	writeJSON(w, http.StatusOK, stats)
}
func Query(w http.ResponseWriter, r *http.Request){
	if r.Method !=http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "only POST method is allowed"}
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body,maxQueryBodyBytes)
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]{"error": "invalid JSON body structure"})
		return
	}
	sqlQuery :=strings.TrimSpace(req.SQL)
	if sqlQuery == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"SQL statement cannot be blank"})
		return
	}
	if !validateSQL(sqlQuery) {
		writeJSON(w,http.StatusForbidden, map[string]string{"error": "security violation: only SELECT queries are authorized"})
		return
	}
	startTime := time.Now()
	mockResult := mockDispatch(sqlQuery)
	latency := time.Since(startTime).Milliseconds()
	response :=QueryResponse{
		Status: "success",
		ExecutedAt: time.Now(),
		LatencyMS: latency,
		ResultCount: len(mockResult),
		Data: mockResult,
	}
	writeJSON(w, http.StatusOK, response)
}
func validateSQL(sql string) bool {
	upperSQL := strings.ToUpper(sql)
	if !strings.HasPrefix(upperSQL, "SELECT") {
		return false
	}
	forbidden := []string{"DROP", "DELETE", "TRUNCATE", "ALTER", "GRANT","REVOKE"}
	for _, keyword := range forbidden {
		if strings.Contains(upperSQL, keyword) {
			return false
		}
	}
	return true 
}
func mockDispatch(sql string) []map[string]interface{}{
	time.Sleep(2 * time.Millisecond)
	return []map[string]interface{}{
		{"id":1, "name": "Enterprise Cluster Alpha", "status": "active"},
		{"id":2, "name": "Global Gateway Mesh", "status": "synchronized"},
	}
}
func writeJsON(w,  http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}