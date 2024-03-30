package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// RawJWTCustomClaims is a type used to deserialize the raw claims from a JWT
type RawJWTCustomClaims struct {
	// RegisteredClaims contains default JWT claims as expressed by the RFC
	jwt.RegisteredClaims

	ClientID          string   `json:"clientId"`
	PreferredUsername string   `json:"preferred_username"`
	Name              string   `json:"name"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	Organization      string   `json:"organization"`
	GroupFull         []string `json:"group-full"`
}

// ToCustomClaims converts a RawJWTCustomClaims to a CustomClaims
func (rawClaims RawJWTCustomClaims) ToCustomClaims() CustomClaims {
	isServiceAccount := len(rawClaims.ClientID) > 0

	if isServiceAccount {
		return CustomClaims{
			RegisteredClaims:  rawClaims.RegisteredClaims,
			IsServiceAccount:  isServiceAccount,
			PreferredUsername: rawClaims.ClientID,
			Name:              rawClaims.ClientID,
			GivenName:         "service",
			FamilyName:        "account",
			Email:             "",
			Organization:      rawClaims.Organization,
		}
	} else {
		return CustomClaims{
			RegisteredClaims:  rawClaims.RegisteredClaims,
			IsServiceAccount:  isServiceAccount,
			PreferredUsername: rawClaims.PreferredUsername,
			Name:              rawClaims.Name,
			GivenName:         rawClaims.GivenName,
			FamilyName:        rawClaims.FamilyName,
			Email:             rawClaims.Email,
			GroupFull:         rawClaims.GroupFull,
			Organization:      rawClaims.Organization,
		}
	}
}

// MockCustomClaims constructs a mock CustomClaims object for testing endpoints that require authentication
func MockCustomClaims() CustomClaims {
	groupFull := []string{""}
	claims := []string{""}

	return CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   "subject",
			Audience:  claims,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * 10)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "1",
		},
		IsServiceAccount:  false,
		PreferredUsername: "preferredUsername",
		Name:              "name",
		GivenName:         "givenName",
		FamilyName:        "familyName",
		Email:             "email",
		GroupFull:         groupFull,
		Organization:      "organization",
	}
}

// CustomClaims is the JWT claim type actually used in the application. It's slightly processed
// from the raw claims on the JWT to make using it easier
type CustomClaims struct {
	jwt.RegisteredClaims

	IsServiceAccount  bool
	PreferredUsername string
	Name              string
	GivenName         string
	FamilyName        string
	Email             string
	GroupFull         []string
	Organization      string
}
