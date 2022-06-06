package mockhttp

type TestStruct struct {
	Name     string
	Input    *Request
	Expected Response
}
