package goweb

import (
	"regexp"
	"strings"
)

var (
	regexpPath *regexp.Regexp
)

func init() {
	var err error
	regexpPath, err = regexp.Compile("/:([^/]+)")
	if err != nil {
		panic(err)
	}
}

type PathConfig struct {
	BasePath   string
	Domain     string
	Regexp     *regexp.Regexp
	ParamNames []string
	PathDepth  int
}

// PathDepth get depth of a url path
func PathDepth(path string) int {
	return strings.Count(path, "/") - 1
}

// ParsePathParam Get path config info of a url pattern
func ParsePathParam(pattern string) *PathConfig {
	i := strings.Index(pattern, "/")
	domain := ""
	path := pattern
	if i != 0 {
		domain = pattern[0:i]
		path = pattern[i:]
	}
	matches := regexpPath.FindAllStringSubmatch(path, -1)
	pd := PathDepth(path)
	if matches == nil {
		return &PathConfig{
			BasePath:  path,
			Domain:    domain,
			PathDepth: pd,
		}
	}

	paramNames := make([]string, 0)
	for _, item := range matches {
		paramNames = append(paramNames, item[1])
	}
	regexp, err := regexp.Compile("^" + regexpPath.ReplaceAllLiteralString(path, "/([^/]+)") + "$")
	if err != nil {
		panic(err)
	}

	subMatchIndex := regexpPath.FindStringSubmatchIndex(path)
	return &PathConfig{
		BasePath:   path[0 : subMatchIndex[0]+1],
		Domain:     domain,
		Regexp:     regexp,
		ParamNames: paramNames,
		PathDepth:  pd,
	}
}

func (p *PathConfig) PatternString() string {
	return p.Domain + p.BasePath
}

// IsPlainPath test if the path pattern is a plain path(without any regular expression)
func (p *PathConfig) IsPlainPath() bool {
	return p.Regexp == nil
}

// Match test if a url path matches the path pattern
// Return test result as bool in first return value
// The second return value contains captured path params as type map[string]string
func (p *PathConfig) Match(path string) (bool, map[string]string) {
	if p.Regexp == nil {
		return p.BasePath == path, nil
	}
	if !p.Regexp.MatchString(path) {
		return false, nil
	}
	matches := p.Regexp.FindStringSubmatch(path)
	params := make(map[string]string, 0)
	for i, item := range matches[1:] {
		params[p.ParamNames[i]] = item
	}
	return true, params
}
