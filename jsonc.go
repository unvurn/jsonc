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
func NewRequest[T any]() httpc.Request[*T] {
	return httpc.NewRequestFunc[*T](defaultResponders[*T]())
}

// jsonResponder JSONレスポンスをデコードするレスポンダー関数
func jsonResponder[T any](res *http.Response) (T, error) {
	defer func() { _ = res.Body.Close() }()
	decoder := json.NewDecoder(res.Body)

	var response T
	err := decoder.Decode(&response)
	if err != nil {
		var zero T
		return zero, err
	}

	return response, nil
}

// defaultResponders JSON対応用デフォルトレスポンダー
//
// この関数は、HTTPレスポンスのContent-Typeヘッダーが"application/json"を含む場合に動作します。
// jsonResponder がジェネリクスによる型変数を必要とするため、変数・シングルトン等による事前用意ではなく
// 関数実行による返値の形式でレスポンダーを生成しています。
func defaultResponders[T any]() []*httpc.ResponderFunc[T] {
	return []*httpc.ResponderFunc[T]{{
		Condition: func(res *http.Response) bool {
			return strings.Contains(res.Header.Get("Content-Type"), "application/json")
		},
		Responder: jsonResponder[T],
	}}
}
