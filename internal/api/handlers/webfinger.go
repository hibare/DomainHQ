package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ggicci/httpin"
	"github.com/hibare/DomainHQ/internal/config"
)

const REL = "http://openid.net/specs/connect/1.0/issuer"

type WebFingerParams struct {
	Resource string `in:"query=resource;required"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type WebFingerResponse struct {
	Subject string `json:"subject"`
	Links   []Link `json:"links"`
}

func WebFinger(w http.ResponseWriter, r *http.Request) {
	requestInput := r.Context().Value(httpin.Input).(*WebFingerParams)

	resource := requestInput.Resource
	parts := strings.SplitN(resource, ":", 2)
	if len(parts) != 2 || parts[0] != "acct" {
		http.Error(w, "Invalid 'resource' parameter", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(parts[1], fmt.Sprintf("@%s", config.Current.WebFinger.Domain)) {
		log.Warnf("Resource '%s' does not match domain '%s'", resource, config.Current.WebFinger.Domain)
		http.Error(w, "Domain not allowed", http.StatusForbidden)
		return
	}

	// ToDo: Validate account with IDP

	resp := &WebFingerResponse{
		Subject: resource,
		Links: []Link{
			{
				Rel:  REL,
				Href: config.Current.WebFinger.Resource,
			},
		},
	}
	log.Infof("Resource '%s' is allowed", resource)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
