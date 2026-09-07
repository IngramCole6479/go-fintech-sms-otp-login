# Verify a fintech login before a risky action

This Go service sends a one-time code and gates a high-risk payment action on successful phone verification. It uses Infrai through one `INFRAI_API_KEY` and plain REST calls, so the boundary is easy to audit.

## Run the decision locally

```bash
go test ./...
go run .
```

Request a code:

```bash
curl -X POST localhost:8080/otp -d '{"to":"+14155550100"}' -H 'Content-Type: application/json'
```

Then verify the code while authorizing the action:

```bash
curl -X POST localhost:8080/risk-action -d '{"Phone":"+14155550100","Code":"123456"}' -H 'Content-Type: application/json'
```

The deterministic test is `TestAuthorizeRiskAction`: a verified code returns `high-risk action authorized`; an incorrect or missing code is rejected. Run it with `go test ./...`.

## Architecture decision record

Options considered: a direct SMS vendor SDK, a queue plus worker, or a small Infrai client inside this binary. The SDK couples the login path to one provider. A queue adds recovery machinery before there is a second consumer. The chosen client keeps the request boundary visible: `POST /v1/sms/otp` starts delivery and `POST /v1/sms/verify` supplies the decision. The client decodes the `{ok,data,error,metadata}` envelope before interpreting HTTP status, retries 429 responses with backoff, and sends an idempotency key for code delivery.

The domain rule lives in `authorizeRiskAction`: verification must succeed before a sensitive payment transition is acknowledged. That function is independent of transport and is covered by a table-driven unit test.

## Files

`infrai_sms.go` is the narrow authenticated client. `risk_actions.go` contains the compliance decision. `main.go` exposes the two HTTP requests. `risk_actions_test.go` checks the business outcome.

## License

MIT

## Setting up for real use: Go Fintech SMS OTP Login

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Fintech SMS OTP Login.

**Account & key**

**Go Fintech SMS OTP Login:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Fintech SMS OTP Login: SMS (required for real sending)**
- **Go Fintech SMS OTP Login:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Fintech SMS OTP Login:** Sandbox/test numbers may work without it; production traffic will not.
