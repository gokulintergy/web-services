package middleware

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"net/http"

	_json "github.com/mappcpd/web-services/cmd/webd/router/handlers/responder"
	j_ "github.com/mappcpd/web-services/internal/platform/jwt"
)

// UserAuthToken is a global AuthToken that is set up by the middleware for convenience
var UserAuthToken j_.AuthToken

// ValidateToken validate the JSON web token passed in the Authorization header. For now
// a POST request to /auth simply returns, without checking the token, as this is
// a request to authenticate and get a new token.
func ValidateToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// pass through when request is preflight http OPTIONS
	if r.Method == http.MethodOptions {
		fmt.Println("Bypassing ValidateToken() middleware for OPTIONS request")
		next(w, r)
		return
	}

	p := _json.Payload{}

	// Get the token from the auth header, 'Bearer' seems useless but this is an OAuth2 standard
	// Authorization: Bearer [jwt]
	a := r.Header.Get("Authorization")
	t, err := j_.FromHeader(a)
	if err != nil {
		p.Message = _json.Message{http.StatusBadRequest, "failure", err.Error()}
		p.Send(w)
		return
	}

	// Set the global AuthToken
	UserAuthToken, err = j_.Check(t)
	if err != nil {
		p.Message = _json.Message{http.StatusUnauthorized, "failure", "Authorization failed: " + err.Error()}
		p.Send(w)
		return
	}

	next(w, r)
}

// AdminScope checks that the auth token belongs to an admin
func AdminScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	p := _json.Payload{}

	if UserAuthToken.CheckScope("admin") == false {
		p.Message = _json.Message{http.StatusUnauthorized, "failed", "Admin Scope Required: token does not belong to an admin user"}
		p.Send(w)
		return
	}

	next(w, r)
}

// MemberScope checks that the auth token belongs to an admin
func MemberScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// pass through when request is preflight http OPTIONS
	if r.Method == http.MethodOptions {
		fmt.Println("Bypassing MemberScope() middleware for OPTIONS request")
		next(w, r)
		return
	}

	p := _json.Payload{}

	if UserAuthToken.CheckScope("member") == false {
		p.Message = _json.Message{http.StatusUnauthorized, "failed", "Member Scope Required: token does not belong to a member user"}
		p.Send(w)
		return
	}

	// mux.Vars is not available here as the path as not been parsed yet... maybe because of the way
	// the middleware has been set up. So, for now just split the path and look for '/members/{id}
	vars := strings.Split(r.URL.Path, "/")

	// We are only interested in validating a member id as part of the path when there
	// is another resource after the id, that is: /members/{memberId}/resource/{resourceId}
	// but it could be /v1/members... so range over vars until we find /members... and make sure
	// there's still at least two items thereafter
	c := len(vars)
	if c > 3 {
		for i := range vars {
			c--
			if string(vars[i]) == "members" && (c-i) >= 2 {
				log.Println("Checking member id on path matches token")
				mid, err := strconv.Atoi(vars[i+1])
				if err != nil {
					p.Message = _json.Message{http.StatusBadRequest, "failed", "Member id in path appears to be invalid"}
					p.Send(w)
					return
				}
				if UserAuthToken.Claims.ID != mid {
					p.Message = _json.Message{http.StatusUnauthorized, "failed", "Member id in path does not match token"}
					p.Send(w)
					return
				}
				break
			}
		}
	}

	next(w, r)
}
