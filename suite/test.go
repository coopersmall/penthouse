package suite

import (
	"fmt"
	"testing"
)

var (
	currentSuite *suite
	currentCtx   *context
	currentT     *runnter
)

func Suite(name string) *suite {
	currentSuite = newSuite(name)
	return currentSuite
}

func Describe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("Describe must be called with an active suite")
	}

	currentCtx = newContext(name)
	fn()

	currentSuite.tests = append(currentSuite.tests, currentCtx)
	return currentSuite
}

func FDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("FDescribe must be called with an active suite")
	}

	currentCtx = newContext(name)
	currentCtx.focus = true
	fn()

	currentSuite.tests = append(currentSuite.tests, currentCtx)
	return currentSuite
}

func XDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("XDescribe must be called with an active suite")
	}

	currentCtx = newContext(name)
	currentCtx.skip = true
	fn()

	currentSuite.tests = append(currentSuite.tests, currentCtx)
	return currentSuite
}

func Before(fn func()) {
	if currentCtx == nil {
		panic("Before must be called with an active context")
	}

	currentCtx.beforeEach(fn)
}

func JustBefore(fn func()) {
	if currentCtx == nil {
		panic("JustBefore must be called with an active context")
	}

	currentCtx.justBeforeEach(fn)
}

func After(fn func()) {
	if currentCtx == nil {
		panic("After must be called with an active context")
	}

	currentCtx.afterEach(fn)
}

func Context(name string, fn func()) {
	if currentCtx == nil {
		panic("Context must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentCtx.name, name))
	oldCtx := currentCtx
	oldCtx.addChild(ctx)
	currentCtx = ctx
	fn()
	currentCtx = oldCtx
}

func XContext(name string, fn func()) {
	if currentCtx == nil {
		panic("XContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentCtx.name, name))
	ctx.skip = true
	oldCtx := currentCtx
	oldCtx.addChild(ctx)
	currentCtx = ctx
	fn()
	currentCtx = oldCtx
}

func FContext(name string, fn func()) {
	if currentCtx == nil {
		panic("FContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentCtx.name, name))
	ctx.focus = true
	oldCtx := currentCtx
	oldCtx.addChild(ctx)
	currentCtx = ctx
	fn()
	currentCtx = oldCtx
}

func It(name string, fn func(Assert)) {
	if currentCtx == nil {
		panic("It must be called inside a Describe or Context")
	}

	currentCtx.it(name, fn)
}

func XIt(name string, fn func(Assert)) {
	if currentCtx == nil {
		panic("XIt must be called inside a Describe or Context")
	}

	currentCtx.xit(name, fn)
}

func FIt(name string, fn func(Assert)) {
	if currentCtx == nil {
		panic("FIt must be called inside a Describe or Context")
	}

	currentCtx.fit(name, fn)
}

func Run(t *testing.T) {
	run(t, currentSuite)
}
