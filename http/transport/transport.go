package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"

	userendpoints "github.com/jumaroar-globant/go-bootcamp/http/endpoints/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

var (
	ErrMissingUserID = errors.New("missing user id")
)

// NewHTTPServer generates a new HTTPServer with its endpoints
func NewHTTPServer(usrEndpoints *userendpoints.UserEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/user/auth").Handler(
		httptransport.NewServer(
			usrEndpoints.Authenticate,
			decodeAuthRequest,
			encodeAuthResponse,
		),
	)

	r.Methods("POST").Path("/user").Handler(
		httptransport.NewServer(
			usrEndpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
	)

	r.Methods("GET").Path("/user/{id}").Handler(
		httptransport.NewServer(
			usrEndpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),
	)

	r.Methods("POST").Path("/user/{id}").Handler(
		httptransport.NewServer(
			usrEndpoints.UpdateUser,
			decodeUpdateUserRequest,
			encodeUpdateUserResponse,
		),
	)

	r.Methods("DELETE").Path("/user/{id}").Handler(
		httptransport.NewServer(
			usrEndpoints.DeleteUser,
			decodeDeleteUserRequest,
			encodeDeleteUserResponse,
		),
	)

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

func decodeCreateUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req shared.User
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*shared.User)
	return json.NewEncoder(w).Encode(res)
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userendpoints.GetUserRequest

	userID := mux.Vars(r)["id"]
	if userID == "" {
		return nil, ErrMissingUserID
	}

	req.UserID = userID

	return req, nil
}

func encodeGetUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*shared.User)
	return json.NewEncoder(w).Encode(res)
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req shared.User
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}

	req.ID = mux.Vars(r)["id"]

	return req, nil
}

func encodeUpdateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(*shared.User)
	return json.NewEncoder(w).Encode(res)
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userendpoints.DeleteUserRequest

	userID := mux.Vars(r)["id"]
	if userID == "" {
		return nil, ErrMissingUserID
	}

	req.UserID = userID

	return req, nil
}

func encodeDeleteUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.(userendpoints.DeleteUserResponse)
	return json.NewEncoder(w).Encode(res)
}
