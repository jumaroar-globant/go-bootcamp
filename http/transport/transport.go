package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"

	userendpoints "github.com/jumaroar-globant/go-bootcamp/http/endpoints/user"
)

func NewHTTPServer(usrEndpoints *userendpoints.UserEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/user/auth").Handler(
		httptransport.NewServer(
			usrEndpoints.Authenticate,
			decodeAuthRequest,
			encodeAuthResponse,
		))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func decodeAuthRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userendpoints.AuthenticationRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeAuthResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(userendpoints.AuthenticationResponse)
	return json.NewEncoder(w).Encode(res)
}
