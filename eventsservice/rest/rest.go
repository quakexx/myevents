package rest

import(
	"lib/persistence"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func ServeAPI(endpoint string, tlsendpoint string, databasehandler persistence.DatabaseHandler)(chan error, chan error) {
	handler := New(databasehandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler) 
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler) 
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler) 

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)
	go func() {
		server := handlers.CORS()(r)
		httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint,"cert.pem","key.pem", server)
	}()

	go func() {
		server := handlers.CORS()(r)
		httpErrChan <- http.ListenAndServe(endpoint,server)
	}()

	return httpErrChan, httptlsErrChan
}
