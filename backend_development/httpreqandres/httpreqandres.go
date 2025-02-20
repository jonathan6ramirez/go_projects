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

func (resp *Response) WriteTo(w io.Writer) (n int64, err error) {
  printf := func(format string, args, ...any) error {
    m, err := fmt.Fprintf(w, format, args...)
    n += int64(m)
    return err
  }
  if err := printf("HTTP/1.1 %d %s\r\n", resp.StatusCode, http.StatusText(resp.StatusCode)); err != nil {
    return n, err
  }

  for _, h := range resp.Headers {
    if err := printf("%s: %s\r\n", h.Key, h.Value); err != nil {
      return n, err
    }
  }

  if err := printf("\r\n%s\r\n", resp.Body); err != nil {
    return n, err
  }

  return n, nil
}

var _, _ fmt.Stringer = (*Request)(nil), (*Response)(nil) // compile-time check that Request and Response implement fmt.Stringer
var _, _ encoding.TextMarshaler = (*Request)(nil), (*Response)(nil)
func (r *Request) String() string { b := new(strings.Builder); r.WriteTo(b); return b.String() }
func (resp *Response) String() string { b := new(strings.Builder); resp.WriteTo(b); return b.String() }
func (r *Request) MarshalText() ([]byte, error) { b := new(bytes.Buffer); r.WriteTo(b); return b.Bytes(), nil }
func (resp *Response) MarshalText() ([]byte, error) { b := new(bytes.Buffer); resp.WriteTo(b); return b.Bytes(), nil }

// ParseRequest parses a HTTP request from the given text.
func ParseRequest(raw string) (r Request, err error) {
  // request has three parts:
  // 1. Request lines
  // 2. Headers
  // 3. Body (optional)
  lines := splitLines(raw)

  log.Println(lines)
  if len(lines) < 3 {
    return Request{}, fmt.Errorf("malformed request: should have at least 3 lines")
  }

  // The first line is special.
  first := strings.Fields(lines[0])
  r.Method, r.Path = first[0], first[1]
  if !strings.HasPrefix(r.Path, "/") {
    return Request{}, fmt.Errorf("malformed request: path should start with /")
  }
  if !strings.Contains(first[2], "HTTP") {
    return Request{}, fmt.Errorf("malformed request: first line should contain HTTP version")
  }
  var foundhost bool
  var bodyStart int
  // then we have Headers. up until an empty line.
  for i := 1; i < len(lines); i++ {
    if lines[i] == "" { // this means an empty line was found
      bodyStart = i + 1
      break
    }
    key, val, ok := strings.Cut(lines[i], ": ")
    if !ok {
      return Request{}, fmt.Errorf("malfromed request: header %q should be of form 'key: value'", lines[i])
    }
    if key == "Host" { // special case: host header is not required.
      foundhost = true
    }
    key = AsTitle(key)

    r.Headers = append(r.Headers, Header{key, val})
  }
  end := len(lines) -1 // recombine the body using normal newlines; skip the last empty line.
  r.Body = strings.Join(lines[bodyStart: end], "\r\n")
  if !foundhost {
    return Request{}, fmt.Errorf("malformed request: missing Host header")
  }

  return r, nil
}

/// ParseResponse parses the given HTTP/1.1 response string into the Response. It returns an error
// if the response is invalid,
// - not a valid integer
// - invalid status code
// - missing status text
// - invalid headers
// it doesn't properly handle multi-line headers, headers with multiple values, or html-encoding, etc.
func ParseResponse(raw string) (resp *Response, err error) {
  // response has three parts:
  // 1. Response line
  // 2. Headers
  // 3. Body (optional)
  lines := splitLines(raw)
  log.Println(lines)

  // The first line is special.
  first := strings.SplitN(lines[0], " ", 3)
  if !strings.Contains(first[0], "HTTP") {
    return nil, fmt.Errorf("malformed Response: first line should contain HTTP version")
  }
  resp = new(Response)
  resp.StatusCode, err = strconv.Atoi(first[1])
  if err != nil {
    return nil, fmt.Errorf("malformed response: expected status code to be an integer, got %q", first[1])
  }
  if first[2] == "" || http.StatusText(resp.StatusCode) != first[2] {
    log.Printf("missing or incorrect status text for status code %d: expected %q, but got %q", resp.StatusCode, http.StatusText(resp.StatusCode), first[2])
  }
  var bodyStart int
  // then we have headers, up until an empty line.
  for i := 1; i < len(lines); i++ {
    log.Println(i, lines[i])
    if lines[i] == "" {
      bodyStart = i + 1
      break
    }
    key, val, ok := strings.Cut(lines[i], ": ")
    if !ok {
      return nil, fmt.Errorf("malformed response: header %q should be of form 'key: value'", lines[i])
    }
    key = AsTitle(key)
    resp.Headers = append(resp.Headers, Header{key, val})
  }
  resp.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\r\n")) // recombine the body using normal newlines.
  return resp, nil
}

// splitLines on the "\r\n" sequence; multiple separators in a row are NOT collapsed.
func splitLines(s string) []string {
  if s == ""  {
    return nil
  }
  var lines []string
  i := 0
  for {
    j := strings.Index(s[i:], "\r\n")
    if j == -1 {
      lines = append(lines, s[i:])
      return lines
    }
    lines = append(lines, s[i:i+j]) // up to but not including the \r\n
    i += j + 2 // skip the \r\n
  }
}

func TestHTTPResponse(t *testing.T) {
  for name, tt := range map[string]struct {
    input string
    want *Response
  }{
    "200 OK (no body)": {
      input: "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n",
      want: &Response{
        StatusCode: 200,
        Headers: []Header{
          {"Content-Length", "0"},
        }
      }
    },
    "404 Not Found (w/ body)": {
      input: "HTTP/1.1 404 Not Found\r\nContent-Length: 11\r\n\r\nHello World\r\n",
      want: &Response{
        StatusCode: 404,
        Headers: []Header{
          {"Content-Length", "11"},
        },
        Body: "Hello World",
      }
    },
  } {
    t.Run(name, func(t *testing.T) {
      got, err := ParseResponse(tt.input)
      if err != nil {
        t.Errorf("ParseResponse(%q) returned error: %v", tt.input, err)
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("ParseResponse(%q) = %#+v, want %#+v", tt.input, got, tt.want)
      }

      if got2, err := ParseResponse(got.String()); err != nil {
        t.Errorf("ParseResponse(%q) returned error: %v", got.String(), err)
      } else if !reflect.DeepEqual(got2, got) {
        t.Errorf("ParseResponse(%q) = %#+v, want %#+v", got.String(), got2, got)
      }
    })
  }
}

func TestHTTPRequest(t *testing.T) {
    for name, tt := range map[string]struct {
        input string
        want  Request
    }{
        "GET (no body)": {
            input: "GET / HTTP/1.1\r\nHost: www.example.com\r\n\r\n",
            want: Request{
                Method: "GET",
                Path:   "/",
                Headers: []Header{
                    {"Host", "www.example.com"},
                },
            },
        },
        "POST (w/ body)": {
            input: "POST / HTTP/1.1\r\nHost: www.example.com\r\nContent-Length: 11\r\n\r\nHello World\r\n",
            want: Request{
                Method: "POST",
                Path:   "/",
                Headers: []Header{
                    {"Host", "www.example.com"},
                    {"Content-Length", "11"},
                },
                Body: "Hello World",
            },
        },
    } {
        t.Run(name, func(t *testing.T) {
            got, err := ParseRequest(tt.input)
            if err != nil {
                t.Errorf("ParseRequest(%q) returned error: %v", tt.input, err)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseRequest(%q) = %#+v, want %#+v", tt.input, got, tt.want)
            }
            // test that the request can be written to a string and parsed back into the same request.
            got2, err := ParseRequest(got.String())
            if err != nil {
                t.Errorf("ParseRequest(%q) returned error: %v", got.String(), err)
            }
            if !reflect.DeepEqual(got, got2) {
                t.Errorf("ParseRequest(%q) = %+v, want %+v", got.String(), got2, got)
            }

        })
    }
}
