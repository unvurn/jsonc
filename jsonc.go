package jsonc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/unvurn/httpc"
)

// NewRequest HTTPリクエストを生成
//
// defaultRespondersを使用してHTTPリクエストを生成します。
func NewRequest[T any]() *httpc.Request[T] {
	return httpc.NewRequestFunc[T](jsonResponder)
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
