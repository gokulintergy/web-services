package rest

import (
	"net/http"

	"github.com/cardiacsociety/web-services/internal/qualification"
)

// Qualifications fetches list of Qualifications
func Qualifications(w http.ResponseWriter, _ *http.Request) {

	p := NewResponder(UserAuthToken.Encoded)

	xq, err := qualification.All(DS)
	if err != nil {
		p.Message = Message{http.StatusInternalServerError, "failed", err.Error()}
		p.Send(w)
		return
	}

	// All good
	p.Message = Message{http.StatusOK, "success", "Data retrieved from " + DS.MySQL.Desc}
	p.Data = xq
	m := make(map[string]interface{})
	m["count"] = len(xq)
	m["description"] = "List of Qualifications"
	p.Meta = m
	p.Send(w)
}
