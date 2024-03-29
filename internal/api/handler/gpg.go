package handler

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/ggicci/httpin"
	"github.com/hibare/DomainHQ/internal/models"
	"github.com/hibare/GoCommon/v2/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	"gorm.io/gorm"
)

type GPGLookupParams struct {
	Op     string `in:"query=op"`
	Search string `in:"query=search;required"`
}

type GPGKeyAddParams struct {
	KeyText string `in:"form=keytext"`
}

const OPGet = "get"

func GPGPubKeyLookup(tx *gorm.DB, w http.ResponseWriter, r *http.Request) {
	requestInput := r.Context().Value(httpin.Input).(*GPGLookupParams)

	if requestInput.Op == OPGet {
		key, err := models.LookupPubKey(tx, requestInput.Search)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				commonHttp.WriteErrorResponse(w, http.StatusNotFound, fmt.Errorf("key not found"))
				return
			}

			slog.Error("Error looking up key", "error", err)
			commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, key.PublicKey)
		return
	} else {
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid op"))
		return
	}
}

func GPGPubKeyAdd(tx *gorm.DB, w http.ResponseWriter, r *http.Request) {
	requestInput := r.Context().Value(httpin.Input).(*GPGKeyAddParams)

	parsedKey, err := models.ParsePubKey(requestInput.KeyText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.AddPubKey(tx, &parsedKey)
	if err != nil {
		slog.Error("Error adding key", "error", err)
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServerError)
		return
	}

	commonHttp.WriteJsonResponse(w, http.StatusOK, "key added")
}
