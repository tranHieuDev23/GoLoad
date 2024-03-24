package servemuxoptions

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func WithAuthCookieToAuthMetadata(authCookieName, authMetadataName string) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			return make(metadata.MD)
		}

		return metadata.New(map[string]string{
			authMetadataName: cookie.Value,
		})
	})
}

func WithAuthMetadataToAuthCookie(
	authMetadataName,
	authCookieName string,
	expiresInDuration time.Duration,
) runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
		metadata, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		authMetadataValues := metadata.HeaderMD.Get(authMetadataName)
		if len(authMetadataValues) == 0 {
			return nil
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    authMetadataValues[0],
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(expiresInDuration),
		})
		return nil
	})
}

func WithRemoveGoAuthMetadata(authMetadataName string) runtime.ServeMuxOption {
	return runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
		if s == authMetadataName {
			return "", false
		}

		return runtime.DefaultHeaderMatcher(s)
	})
}
