package functions

import (
	"net/http"

	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
)

func GetCtxRequestID(r *http.Request) (string, error) {

	requestID, ok := r.Context().Value(mc.AccessKey(mc.RequestID)).(string)
	if !ok {
		// we need to authenticate requests using unique keys | remote address is OK
		return r.RemoteAddr, me.ErrNoRequestIdInContext
	}
	return requestID, nil
}
