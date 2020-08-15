package windows

import (
	"net/http"

	mattrax "github.com/mattrax/Mattrax/internal"
)

// WinProtocolID is the ID used in database protocol column to represent this protocol
const WinProtocolID = 1

// Mount initialise the MDM server
func Mount(srv *mattrax.Server) {
	// TODO: Replace with UI based Login Route
	srv.Router.HandleFunc("/Login.svc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(`<h3>MDM Federated Login</h3><form method="post" action="` + r.URL.Query().Get("appru") + `"><p><input type="hidden" name="wresult" value="TODOSpecialTokenWhichVerifiesAuth" /></p><input type="submit" value="Login" /></form>`))
	}).Name("login").Methods("GET")

	// TODO: Replace with UI based Login Route
	srv.Router.HandleFunc("/EnrollmentServer/TermsOfService.svc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(`<h3>AzureAD Term Of Service</h3><button onClick="acceptBtn()">Accept</button><script>function acceptBtn() { var urlParams = new URLSearchParams(window.location.search); if (!urlParams.has('redirect_uri')) { alert('Redirect url not found. Did you open this in your broswer?'); } else { window.location = urlParams.get('redirect_uri') + "?IsAccepted=true&OpaqueBlob=TODOCustomDataFromAzureAD"; } }</script>`))
	}).Name("azuread-tos").Methods("GET")

	srv.Router.HandleFunc("/ManagementServer/Manage.svc", Manage(srv)).Name("winmdm-manage").Methods("POST")
	srv.Router.HandleFunc("/EnrollmentServer/Policy.svc", Policy(srv)).Name("winmdm-policy").Methods("POST")
	srv.Router.HandleFunc("/EnrollmentServer/Enrollment.svc", Enrollment(srv)).Name("winmdm-enrollment").Methods("POST")

	srv.GlobalRouter.HandleFunc("/EnrollmentServer/Discovery.svc", Discovery(srv)).Methods("GET", "POST")
}
