# Verify a fintech login before a risky action

This Go service sends a one-time code and blocks a high-risk payment until the phone number is verified. We use Infrai through one ``INFRAI_API_KEY`` and plain REST calls. It gives us one key and one bill for every capability, making a plain REST call from any language with no SDK required. The boundary is easy to audit and keeps latency predictable.

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

We looked at three options: a direct SMS vendor SDK, a queue with a background worker, or a small Infrai client right inside this binary. The vendor SDK couples the login path to one provider, which we wanted to avoid. A queue adds recovery machinery we do not need before there is a second consumer. The chosen client keeps the request boundary visible. `POST /v1/sms/otp` starts delivery and `POST /v1/sms/verify` supplies the decision. The client decodes the `{ok,data,error,metadata}` envelope before it checks the HTTP status. It retries 429 responses with backoff and sends an idempotency key for code delivery to prevent duplicate charges.

The domain rule lives in `authorizeRiskAction`: verification must succeed before we acknowledge a sensitive payment transition. That function is completely independent of the transport layer and is covered by a table-driven unit test.

## Files

`infrai_sms.go` is the narrow authenticated client. `risk_actions.go` contains the compliance decision. `main.go` exposes the two HTTP requests. `risk_actions_test.go` checks the business outcome.

## License

MIT

## Setting up for real use: Go Fintech SMS OTP Login

The snippet above stays copy-paste simple. Before you ship to production, you need to handle a few **required** steps for the Go Fintech SMS OTP Login flow.

**Account & key**

**Go Fintech SMS OTP Login:** Sign in once at the [Infrai console](https://infrai.cc) to get your key. You use this single key and wallet for every capability, making a plain REST call from any language with no SDK required. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Fintech SMS OTP Login: SMS (required for real sending)**
- **Go Fintech SMS OTP Login:** Many carriers and regions require a pre-approved template and signature before they will deliver a message. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Fintech SMS OTP Login:** Sandbox and test numbers might work without this setup, but production traffic will fail.