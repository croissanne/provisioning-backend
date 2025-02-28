package ctxval

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	ucontext "github.com/Unleash/unleash-client-go/v3/context"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type commonKeyId int

const (
	loggerCtxKey         commonKeyId = iota
	requestIdCtxKey      commonKeyId = iota
	accountIdCtxKey      commonKeyId = iota
	unleashContextCtxKey commonKeyId = iota
)

var MissingAccountInContextError = errors.New("operation requires account_id in context")

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set.
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(loggerCtxKey) == nil {
		return &log.Logger
	}
	return ctx.Value(loggerCtxKey).(*zerolog.Logger)
}

func WithIdentity(ctx context.Context, id identity.XRHID) context.Context {
	return context.WithValue(ctx, identity.Key, id)
}

func WithIdentityFrom64(ctx context.Context, id string) (context.Context, error) {
	idRaw, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("unable to b64 decode x-rh-identity %w", err)
	}

	var jsonData identity.XRHID
	err = json.Unmarshal(idRaw, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json %w", err)
	}

	return context.WithValue(ctx, identity.Key, jsonData), nil
}

func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// TraceId returns request id or an empty string when not set.
func TraceId(ctx context.Context) string {
	if ctx.Value(requestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(requestIdCtxKey).(string)
}

func WithTraceId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}

// AccountId returns current account model or panics when not set
func AccountId(ctx context.Context) int64 {
	value := ctx.Value(accountIdCtxKey)
	if value == nil {
		panic(MissingAccountInContextError)
	}
	return value.(int64)
}

// AccountIdOrNil returns current account model or 0 when not set.
func AccountIdOrNil(ctx context.Context) int64 {
	value := ctx.Value(accountIdCtxKey)
	if value == nil {
		return 0
	}
	return value.(int64)
}

func WithAccountId(ctx context.Context, accountId int64) context.Context {
	return context.WithValue(ctx, accountIdCtxKey, accountId)
}

// UnleashContext returns unleash context or an empty context when not set.
func UnleashContext(ctx context.Context) ucontext.Context {
	if ctx.Value(unleashContextCtxKey) == nil {
		return ucontext.Context{}
	}
	return ctx.Value(unleashContextCtxKey).(ucontext.Context)
}

func WithUnleashContext(ctx context.Context, uctx ucontext.Context) context.Context {
	return context.WithValue(ctx, unleashContextCtxKey, uctx)
}
