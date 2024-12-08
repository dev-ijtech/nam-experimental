package nam

import "fmt"

type Validator interface {
	Valid() ProblemSet
}

type Problem struct {
	Name   string
	Reason string
}

type ProblemSet struct {
	Set []Problem
}

func (p *ProblemSet) Add(name string, reason string) {
	p.Set = append(p.Set, Problem{name, reason})
}

func (p *ProblemSet) String() string {
	problemString := ""

	for i, problem := range p.Set {
		problemString += fmt.Sprintf("%d. %s %s\n", i+1, problem.Name, problem.Reason)
	}

	return problemString
}
