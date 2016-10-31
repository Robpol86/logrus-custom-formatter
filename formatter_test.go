package lcf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"github.com/stretchr/testify/require"
)

func TestNewFormatterColors(t *testing.T) {
	defer ResetLogger() // Cleanup after test.
	for _, tc := range []string{"no built-in fields", "force colors", "disable colors"} {
		t.Run(tc, func(t *testing.T) {
			assert := require.New(t)

			// Run.
			_, stderr, err := WithCapSys(func() {
				ResetLogger()
				formatter := NewFormatter(Basic, nil)
				formatter.ForceColors = tc == "force colors"
				formatter.DisableColors = tc == "disable colors"
				logrus.SetFormatter(formatter)
				logrus.SetLevel(logrus.DebugLevel)
				LogMsgs(tc != "no built-in fields")
			})
			assert.NoError(err)
			actual := strings.Split(stderr, `\n`)

			// Determine expected from test case.
			var expected []string
			if tc == "no built-in fields" {
				expected = []string{
					"DEBUG::Sample debug 1.",
					"DEBUG::Sample debug 2. a=b c=10",
					"INFO::Sample info 1.",
					"INFO::Sample info 2. a=b c=10",
					"WARNING::Sample warn 1.",
					"WARNING::Sample warn 2. a=b c=10",
					"ERROR::Sample error 1.",
					"ERROR::Sample error 2. a=b c=10",
					"",
				}
			} else if tc == "force colors" {
				expected = []string{
					"\033[37mDEBUG\033[0m:LogMsgs:Sample debug 1.",
					"\033[37mDEBUG\033[0m:LogMsgs:Sample debug 2. \033[37ma\033[0m=b \033[37mc\033[0m=10",
					"\033[34mINFO\033[0m:LogMsgs:Sample info 1.",
					"\033[34mINFO\033[0m:LogMsgs:Sample info 2. \033[34ma\033[0m=b \033[34mc\033[0m=10",
					"\033[33mWARNING\033[0m:LogMsgs:Sample warn 1.",
					"\033[33mWARNING\033[0m:LogMsgs:Sample warn 2. \033[33ma\033[0m=b \033[33mc\033[0m=10",
					"\033[31mERROR\033[0m:LogMsgs:Sample error 1.",
					"\033[31mERROR\033[0m:LogMsgs:Sample error 2. \033[31ma\033[0m=b \033[31mc\033[0m=10",
					"",
				}
			} else {
				expected = []string{
					"DEBUG:LogMsgs:Sample debug 1.",
					"DEBUG:LogMsgs:Sample debug 2. a=b c=10",
					"INFO:LogMsgs:Sample info 1.",
					"INFO:LogMsgs:Sample info 2. a=b c=10",
					"WARNING:LogMsgs:Sample warn 1.",
					"WARNING:LogMsgs:Sample warn 2. a=b c=10",
					"ERROR:LogMsgs:Sample error 1.",
					"ERROR:LogMsgs:Sample error 2. a=b c=10",
					"",
				}
			}

			// Verify.
			assert.Equal(expected, actual)
		})
	}
}

func runFormatterTest(assert *require.Assertions, template string, toFile bool) []string {
	var logFile string
	if toFile {
		tmpdir, err := ioutil.TempDir("", "")
		logFile = filepath.Join(tmpdir, "sample.log")
		assert.NoError(err)
	}

	// Run.
	defer ResetLogger() // Cleanup after test.
	_, stderr, err := WithCapSys(func() {
		ResetLogger()
		logrus.SetFormatter(NewFormatter(template, nil))
		logrus.SetLevel(logrus.DebugLevel)
		if toFile {
			pathMap := lfshook.PathMap{}
			for _, level := range logrus.AllLevels {
				pathMap[level] = logFile
			}
			logrus.AddHook(lfshook.NewHook(pathMap))
			logrus.SetOutput(ioutil.Discard)
		}
		LogMsgs(true)
	})
	assert.NoError(err)

	// Read.
	if toFile {
		assert.Empty(stderr)
		contents, err := ioutil.ReadFile(logFile)
		assert.NoError(err)
		return strings.Split(string(contents), `\n`)
	}
	return strings.Split(stderr, `\n`)
}

func TestNewFormatterBasic(t *testing.T) {
	for _, toFile := range []bool{false, true} {
		t.Run(fmt.Sprintf("toFile:%v", toFile), func(t *testing.T) {
			assert := require.New(t)
			actual := runFormatterTest(assert, Basic, toFile)
			expected := []string{
				"DEBUG:LogMsgs:Sample debug 1.",
				"DEBUG:LogMsgs:Sample debug 2. a=b c=10",
				"INFO:LogMsgs:Sample info 1.",
				"INFO:LogMsgs:Sample info 2. a=b c=10",
				"WARNING:LogMsgs:Sample warn 1.",
				"WARNING:LogMsgs:Sample warn 2. a=b c=10",
				"ERROR:LogMsgs:Sample error 1.",
				"ERROR:LogMsgs:Sample error 2. a=b c=10",
				"",
			}
			assert.Equal(expected, actual)
		})
	}
}

func TestNewFormatterMessage(t *testing.T) {
	for _, toFile := range []bool{false, true} {
		t.Run(fmt.Sprintf("toFile:%v", toFile), func(t *testing.T) {
			assert := require.New(t)
			actual := runFormatterTest(assert, Message, toFile)
			expected := []string{
				"Sample debug 1.",
				"Sample debug 2. a=b c=10",
				"Sample info 1.",
				"Sample info 2. a=b c=10",
				"Sample warn 1.",
				"Sample warn 2. a=b c=10",
				"Sample error 1.",
				"Sample error 2. a=b c=10",
				"",
			}
			assert.Equal(expected, actual)
		})
	}
}

func TestNewFormatterDetailed(t *testing.T) {
	reTimestamp := regexp.MustCompile(`^\d{4}-\d\d-\d\d \d\d:\d\d:\d\d\.\d{3}`)

	for _, toFile := range []bool{false, true} {
		t.Run(fmt.Sprintf("toFile:%v", toFile), func(t *testing.T) {
			assert := require.New(t)
			actual := runFormatterTest(assert, Detailed, toFile)
			for i, str := range actual {
				if str != "" {
					actual[i] = reTimestamp.ReplaceAllString(str, "2016-10-30 19:12:17.149")
				}
			}
			expected := []string{
				"2016-10-30 19:12:17.149 %s DEBUG LogMsgs              Sample debug 1.",
				"2016-10-30 19:12:17.149 %s DEBUG LogMsgs              Sample debug 2. a=b c=10",
				"2016-10-30 19:12:17.149 %s INFO  LogMsgs              Sample info 1.",
				"2016-10-30 19:12:17.149 %s INFO  LogMsgs              Sample info 2. a=b c=10",
				"2016-10-30 19:12:17.149 %s WARN  LogMsgs              Sample warn 1.",
				"2016-10-30 19:12:17.149 %s WARN  LogMsgs              Sample warn 2. a=b c=10",
				"2016-10-30 19:12:17.149 %s ERROR LogMsgs              Sample error 1.",
				"2016-10-30 19:12:17.149 %s ERROR LogMsgs              Sample error 2. a=b c=10",
				"",
			}
			for i, str := range expected {
				if str != "" {
					expected[i] = fmt.Sprintf(str, fmt.Sprintf("%-5d", os.Getpid()))
				}
			}
			assert.Equal(expected, actual)
		})
	}
}
