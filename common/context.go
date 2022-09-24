package common

import "context"

const (
	ContextKeyRequestID = "request_id"
)

// GetRequestIdFromCtx returns the string request id from context
func GetRequestIdFromCtx(ctx context.Context) string {
	reqIDRaw := ctx.Value(ContextKeyRequestID)
	reqID, ok := reqIDRaw.(string)
	if !ok {
		return ""
	}

	return reqID
}
