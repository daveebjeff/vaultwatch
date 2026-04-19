package notify

import (
	"encoding/json"
	"net/http"
	"time"
)

// StatusPageHandler returns an http.HandlerFunc that serialises the current
// snapshot state from sn as a JSON response. Intended for lightweight
// internal status endpoints.
func StatusPageHandler(sn *SnapshotNotifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := sn.All()
		type entry struct {
			Path       string    `json:"path"`
			Status     string    `json:"status"`
			ReceivedAt time.Time `json:"received_at"`
		}
		result := make([]entry, 0, len(all))
		for _, snap := range all {
			result = append(result, entry{
				Path:       snap.Message.Path,
				Status:     string(snap.Message.Status),
				ReceivedAt: snap.ReceivedAt,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
