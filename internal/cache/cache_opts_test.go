package cache

import (
	"testing"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	h "github.com/buildpacks/pack/testhelpers"
)

type CacheOptTestCase struct {
	name       string
	input      string
	output     string
	shouldFail bool
}

func TestMetadata(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "Metadata", testCacheOpts, spec.Sequential(), spec.Report(report.Terminal{}))
}

func testCacheOpts(t *testing.T, when spec.G, it spec.S) {
	when("image cache format options are passed", func() {
		it("with complete options", func() {
			testcases := []CacheOptTestCase{
				{
					name:   "Build cache as Image",
					input:  "type=build;format=image;name=io.test.io/myorg/my-cache:build",
					output: "type=build;format=image;name=io.test.io/myorg/my-cache:build;type=launch;format=volume;name=;",
				},
				{
					name:   "Launch cache as Image",
					input:  "type=launch;format=image;name=io.test.io/myorg/my-cache:build",
					output: "type=build;format=volume;name=;type=launch;format=image;name=io.test.io/myorg/my-cache:build;",
				},
			}

			for _, testcase := range testcases {
				var cacheFlags CacheOpts
				t.Logf("Testing cache type: %s", testcase.name)
				err := cacheFlags.Set(testcase.input)
				h.AssertNil(t, err)
				h.AssertEq(t, testcase.output, cacheFlags.String())
			}
		})

		it("with missing options", func() {
			successTestCases := []CacheOptTestCase{
				{
					name:   "Build cache as Image missing: type",
					input:  "format=image;name=io.test.io/myorg/my-cache:build",
					output: "type=build;format=image;name=io.test.io/myorg/my-cache:build;type=launch;format=volume;name=;",
				},
				{
					name:   "Build cache as Image missing: format",
					input:  "type=build;name=io.test.io/myorg/my-cache:build",
					output: "type=build;format=volume;name=io.test.io/myorg/my-cache:build;type=launch;format=volume;name=;",
				},
				{
					name:       "Build cache as Image missing: name",
					input:      "type=build;format=image",
					output:     "cache 'name' is required",
					shouldFail: true,
				},
				{
					name:   "Build cache as Image missing: type, format",
					input:  "name=io.test.io/myorg/my-cache:build",
					output: "type=build;format=volume;name=io.test.io/myorg/my-cache:build;type=launch;format=volume;name=;",
				},
				{
					name:   "Build cache as Image missing: format, name",
					input:  "type=build",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=;",
				},
				{
					name:       "Build cache as Image missing: type, name",
					input:      "format=image",
					output:     "cache 'name' is required",
					shouldFail: true,
				},
				{
					name:       "Launch cache as Image missing: name",
					input:      "type=launch;format=image",
					output:     "cache 'name' is required",
					shouldFail: true,
				},
			}

			for _, testcase := range successTestCases {
				var cacheFlags CacheOpts
				t.Logf("Testing cache type: %s", testcase.name)
				if testcase.name == "Everything missing" {
					print("i am here")
				}
				err := cacheFlags.Set(testcase.input)

				if testcase.shouldFail {
					h.AssertError(t, err, testcase.output)
				} else {
					h.AssertNil(t, err)
					output := cacheFlags.String()
					h.AssertEq(t, testcase.output, output)
				}
			}
		})

		it("with invalid options", func() {
			testcases := []CacheOptTestCase{
				{
					name:       "Invalid cache type",
					input:      "type=invalid_cache;format=image;name=io.test.io/myorg/my-cache:build",
					output:     "invalid cache type 'invalid_cache'",
					shouldFail: true,
				},
				{
					name:       "Invalid cache format",
					input:      "type=launch;format=invalid_format;name=io.test.io/myorg/my-cache:build",
					output:     "invalid cache format 'invalid_format'",
					shouldFail: true,
				},
				{
					name:       "Not a key=value pair",
					input:      "launch;format=image;name=io.test.io/myorg/my-cache:build",
					output:     "invalid field 'launch' must be a key=value pair",
					shouldFail: true,
				},
				{
					name:       "Extra semicolon",
					input:      "type=launch;format=image;name=io.test.io/myorg/my-cache:build;",
					output:     "invalid field '' must be a key=value pair",
					shouldFail: true,
				},
			}

			for _, testcase := range testcases {
				var cacheFlags CacheOpts
				t.Logf("Testing cache type: %s", testcase.name)
				err := cacheFlags.Set(testcase.input)
				h.AssertError(t, err, testcase.output)
			}
		})
	})

	when("volume cache format options are passed", func() {
		it("with complete options", func() {
			testcases := []CacheOptTestCase{
				{
					name:   "Build cache as Volume",
					input:  "type=build;format=volume;name=test-build-volume-cache",
					output: "type=build;format=volume;name=test-build-volume-cache;type=launch;format=volume;name=;",
				},
				{
					name:   "Launch cache as Volume",
					input:  "type=launch;format=volume;name=test-launch-volume-cache",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=test-launch-volume-cache;",
				},
			}

			for _, testcase := range testcases {
				var cacheFlags CacheOpts
				t.Logf("Testing cache type: %s", testcase.name)
				err := cacheFlags.Set(testcase.input)
				h.AssertNil(t, err)
				h.AssertEq(t, testcase.output, cacheFlags.String())
			}
		})

		it("with missing options", func() {
			successTestCases := []CacheOptTestCase{
				{
					name:   "Launch cache as Volume missing: format",
					input:  "type=launch;name=test-launch-volume",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=test-launch-volume;",
				},
				{
					name:   "Launch cache as Volume missing: name",
					input:  "type=launch;format=volume",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=;",
				},
				{
					name:   "Launch cache as Volume missing: format, name",
					input:  "type=launch",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=;",
				},
				{
					name:   "Launch cache as Volume missing: type, name",
					input:  "format=volume",
					output: "type=build;format=volume;name=;type=launch;format=volume;name=;",
				},
			}

			for _, testcase := range successTestCases {
				var cacheFlags CacheOpts
				t.Logf("Testing cache type: %s", testcase.name)
				err := cacheFlags.Set(testcase.input)

				if testcase.shouldFail {
					h.AssertError(t, err, testcase.output)
				} else {
					h.AssertNil(t, err)
					output := cacheFlags.String()
					h.AssertEq(t, testcase.output, output)
				}
			}
		})
	})
}
