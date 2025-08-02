# About

[![Build](https://github.com/rgl/github-actions-validate-jwt/actions/workflows/build.yml/badge.svg)](https://github.com/rgl/github-actions-validate-jwt/actions/workflows/build.yml)

This validates a GitHub Actions ID Token JWT using the keys available at its JWKS endpoint.

A GitHub Actions ID Token JWT is a secret string that can be used to authenticate a particular GitHub Actions Workflow Job in 3rd party services (like HashiCorp Vault or Microsoft Entra ID).

A GitHub Actions ID Token JWT is [requested at runtime by the GitHub Actions Workflow Job to the GitHub OIDC Identity Provider](https://docs.github.com/en/actions/concepts/security/openid-connect).

A JWT is a structured string separated by dot characters; for example, a custom ID token JWT, something alike:

```
eyJhbGciOiJSUzI1NiIsImtpZCI6IjM4ODI2YjE3LTZhMzAtNWY5Yi1iMTY5LThiZWI4MjAyZjcyMyIsInR5cCI6IkpXVCIsIng1dCI6InlrTmFZNHFNX3RhNGsyVGdaT0NFWUxrY1lsQSJ9.eyJhY3RvciI6InJnbCIsImFjdG9yX2lkIjoiNDMzNTYiLCJhdWQiOiJodHRwczovL2V4YW1wbGUuY29tIiwiYmFzZV9yZWYiOiIiLCJldmVudF9uYW1lIjoicHVzaCIsImV4cCI6MTc0MzI2NzgyNywiaGVhZF9yZWYiOiIiLCJpYXQiOjE3NDMyNDYyMjcsImlzcyI6Imh0dHBzOi8vdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJqb2Jfd29ya2Zsb3dfcmVmIjoicmdsL2dpdGh1Yi1hY3Rpb25zLXZhbGlkYXRlLWp3dC8uZ2l0aHViL3dvcmtmbG93cy9idWlsZC55bWxAcmVmcy9oZWFkcy9tYWluIiwiam9iX3dvcmtmbG93X3NoYSI6IjBmNjYxYTRhMDExZDc3Y2IyMWZhNGRjNzAxM2Q5YjZmOTI4NjY2ZWQiLCJqdGkiOiI1ODc3MmM5Zi00ZDc1LTRmYzgtODk1Ni03ZGQ1MjI1NjlhMDAiLCJuYmYiOjE3NDMyNDU5MjcsInJlZiI6InJlZnMvaGVhZHMvbWFpbiIsInJlZl9wcm90ZWN0ZWQiOiJmYWxzZSIsInJlZl90eXBlIjoiYnJhbmNoIiwicmVwb3NpdG9yeSI6InJnbC9naXRodWItYWN0aW9ucy12YWxpZGF0ZS1qd3QiLCJyZXBvc2l0b3J5X2lkIjoiOTU3MDE0NDY2IiwicmVwb3NpdG9yeV9vd25lciI6InJnbCIsInJlcG9zaXRvcnlfb3duZXJfaWQiOiI0MzM1NiIsInJlcG9zaXRvcnlfdmlzaWJpbGl0eSI6InB1YmxpYyIsInJ1bl9hdHRlbXB0IjoiMSIsInJ1bl9pZCI6IjE0MTQ0OTk2NTEyIiwicnVuX251bWJlciI6IjMiLCJydW5uZXJfZW52aXJvbm1lbnQiOiJnaXRodWItaG9zdGVkIiwic2hhIjoiMGY2NjFhNGEwMTFkNzdjYjIxZmE0ZGM3MDEzZDliNmY5Mjg2NjZlZCIsInN1YiI6InJlcG86cmdsL2dpdGh1Yi1hY3Rpb25zLXZhbGlkYXRlLWp3dDpyZWY6cmVmcy9oZWFkcy9tYWluIiwid29ya2Zsb3ciOiJCdWlsZCIsIndvcmtmbG93X3JlZiI6InJnbC9naXRodWItYWN0aW9ucy12YWxpZGF0ZS1qd3QvLmdpdGh1Yi93b3JrZmxvd3MvYnVpbGQueW1sQHJlZnMvaGVhZHMvbWFpbiIsIndvcmtmbG93X3NoYSI6IjBmNjYxYTRhMDExZDc3Y2IyMWZhNGRjNzAxM2Q5YjZmOTI4NjY2ZWQifQ.wITH_EaL8MXkr7rIPjHbAcuT31fKzjocC6hJZn_zhpce2kqhIrPKnYM0vvM4M34cazahfIwZGXMxcojvPPr9X5tleQYbZtabrCmQMw9HMkRxTpw3NBWJDud8oxofbPiKhWwlCz7b2PIEzp31SNsAvD1D1hyN0MFE2wOPLsm3swouxBNYmpyS565cJQx7V5v2VFgrhTIluYPVpFVj-3K4NY4WQrKlBcEQjky26S-P0m6Ksfr7J0grVxkvx2WKqn5YAjJe1FiMpZcZeID1LYXWpc2K7VjPkYlBnRbxCbX9S0ZTd_g4GXFZRQNYcDyHQ2t1jDi5gNK6o6QPCaHhJaDuNg
```

When split by dot and decoded it has a header, payload and signature.

In this case, the header is:

```json
{
  "alg": "RS256",
  "kid": "38826b17-6a30-5f9b-b169-8beb8202f723",
  "typ": "JWT",
  "x5t": "ykNaY4qM_ta4k2TgZOCEYLkcYlA"
}
```

The payload is:

```json
{
  "actor": "rgl",
  "actor_id": "43356",
  "aud": "https://example.com",
  "base_ref": "",
  "event_name": "push",
  "exp": 1743267827, // NB valid for 6h.
  "head_ref": "",
  "iat": 1743246227,
  "iss": "https://token.actions.githubusercontent.com",
  "job_workflow_ref": "rgl/github-actions-validate-jwt/.github/workflows/build.yml@refs/heads/main",
  "job_workflow_sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed",
  "jti": "58772c9f-4d75-4fc8-8956-7dd522569a00",
  "nbf": 1743245927,
  "ref": "refs/heads/main",
  "ref_protected": "false",
  "ref_type": "branch",
  "repository": "rgl/github-actions-validate-jwt",
  "repository_id": "957014466",
  "repository_owner": "rgl",
  "repository_owner_id": "43356",
  "repository_visibility": "public",
  "run_attempt": "1",
  "run_id": "14144996512",
  "run_number": "3",
  "runner_environment": "github-hosted",
  "sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed",
  "sub": "repo:rgl/github-actions-validate-jwt:ref:refs/heads/main",
  "workflow": "Build",
  "workflow_ref": "rgl/github-actions-validate-jwt/.github/workflows/build.yml@refs/heads/main",
  "workflow_sha": "0f661a4a011d77cb21fa4dc7013d9b6f928666ed"
}
```

And the signature is the value from the 3rd part of the JWT string.

Before a JWT can be used it must be validated. In this particular example the JWT can be validated with:

```go
RSASHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  gitHubJwtKeySet.getPublicKey(header.kid))
```

The above public key should be retrieved from the [GitHub Actions JWKS endpoint](https://token.actions.githubusercontent.com/.well-known/jwks).

To see how all of this can be done read the [main.go](main.go) and the [build.yml github actions workflow](.github/workflows/build.yml) file.
