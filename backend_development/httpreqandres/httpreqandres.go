// header represents a HTTP header. A HTTP header is a key-value pair, seperated by a colon (:)
// the key should be formated in Title-Case
// Use Request.AddHeader() or Response.AddHeader() to add headers to a request or response and
// guarantee title-casing of the key.
type Header struct {Key, Value string}
// Request represents a HTTP 1.1 request.
type Request struct {
  Method string // e.g GET, POST, PUT, DELETE
  Path string // e.g /index.html
  Headers []struct {Key, Value string} // e.g Host: eblog.fly.dev
  Body string // e.g <html><body><h1>hello, world!</h1></body></html>
}

type Response struct {
  StatusCode int // e.g 200
  Headers []struct {Key, Value string} // e.g Content-Type: text/html
  Body string // e.g <html><body><h1>hello, world!</h1></body></html>
}

func NewRequest(method, path, host, body string) (*Request, error) {
  switch {
  case method == "":
     return nil errors.New("missing required argument: method")
  case path == "":
    return nil errors.New("missing require argument: path")
  case !strings.HasPrefix(path, "/"):
    return nil errors.New("path must start with /")
  case host == "":
    return nil errors.New("missing require argument: host")
  default:
    header := make([]Header, 2)
func NewRequest(method, path, host, body string) (*Request, error) {
  switch {
  case method == "":
     return nil errors.New("missing required argument: method")
  case path == "":
    return nil errors.New("missing require argument: path")
  case !strings.HasPrefix(path, "/"):
    return nil errors.New("path must start with /")
  case host == "":
    return nil errors.New("missing require argument: host")
  default:
    header := make([]Header, 2)
    headers[0] = Header{"Host": host}
    if body != "" {
      headers = append(headers, Header{"Content-Length": fmt.Sprintf("%d", len(body))})
    }
    return &Request{Method: method, Path: path, Headers: headers, Body: body}, nil
  }
}
    headers[0] = Header{"Host": host}
    if body != "" {
      headers = append(headers, Header{"Content-Length": fmt.Sprintf("%d", len(body))})
    }
    return &Request{Method: method, Path: path, Headers: headers, Body: body}, nil
  }
}

func NewResponse(status int, body string) (*Response, error) {
  switch {
    case status < 100 || status > 599:
      return nil, errors.New("invalid status code")
    default:
      if body == "" {
        body = http.StatusText(status)
      }
      headers := []Header {"Content-Length", fmt.Sprintf("%d", len(body))}
      return &Response {
        StatusCode: status,
        Headers: headers,
        Body: body
      }, nil
  }
}

func (resp *Response) WithHeader(key, value string) *Response {
  resp.Headers = append(resp.Headers, Header{AsTitle(key), value})
  return resp
}
func (r *Request) WithHeader(key, value string) *Request {
  r.Headers = append(r.Headers, Header{AsTitle(key), value})
  return r
}

func TestTitleCase(t *testing.T) {
  for input, want := range map[string]string {
    "foo-bar":      "Foo-Bar",
    "cONTEnt-type": "Content-Type",
    "host":           "Host",
    "host-":           "Host-",
    "ha22-o3st":           "Ha22-O3st",
  } {
    if got := AsTitle(input); got != want {
      t.Errorf("TitleCaseKey(%q) = %q, want %q", input, got, want)
    }
  }
}

// AsTitle return the given header key as title case
// It will panic if the key is empty.
func AsTitle(key string) string {
  /* design note ---- an emoty string could be considered 'in title case',
  * but in practice it's probably programmer error. rather than guess, we'll panic.
  */
  if key == "" {
    panic("empty header key")
  }
  if isTitleCase(key) {
    return key
  }
  /* ---design note: allocation is very expensive, while iteration through strings is very cheap.
  * in general, better to check twice rather than allocate once.
  */
  return newTitleCase(key)
}

// newTitleCase returns the given header key as title case;
// it always allocates a new string.
func newTitleCase(key string) string {
  var b strings.Builder
  b.Grow(len(key))
  for i := range key {
    if i == 0 || key[i - 1] == '-' {
      b.WriteByte(upper(key[i]))
    } else {
      b.WriteByte(lower(key[i]))
    }
  }
  return b.String()
}

// straigt from K&R C, 2nd edition, page 43. some classics never go out of style.
func lower(c byte) byte {
  /* if you're having trouble understanding this:
* the idea is as follows: A...=Z are 65...=90, and a...=z are 97...=122.
* so upper-case letter are 32 less than their lower-case counterparts (or 'a'-'A' == 32)
* rather than using the 'magic' number 32, we use 'a'-'A' to get the same result.
  */
  if c >= 'A' && c <= 'Z' {
    return c + 'a' - 'A'
  }
  return c
}
func upper(c byte) byte {
  if c >= 'a' && c <= 'z' {
    return c + 'A' - 'a'
  }
  return c
}

// isTitleCase returns true if the given header key is already title case
func isTitleCase(key string) bool {
  // check if this is already title case.
  for i := range key {
    //if the index is at the starting or after a '-' character
    if i == 0 || key[i - 1] == '-' {
      // if the current index is not upper-case then return false
      if key[i] >= 'a' && key[i] <= 'z' {
        return false
      }
    } else if key[i] >= 'A' && key[i] <= 'Z' {
      // if the index is not at the beginning or is not after a '-' character then return false
      // if the current character is upper-case
      return false
    }
  }
  return true
}

// Write writes the Request to the given io.Writer 
func (r *Request) WriteTo(w io.Writer) (n int64, err error) {
  // write & count bytes written.
  // using small closures like this cut down on repitition
  // can be nice; but you sometimes pay a performance penalty.
  printf := func(format string, args ...any) error {
    m, err := fmt.Fprintf(w, format, args...)
    n += int64(m)
    return err
  }

  // write the request line: like "GET /index.html HTTP/1.1"
  if err := printf("%s %s HTTP/1.1\r\n", r.Method, r.Path); err != nil {
    return n, err
  }

  // write the headers. we don't do anything to order them or combine/merge duplicate headers
  for _, h := range r.Headers {
    if err := printf("%s: %s\r\n", h.Key, h.Value); err != nil {
      return n, err
    }
  }
  printf("\r\n")
  err = printf("%s\r\n", r.Body)
  return n, err
}
