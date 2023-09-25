package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ggicci/httpin"
	"github.com/hibare/DomainHQ/internal/config"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	"github.com/rs/zerolog/log"
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
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'resource' parameter"))
		return
	}

	if !strings.HasSuffix(parts[1], fmt.Sprintf("@%s", config.Current.WebFinger.Domain)) {
		log.Warn().Msgf("Resource '%s' does not match domain '%s'", resource, config.Current.WebFinger.Domain)
		commonHttp.WriteErrorResponse(w, http.StatusForbidden, fmt.Errorf("domain not allowed"))
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
	log.Info().Msgf("Resource '%s' is allowed", resource)

	commonHttp.WriteJsonResponse(w, http.StatusOK, resp)
}
