package ginutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLocaleFromAcceptLanguage(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "empty", header: "", want: LocaleEN},
		{name: "en first", header: "en-US,en;q=0.9", want: LocaleEN},
		{name: "zh first", header: "zh-CN,zh;q=0.9,en;q=0.8", want: LocaleZH},
		{name: "weighted preference zh", header: "en;q=0.1, zh;q=0.9", want: LocaleZH},
		{name: "weighted preference en", header: "zh;q=0.1, en;q=0.9", want: LocaleEN},
		{name: "q zero ignored", header: "zh;q=0, en;q=0.8", want: LocaleEN},
		{name: "unknown fallback en", header: "fr-FR,ja;q=0.9", want: LocaleEN},
		{name: "underscore", header: "zh_CN", want: LocaleZH},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LocaleFromAcceptLanguage(tt.header)
			if got != tt.want {
				t.Fatalf("LocaleFromAcceptLanguage(%q)=%q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

func TestBodyForAcceptLanguage_LocalizedFallback(t *testing.T) {
	err := NewLocalizedBusinessGinError(194, "no permission", map[string]string{"zh": "沒有權限"})

	if got, want := err.BodyForAcceptLanguage("").Message, "no permission"; got != want {
		t.Fatalf("default body message=%q, want %q", got, want)
	}
	if got, want := err.BodyForAcceptLanguage("zh-TW,zh;q=0.9").Message, "沒有權限"; got != want {
		t.Fatalf("zh body message=%q, want %q", got, want)
	}
}

func TestLocalizedConstructor_ZhHant(t *testing.T) {
	err := NewLocalizedBusinessGinError(194, "no permission", map[string]string{"zh-Hant": "沒有權限"})

	body := err.BodyForAcceptLanguage("zh-Hant")
	if got, want := body.Message, "沒有權限"; got != want {
		t.Fatalf("zh-Hant body message=%q, want %q", got, want)
	}
}

func TestWithMessage_LocalizesEachSegment(t *testing.T) {
	err := NewLocalizedBadRequestGinError(142, "invalid input", map[string]string{"zh": "無效輸入"}).
		WithMessage("fiatType is required", map[string]string{"zh": "fiatType 為必填項"})

	if got, want := err.Body().Message, "invalid input: fiatType is required"; got != want {
		t.Fatalf("raw message=%q, want %q", got, want)
	}
	if got, want := err.BodyForAcceptLanguage("zh-Hant").Message, "無效輸入: fiatType 為必填項"; got != want {
		t.Fatalf("localized message=%q, want %q", got, want)
	}
}

func TestRenderError_UsesAcceptLanguage(t *testing.T) {
	ctx, recorder := buildGinErrorTestContext("zh-TW")

	RenderError(ctx, NewLocalizedBusinessGinError(194, "no permission", map[string]string{"zh": "沒有權限"}))

	if got, want := recorder.Code, HTTP_STATUS_BUISINESS; got != want {
		t.Fatalf("status=%d, want %d", got, want)
	}
	var body GinErrorBody
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got, want := body.Message, "沒有權限"; got != want {
		t.Fatalf("message=%q, want %q", got, want)
	}
}

func TestRenderResponse_UsesAcceptLanguageForGinError(t *testing.T) {
	ctx, recorder := buildGinErrorTestContext("zh-CN,zh;q=0.9")

	RenderResponse(ctx, nil, NewLocalizedBusinessGinError(194, "no permission", map[string]string{"zh": "沒有權限"}))

	if got, want := recorder.Code, HTTP_STATUS_BUISINESS; got != want {
		t.Fatalf("status=%d, want %d", got, want)
	}
	var body GinErrorBody
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got, want := body.Message, "沒有權限"; got != want {
		t.Fatalf("message=%q, want %q", got, want)
	}
}

func TestRenderResponse_PreservesGenericNormalMessage(t *testing.T) {
	ctx, recorder := buildGinErrorTestContext("zh-TW")

	RenderResponse(ctx, nil, errors.New("unknown BANK_TYPE: foo"))

	if got, want := recorder.Code, HTTP_STATUS_BUISINESS; got != want {
		t.Fatalf("status=%d, want %d", got, want)
	}
	var body GinErrorBody
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got, want := body.Code, ERR_CORDE_NORMAL; got != want {
		t.Fatalf("code=%d, want %d", got, want)
	}
	if got, want := body.Message, "unknown BANK_TYPE: foo"; got != want {
		t.Fatalf("message=%q, want %q", got, want)
	}
}

func buildGinErrorTestContext(acceptLanguage string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if acceptLanguage != "" {
		req.Header.Set("Accept-Language", acceptLanguage)
	}
	ctx.Request = req
	return ctx, recorder
}
