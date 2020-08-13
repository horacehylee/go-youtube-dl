package decipher

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// jsFunction for decrypt operations function calls
type jsFunction struct {
	Name  string
	Param int
}

var (
	jsFuncPattern = regexp.MustCompile(`[^,()]+`) // not include variable a
)

func parseJsFunction(s string) (jsFunction, error) {
	var f jsFunction
	matches := jsFuncPattern.FindAllStringSubmatch(s, -1)
	if matches == nil || len(matches) < 3 {
		return f, fmt.Errorf("failed to parse JS function with pattern: %v", jsFuncPattern.String())
	}

	ss := strings.Split(matches[0][0], ".")
	if len(ss) < 2 {
		return f, fmt.Errorf("failed to split JS function name")
	}
	f.Name = ss[1]

	p, err := strconv.Atoi(matches[2][0])
	if err != nil {
		return f, err
	}
	f.Param = p
	return f, nil
}
