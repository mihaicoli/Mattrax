package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

func Devices(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := srv.DB.GetDevices(r.Context())
		if err != nil {
			log.Printf("[GetDevices Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(devices); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Device(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		device, err := srv.DB.GetBasicDevice(r.Context(), int32(id))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDevice Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(device); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeviceInformation(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		device, err := srv.DB.GetDevice(r.Context(), int32(id))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetDevice Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]map[string]interface{}{
			"Device Information": {
				"Computer Name": device.Name,
				// "Serial Number":
			},
			"Software Information": {
				"Operating System":         "Windows 10", // TODO
				"Operating System Version": device.OperatingSystem,
			},
			"MDM": {
				"Last Seen":        device.Lastseen,
				"Last Seen Status": device.LastseenStatus,
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeviceScope(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		groups, err := srv.DB.GetBasicDeviceScopedGroups(r.Context(), int32(id))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDeviceScopedGroups Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		policies, err := srv.DB.GetBasicDeviceScopedPolicies(r.Context(), int32(id))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDeviceScopedPolicies Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"groups":   groups,
			"policies": policies,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
