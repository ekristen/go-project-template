package server

type customErr struct {
	Message string                 `json:"msg"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// anotherErr is another custom error.
type anotherErr struct {
	Foo int `json:"foo"`
}

func (anotherErr) Error() string {
	return "foo happened"
}
