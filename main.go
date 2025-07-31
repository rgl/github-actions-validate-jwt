package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("the %s environment variable must be set", name)
	}
	return v
}

func main() {
	issuer := "https://token.actions.githubusercontent.com"
	audience := "https://example.com"

	// see https://token.actions.githubusercontent.com/.well-known/openid-configuration
	jwksURL := fmt.Sprintf("%s/.well-known/jwks", issuer)

	if len(os.Args) != 2 {
		log.Fatalf("you must pass the github actions workflow job id token environment variable name as the single command line argument")
	}
	jobJWT := getEnv(os.Args[1])

	// fetch the github actions jwt key set.
	//
	// a key set is public object alike:
	//
	// 		{
	// 			"keys": [
	// 				{
	// 					"kty": "RSA",
	// 					"alg": "RS256",
	// 					"use": "sig",
	// 					"kid": "38826b17-6a30-5f9b-b169-8beb8202f723",
	// 					"n": "5Manmy-zwsk3wEftXNdKFZec4rSWENW4jTGevlvAcU9z3bgLBogQVvqYLtu9baVm2B3rfe5onadobq8po5UakJ0YsTiiEfXWdST7YI2Sdkvv-hOYMcZKYZ4dFvuSO1vQ2DgEkw_OZNiYI1S518MWEcNxnPU5u67zkawAGsLlmXNbOylgVfBRJrG8gj6scr-sBs4LaCa3kg5IuaCHe1pB-nSYHovGV_z0egE83C098FfwO1dNZBWeo4Obhb5Z-ZYFLJcZfngMY0zJnCVNmpHQWOgxfGikh3cwi4MYrFrbB4NTlxbrQ3bL-rGKR5X318veyDlo8Dyz2KWMobT4wB9U1Q",
	// 					"e": "AQAB",
	// 					"x5c": [
	// 						"MIIDKzCCAhOgAwIBAgIUDnwm6eRIqGFA3o/P1oBrChvx/nowDQYJKoZIhvcNAQELBQAwJTEjMCEGA1UEAwwaYWN0aW9ucy5zZWxmLXNpZ25lZC5naXRodWIwHhcNMjQwMTIzMTUyNTM2WhcNMzQwMTIwMTUyNTM2WjAlMSMwIQYDVQQDDBphY3Rpb25zLnNlbGYtc2lnbmVkLmdpdGh1YjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOTGp5svs8LJN8BH7VzXShWXnOK0lhDVuI0xnr5bwHFPc924CwaIEFb6mC7bvW2lZtgd633uaJ2naG6vKaOVGpCdGLE4ohH11nUk+2CNknZL7/oTmDHGSmGeHRb7kjtb0Ng4BJMPzmTYmCNUudfDFhHDcZz1Obuu85GsABrC5ZlzWzspYFXwUSaxvII+rHK/rAbOC2gmt5IOSLmgh3taQfp0mB6Lxlf89HoBPNwtPfBX8DtXTWQVnqODm4W+WfmWBSyXGX54DGNMyZwlTZqR0FjoMXxopId3MIuDGKxa2weDU5cW60N2y/qxikeV99fL3sg5aPA8s9iljKG0+MAfVNUCAwEAAaNTMFEwHQYDVR0OBBYEFIPALo5VanJ6E1B9eLQgGO+uGV65MB8GA1UdIwQYMBaAFIPALo5VanJ6E1B9eLQgGO+uGV65MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAGS0hZE+DqKIRi49Z2KDOMOaSZnAYgqq6ws9HJHT09MXWlMHB8E/apvy2ZuFrcSu14ZLweJid+PrrooXEXEO6azEakzCjeUb9G1QwlzP4CkTcMGCw1Snh3jWZIuKaw21f7mp2rQ+YNltgHVDKY2s8AD273E8musEsWxJl80/MNvMie8Hfh4n4/Xl2r6t1YPmUJMoXAXdTBb0hkPy1fUu3r2T+1oi7Rw6kuVDfAZjaHupNHzJeDOg2KxUoK/GF2/M2qpVrd19Pv/JXNkQXRE4DFbErMmA7tXpp1tkXJRPhFui/Pv5H9cPgObEf9x6W4KnCXzT3ReeeRDKF8SqGTPELsc="
	// 					],
	// 					"x5t": "ykNaY4qM_ta4k2TgZOCEYLkcYlA"
	// 				},
	// 			]
	// 		}
	log.Printf("Getting the GitHub Actions JWT public key set from the JWKS endpoint at %s...", jwksURL)
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		log.Fatalf("failed to parse JWK from %s: %v", jwksURL, err)
	}
	if keySet.Len() < 1 {
		log.Fatalf("%s did not return any key", jwksURL)
	}

	// parse and validate the job id token jwt against the github actions jwt key set.
	//
	// a job jwt is a private string alike:
	//
	// 		eyJhbGciOiJSUzI1NiIsImtpZCI6IjM4ODI2YjE3LTZhMzAtNWY5Yi1iMTY5LThiZWI4MjAyZjcyMyIsInR5cCI6IkpXVCIsIng1dCI6InlrTmFZNHFNX3RhNGsyVGdaT0NFWUxrY1lsQSJ9.eyJhY3RvciI6InJnbCIsImFjdG9yX2lkIjoiNDMzNTYiLCJhdWQiOiJodHRwczovL2V4YW1wbGUuY29tIiwiYmFzZV9yZWYiOiIiLCJldmVudF9uYW1lIjoicHVzaCIsImV4cCI6MTc0MzI2NzgyNywiaGVhZF9yZWYiOiIiLCJpYXQiOjE3NDMyNDYyMjcsImlzcyI6Imh0dHBzOi8vdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJqb2Jfd29ya2Zsb3dfcmVmIjoicmdsL2dpdGh1Yi1hY3Rpb25zLXZhbGlkYXRlLWp3dC8uZ2l0aHViL3dvcmtmbG93cy9idWlsZC55bWxAcmVmcy9oZWFkcy9tYWluIiwiam9iX3dvcmtmbG93X3NoYSI6IjBmNjYxYTRhMDExZDc3Y2IyMWZhNGRjNzAxM2Q5YjZmOTI4NjY2ZWQiLCJqdGkiOiI1ODc3MmM5Zi00ZDc1LTRmYzgtODk1Ni03ZGQ1MjI1NjlhMDAiLCJuYmYiOjE3NDMyNDU5MjcsInJlZiI6InJlZnMvaGVhZHMvbWFpbiIsInJlZl9wcm90ZWN0ZWQiOiJmYWxzZSIsInJlZl90eXBlIjoiYnJhbmNoIiwicmVwb3NpdG9yeSI6InJnbC9naXRodWItYWN0aW9ucy12YWxpZGF0ZS1qd3QiLCJyZXBvc2l0b3J5X2lkIjoiOTU3MDE0NDY2IiwicmVwb3NpdG9yeV9vd25lciI6InJnbCIsInJlcG9zaXRvcnlfb3duZXJfaWQiOiI0MzM1NiIsInJlcG9zaXRvcnlfdmlzaWJpbGl0eSI6InB1YmxpYyIsInJ1bl9hdHRlbXB0IjoiMSIsInJ1bl9pZCI6IjE0MTQ0OTk2NTEyIiwicnVuX251bWJlciI6IjMiLCJydW5uZXJfZW52aXJvbm1lbnQiOiJnaXRodWItaG9zdGVkIiwic2hhIjoiMGY2NjFhNGEwMTFkNzdjYjIxZmE0ZGM3MDEzZDliNmY5Mjg2NjZlZCIsInN1YiI6InJlcG86cmdsL2dpdGh1Yi1hY3Rpb25zLXZhbGlkYXRlLWp3dDpyZWY6cmVmcy9oZWFkcy9tYWluIiwid29ya2Zsb3ciOiJCdWlsZCIsIndvcmtmbG93X3JlZiI6InJnbC9naXRodWItYWN0aW9ucy12YWxpZGF0ZS1qd3QvLmdpdGh1Yi93b3JrZmxvd3MvYnVpbGQueW1sQHJlZnMvaGVhZHMvbWFpbiIsIndvcmtmbG93X3NoYSI6IjBmNjYxYTRhMDExZDc3Y2IyMWZhNGRjNzAxM2Q5YjZmOTI4NjY2ZWQifQ.wITH_EaL8MXkr7rIPjHbAcuT31fKzjocC6hJZn_zhpce2kqhIrPKnYM0vvM4M34cazahfIwZGXMxcojvPPr9X5tleQYbZtabrCmQMw9HMkRxTpw3NBWJDud8oxofbPiKhWwlCz7b2PIEzp31SNsAvD1D1hyN0MFE2wOPLsm3swouxBNYmpyS565cJQx7V5v2VFgrhTIluYPVpFVj-3K4NY4WQrKlBcEQjky26S-P0m6Ksfr7J0grVxkvx2WKqn5YAjJe1FiMpZcZeID1LYXWpc2K7VjPkYlBnRbxCbX9S0ZTd_g4GXFZRQNYcDyHQ2t1jDi5gNK6o6QPCaHhJaDuNg
	//
	// and decoded as a private object is alike:
	//
	// 		header:
	//
	// 			{
	// 				"alg": "RS256",
	// 				"kid": "38826b17-6a30-5f9b-b169-8beb8202f723",
	// 				"typ": "JWT",
	// 				"x5t": "ykNaY4qM_ta4k2TgZOCEYLkcYlA"
	// 			}
	//
	// 		payload:
	//
	//			{
	// 				"actor": "rgl",
	// 				"actor_id": "43356",
	// 				"aud": "https://example.com",
	// 				"base_ref": "",
	// 				"event_name": "push",
	// 				"exp": 1743267827, // NB valid for 6h.
	// 				"head_ref": "",
	// 				"iat": 1743246227,
	// 				"iss": "https://token.actions.githubusercontent.com",
	// 				"job_workflow_ref": "rgl/github-actions-validate-jwt/.github/workflows/build.yml@refs/heads/main",
	// 				"job_workflow_sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed",
	// 				"jti": "58772c9f-4d75-4fc8-8956-7dd522569a00",
	// 				"nbf": 1743245927,
	// 				"ref": "refs/heads/main",
	// 				"ref_protected": "false",
	// 				"ref_type": "branch",
	// 				"repository": "rgl/github-actions-validate-jwt",
	// 				"repository_id": "957014466",
	// 				"repository_owner": "rgl",
	// 				"repository_owner_id": "43356",
	// 				"repository_visibility": "public",
	// 				"run_attempt": "1",
	// 				"run_id": "14144996512",
	// 				"run_number": "3",
	// 				"runner_environment": "github-hosted",
	// 				"sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed",
	// 				"sub": "repo:rgl/github-actions-validate-jwt:ref:refs/heads/main",
	// 				"workflow": "Build",
	// 				"workflow_ref": "rgl/github-actions-validate-jwt/.github/workflows/build.yml@refs/heads/main",
	// 				"workflow_sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed"
	// 			}
	//
	//		signature:
	//
	//			the value is the 3rd part of the jwt.
	//
	//			in this particular example the jwt can be validated with:
	//
	//				RSASHA256(
	//   				base64UrlEncode(header) + "." + base64UrlEncode(payload),
	//					gitHubActionsJwtKeySet.getKey(header.kid))
	log.Println("Validating GitHub Actions job JWT...")
	token, err := jwt.ParseString(jobJWT, jwt.WithAudience(audience), jwt.WithIssuer(issuer), jwt.WithKeySet(keySet))
	if err != nil {
		log.Fatalf("failed to validate the jwt: %v", err)
	}

	var repository string
	if err = token.Get("repository", &repository); err != nil {
		log.Fatalf("failed to get the repository claim from the jwt: %v", err)
	}
	log.Printf("jwt is valid for repository %s", repository)

	// dump the jwt claims (sorted by claim name).
	claims := make([]string, 0, len(token.Keys()))
	for _, k := range token.Keys() {
		var v any
		if err = token.Get(k, &v); err != nil {
			log.Fatalf("failed to get the %s claim from the jwt: %v", k, err)
		}
		switch v := v.(type) {
		case string:
			claims = append(claims, fmt.Sprintf("%s=%s", k, v))
		case []string:
			for _, v := range v {
				claims = append(claims, fmt.Sprintf("%s=%s", k, v))
			}
		case time.Time:
			claims = append(claims, fmt.Sprintf("%s=%s", k, v.Format("2006-01-02T15:04:05-0700")))
		default:
			log.Printf("WARNING: skipping the %s claim with type %T value %v", k, v, v)
			continue
		}
	}
	sort.Strings(claims)
	for _, claim := range claims {
		log.Printf("jwt claim: %s", claim)
	}
}
