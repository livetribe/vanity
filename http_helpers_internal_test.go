package vanity

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost_not_forwarded(t *testing.T) {
	r := &http.Request{Host: "host", Header: make(http.Header)}
	h := host(r)
	assert.Equal(t, "host", h)
}

func TestHost_forwarded(t *testing.T) {
	r := &http.Request{Host: "host", Header: http.Header{xForwardedHost: {"forwarded"}}}
	h := host(r)
	assert.Equal(t, "forwarded", h)
}

func TestTemplatize(t *testing.T) {
	expected := `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <meta name="go-import" content="a b c">
  <meta name="go-source" content="a c c/tree/master{/dir} c/blob/master{/dir}/{file}#L{line}">
</head>
</html>
`
	body, err := templatize("a", "b", "c")
	assert.NoError(t, err)
	assert.Equal(t, expected, string(body))
}
