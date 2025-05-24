package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"RushOrder/session"
)

func CustomerLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nama := r.FormValue("nama")
	mejaStr := r.FormValue("meja")
	meja, err := strconv.Atoi(mejaStr)
	if err != nil {
		http.Error(w, "Nomor meja tidak valid", http.StatusBadRequest)
		return
	}

	err = session.CreateSession(w, r, nama, meja)
	if err != nil {
		http.Error(w, "Gagal membuat sesi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}

func GetCustomerSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	customerSession, err := session.GetSession(r)
	if err != nil {
		http.Error(w, "Gagal mendapatkan sesi", http.StatusInternalServerError)
		return
	}
	if customerSession == nil {
		http.Error(w, "Sesi tidak ditemukan", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customerSession)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session.ClearCustomerSession(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
