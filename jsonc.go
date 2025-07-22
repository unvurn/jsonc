package jsonc

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/unvurn/httpc"
)

// NewRequest HTTPリクエストを生成
//
// defaultRespondersを使用してHTTPリクエストを生成します。
func NewRequest[T any]() *httpc.Request[T] {
	return httpc.NewRequest[T]().Decoder("application/json", jsonDecoder[T]).Encoder("application/json", jsonEncoder)
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

func jsonEncoder(params any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}

	return &buf, nil
}
