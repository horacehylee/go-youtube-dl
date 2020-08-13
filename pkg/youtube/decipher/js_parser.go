package decipher

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// jsFunction for decrypt operations function calls
type jsFunction struct {
	name  string
	param int
}

var (
	jsFuncPattern = regexp.MustCompile(`[^,()]+`) // not include variable a
)

func parseJsFunction(s string) (jsFunction, error) {
	var f jsFunction
	matches := jsFuncPattern.FindAllStringSubmatch(s, -1)
	if matches == nil || len(matches) < 3 {
		fmt.Println(matches[1])
		return f, fmt.Errorf("failed to parse JS function with pattern: %v", jsFuncPattern.String())
	}

	ss := strings.Split(matches[0][0], ".")
	if len(ss) < 2 {
		return f, fmt.Errorf("failed to split JS function name")
	}
	f.name = ss[1]

	p, err := strconv.Atoi(matches[2][0])
	if err != nil {
		return f, nil
	}
	f.param = p
	return f, nil
}
