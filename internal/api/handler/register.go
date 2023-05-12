package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-diploma/internal/api/models"
	"go-diploma/internal/service/auth"
	"io"
	"log"
	"net/http"
)

func (h *LoyaltyHandler) Register(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		msg := fmt.Sprintf("Unsupported content type [%s]", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	cred := models.CredentialsJSON{}

	dec := json.NewDecoder(r.Body)
	if errDec := dec.Decode(&cred); errDec != nil && errDec != io.EOF {
		msg := fmt.Sprintf("Cannot parse credentials: %s", errDec.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	log.Println(cred)

	errReg := h.AuthSrv.Register(cred.ToCredentials())
	if errors.Is(errReg, auth.ErrLoginAlreadyTaken) {
		http.Error(w, "Login already taken", http.StatusConflict)
		return
	}
	if errReg != nil {
		log.Printf("Cannot register user [%v]: %s", cred, errReg.Error())
		http.Error(w, "Cannot register user because of error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/api/user/login", http.StatusTemporaryRedirect)
}
