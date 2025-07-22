package jsonc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/unvurn/httpc"
)

type Request[T any] struct {
	httpc.Request[T]
}

// NewRequest HTTPリクエストを生成
//
// defaultRespondersを使用してHTTPリクエストを生成します。
func NewRequest[T any]() *Request[T] {
	return &Request[T]{*httpc.NewRequest[T]().Decoder("application/json", jsonDecoder[T])}
}

// jsonResponder JSONレスポンスをデコードするレスポンダー関数
func jsonDecoder[T any](data []byte) (T, error) {
	var zero T

	var response T
	err := json.Unmarshal(data, &response)
	if err != nil {
		return zero, err
	}

	return response, nil
}

func (r *Request[T]) PostJSON(ctx context.Context, u string, params any) (T, error) {
	var zero T
	response, err := r.DoFunc(ctx, http.MethodPost, u, "application/json", func() (io.Reader, error) {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(params); err != nil {
			return nil, err
		}

		return &buf, nil
	})
	if err != nil {
		return zero, err
	}
	var v T
	err = response.As(&v)
	if err != nil {
		return zero, err
	}
	return v, nil
}
