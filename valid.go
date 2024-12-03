package nam

type Validator interface {
	Valid() (problems map[string]string)
}
