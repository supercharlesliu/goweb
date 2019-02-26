package goweb

import (
	"fmt"
	"testing"
)

func assert(expr bool, failureWords string) {
	if !expr {
		panic(failureWords)
	}
}

func TestParsePathParam(t *testing.T) {
	p0 := ParsePathParam("/sites")
	p1 := ParsePathParam("/sites/")
	p2 := ParsePathParam("/sites/:siteId/")
	p3 := ParsePathParam("/sites/:siteId/plugins/:pluginId")

	// Test path level
	assert(p0.PathDepth == 0, "path level wrong")
	assert(p1.PathDepth == 1, "path level wrong")
	assert(p2.PathDepth == 2, "path level wrong")
	assert(p3.PathDepth == 3, "path level wrong")

	// Test IsPlainPath
	assert(p1.IsPlainPath(), "IsPlainPath wrong")
	assert(!p2.IsPlainPath() && !p3.IsPlainPath(), "IsPlainPath wrong")

	// Test match
	match, params := p0.Match("/sites/1234")
	assert(!match, "match wrong")

	match, params = p1.Match("/sites/1234")
	assert(!match, "match wrong")

	match, params = p2.Match("/sites/1234")
	assert(!match, "match wrong")

	match, params = p2.Match("/sites/1234/")
	assert(match && params["siteId"] == "1234", "match wrong")

	match, params = p3.Match("/sites/1234/")
	assert(!match, "match wrong")

	match, params = p3.Match("/sites/1234/plugins/")
	assert(!match, "match wrong")

	match, params = p3.Match("/sites/1234/plugins/1234/")
	assert(!match, "match wrong")

	match, params = p3.Match("/sites/1234/plugins/1234")
	fmt.Println(params)
	assert(match && params["siteId"] == "1234" && params["pluginId"] == "1234", "match wrong")
}
