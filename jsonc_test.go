package jsonc_test

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unvurn/httpc/form"

	"github.com/unvurn/jsonc"
)

const httpbinEndpoint = "https://httpbin.org"

type httpbinGetResponse[T any] struct {
	Args    T                 `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

type httpbinPostResponse struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

type httpbinPostFormResponse[T any] struct {
	httpbinPostResponse
	Form  T                 `json:"form"`
	Files map[string]string `json:"files"`
}

type httpbinPostJSONResponse[T any] struct {
	httpbinPostResponse
	JSON T `json:"json"`
}

func TestHttpbin_Get(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "get")
	resp, err := jsonc.NewRequest[httpbinGetResponse[any]]().Get(context.Background(), u)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

type params struct {
	Name         string  `schema:"name" json:"name"`
	Age          int     `schema:"age" json:"age"`
	Scores       []int   `schema:"scores" json:"scores"`
	Description  string  `schema:"description" json:"description"`
	Description2 string  `schema:"description_2,omitempty" json:"description_2,omitempty"`
	ExtraNote    *string `schema:"extra_note" json:"extra_note"`
	ExtraNote2   *string `schema:"extra_note_2,omitempty" json:"extra_note_2,omitempty"`
}

func (p *params) UnmarshalJSON(b []byte) error {
	var temp struct {
		Name         string   `json:"name"`
		Age          string   `json:"age"`
		Scores       []string `json:"scores"`
		Description  string   `json:"description"`
		Description2 string   `json:"description_2,omitempty"`
		ExtraNote    *string  `json:"extra_note"`
		ExtraNote2   *string  `json:"extra_note_2,omitempty"`
	}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	p.Name = temp.Name
	p.Age, err = strconv.Atoi(temp.Age)
	if err != nil {
		return err
	}
	for _, score := range temp.Scores {
		s, err := strconv.Atoi(score)
		if err != nil {
			return err
		}
		p.Scores = append(p.Scores, s)
	}
	p.Description = temp.Description
	p.Description2 = temp.Description2
	if temp.ExtraNote == nil || *temp.ExtraNote == "null" {
		p.ExtraNote = nil
	} else {
		p.ExtraNote = temp.ExtraNote
	}
	if temp.ExtraNote2 == nil || *temp.ExtraNote2 == "null" {
		p.ExtraNote2 = nil
	} else {
		p.ExtraNote2 = temp.ExtraNote2
	}
	return nil
}

func TestHttpbin_PostForm(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostFormResponse[params]]().Post(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.Form.Name)
	assert.Equal(t, 25, resp.Form.Age)
	assert.Equal(t, []int{100, 90, 80}, resp.Form.Scores)
	assert.Len(t, resp.Form.Scores, 3)
	assert.NotEmpty(t, resp.Headers)
	assert.Equal(t, "80", resp.Headers["Content-Length"])
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.URL)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_PostFileUpload(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostFormResponse[params]]().Post(context.Background(), u, p,
		form.Bytes("data1", "data1.txt", []byte("This is data1 content.")),
		form.File("data2", "testdata/samples/dummy.pdf"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.Form.Name)
	assert.Equal(t, 25, resp.Form.Age)
	assert.Len(t, resp.Form.Scores, 3)
	assert.Equal(t, []int{100, 90, 80}, resp.Form.Scores)
	assert.Len(t, resp.Files, 2)
	assert.Equal(t, "This is data1 content.", resp.Files["data1"])
	assert.Len(t, resp.Files["data2"], 17725)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.URL)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}

func TestHttpbin_PostJSON(t *testing.T) {
	u, _ := url.JoinPath(httpbinEndpoint, "post")
	type params struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Scores []int  `json:"scores"`
	}
	p := params{
		Name:   "Jane Doe",
		Age:    25,
		Scores: []int{100, 90, 80},
	}
	resp, err := jsonc.NewRequest[httpbinPostJSONResponse[params]]().PostJSON(context.Background(), u, p)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Jane Doe", resp.JSON.Name)
	assert.Equal(t, 25, resp.JSON.Age)
	assert.Len(t, resp.JSON.Scores, 3)
	assert.Equal(t, []int{100, 90, 80}, resp.JSON.Scores)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Origin)
	assert.Equal(t, u, resp.URL)
	assert.Contains(t, resp.Headers["User-Agent"], "Go-http-client/2.0")
}
