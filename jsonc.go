package jsonc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/unvurn/httpc"
)

type Request[T any] struct {
	httpc.Request[T]
}

// NewRequest HTTPリクエストを生成
//
// defaultRespondersを使用してHTTPリクエストを生成します。
func NewRequest[T any]() *Request[T] {
	return &Request[T]{*httpc.NewRequestFunc[T](jsonResponder)}
}

// jsonResponder JSONレスポンスをデコードするレスポンダー関数
//
// HTTPレスポンスのContent-Typeヘッダーが"application/json"を含む場合に動作します。
func jsonResponder[T any](res *http.Response) (T, error) {
	var zero T
	if res.StatusCode != http.StatusOK || !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return zero, nil
	}

	decoder := json.NewDecoder(res.Body)

	var response T
	err := decoder.Decode(&response)
	if err != nil {
		return zero, err
	}

	return response, nil
}

func (r *Request[T]) PostJSON(ctx context.Context, u string, params any) (T, error) {
	return r.DoFunc(ctx, http.MethodPost, u, "application/json", func() (io.Reader, error) {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(params); err != nil {
			return nil, err
		}

		return &buf, nil
	})
}
