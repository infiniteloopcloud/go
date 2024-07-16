package middlewares

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var fileDetails = regexp.MustCompile(`^(/([\w-._]+/)*([\w-]+\.[a-z]+)):(\d+)\s?(\+0x[0-9a-fA-F]+$)?`)
var tabsNewlinesRegexp = regexp.MustCompile(`[\t\n]+`)

type stackFlags struct {
	Vendor  bool
	Builtin bool
}
type stackLine struct {
	FunctionName  string
	FilePath      string
	FilePathShort string
	LineNumber    int
	StackPosition string
	Flags         stackFlags
}

func (sl stackLine) String() string {
	return fmt.Sprintf("%s:%d:%s", sl.FilePathShort, sl.LineNumber, sl.FunctionName)
}

type panicParser struct{}

func (p panicParser) Parse(stack []byte) ([]stackLine, error) {
	str := strings.TrimSpace(string(stack))
	// cut the header
	_, str, _ = strings.Cut(str, "\n")
	// replace new lines with tabs
	str = strings.ReplaceAll(str, "\n", "\t")
	// replace tab duplications with single tabs
	str = tabsNewlinesRegexp.ReplaceAllString(str, "\t")

	// nolint: prealloc
	var result []stackLine
	var tabCounter int
	var latestTabIndex int
	for i, r := range []rune(str) {
		if r != '\t' {
			// skip chars until reaching a tab
			continue
		}
		tabCounter++
		if tabCounter%2 != 0 {
			// every second tab matters, odd tabs will be skipped
			continue
		}

		// take the part of the string between the even tabs
		line := p.substring(str, latestTabIndex, i)

		functionNameParams, fileData, found := strings.Cut(line, "\t")
		if !found {
			return nil, fmt.Errorf("error separating function name from file details: %s", line)
		}

		// remove last (...)
		functionName, _, found := p.cutLast(functionNameParams, "(")
		if !found {
			return nil, fmt.Errorf("error cutting params from the func definition: %s", functionNameParams)
		}

		fileDetailsRegexResult := fileDetails.FindAllStringSubmatch(fileData, -1)
		if len(fileDetailsRegexResult) == 0 || len(fileDetailsRegexResult[0]) != 6 {
			return nil, fmt.Errorf("error parsing file details of the stack line: %s", fileData)
		}

		lineNum, err := strconv.Atoi(fileDetailsRegexResult[0][4])
		if err != nil {
			return nil, fmt.Errorf("error parsing line number: %w", err)
		}

		groot := runtime.GOROOT()
		result = append(result, stackLine{
			FunctionName:  functionName + "(...)",
			FilePath:      fileDetailsRegexResult[0][1],
			FilePathShort: fileDetailsRegexResult[0][2] + fileDetailsRegexResult[0][3],
			LineNumber:    lineNum,
			StackPosition: fileDetailsRegexResult[0][5],
			Flags: stackFlags{
				Vendor:  strings.Contains(fileDetailsRegexResult[0][1], "/vendor/"),
				Builtin: strings.Contains(fileDetailsRegexResult[0][1], groot),
			},
		})
		latestTabIndex = i
	}

	return result, nil
}

func (panicParser) substring(str string, i, j int) string {
	if i > 0 {
		i += 1
	}
	return str[i:j]
}

func (panicParser) cutLast(str, sep string) (string, string, bool) {
	idx := strings.LastIndex(str, sep)
	if idx < 0 {
		return str, "", false
	}
	return str[:idx], str[idx+1:], true
}
