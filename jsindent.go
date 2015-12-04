package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	for _, v := range stringChars {
		stringCharsMap[string(v)] = struct{}{}
	}

	for _, v := range openChars {
		openCharsMap[string(v)] = struct{}{}
	}

	for _, v := range closeChars {
		closeCharsMap[string(v)] = struct{}{}
	}

	for _, v := range operators {
		operatorsMap[v] = struct{}{}
	}

	for _, v := range parensKeywords {
		parensKeywordsMap[v] = struct{}{}
	}

	for _, v := range keywords {
		keywordsMap[v] = struct{}{}
	}
}

type parseState int

//go:generate stringer -type=parseState
const (
	stateNone parseState = iota
	stateIdentifier
	stateOperator
	stateStringSingle
	stateStringDouble
	stateStringEscape
	stateBrace
	stateParen
	stateArray
	stateCommentStart // is a /
	stateCommentEnd   // is a *
	stateComment      // is a /
	stateCommentBlock // is a *
)

const spacesPerIndent = 2

func doIndent(inputStr string) string {
	outBuf := &bytes.Buffer{}
	identBuf := &bytes.Buffer{}
	currentLine := &bytes.Buffer{}
	var lastLine *bytes.Buffer
	input := []byte(inputStr)
	state := stateNone
	stack := NewStack()

	indentLevel := 0
	currentOpenCounter := 0
	lastOpenCounter := 0

	sawUglyKeyword := false
	lastLineSawUglyKeyword := false
	indentNextLine := false

	foundDo := false
	parensKeywordFound := false
	lastLineParensKeywordFound := false

	lastLineEndsWithParen := false

	countOpenerBeforeUglyKeyword := 0
	lastLineCountOpenerBeforeUglyKeyword := 0

	var newLineHandler = func() {
		if lastLine == nil {
			lastLine = &bytes.Buffer{}
		} else if lastLine.Len() == 0 {
			outBuf.WriteByte('\n')
		} else {
			line := bytes.TrimSpace(lastLine.Bytes())
			if len(line) != 0 {
				currentIndentLevel := indentLevel
				if indentNextLine {
					// if lastLineEndsWithParen || !lastLineParensKeywordFound {
					fmt.Printf("helloline: %s, currentline: %s\n", line, currentLine.Bytes())
					currentIndentLevel++
					// }
					indentNextLine = false
				}

				lastLineEndsWithParen = bytes.HasSuffix(line, []byte{')'})
				if lastLineCountOpenerBeforeUglyKeyword == 0 && (lastLineSawUglyKeyword || (lastLineParensKeywordFound && lastLineEndsWithParen)) {
					indentNextLine = true
				}

				for spaces := currentIndentLevel * spacesPerIndent; spaces > 0; spaces-- {
					outBuf.WriteByte(' ')
				}

				_, _ = outBuf.Write(line)
			}

			outBuf.WriteByte('\n')

			if lastOpenCounter > 0 {
				indentLevel++
			} else if currentOpenCounter < 0 {
				indentLevel--
			}

			fmt.Printf("lastLineOpens: %d, currentLineOpens: %d, indentLevel: %d, indentNextLine: %t, lastLineEndsWithParen: %t", lastOpenCounter, currentOpenCounter, indentLevel, indentNextLine, lastLineEndsWithParen)
			fmt.Printf("\nsawUglyKeyword: %t, lastLineSawUglyKeyword: %t, openBeforeKeyword: %d, lastLineOpenBeforeKeyword: %d\nlastLine: %s\n", sawUglyKeyword, lastLineSawUglyKeyword, countOpenerBeforeUglyKeyword, lastLineCountOpenerBeforeUglyKeyword, line)
		}

		lastOpenCounter = currentOpenCounter
		currentOpenCounter = 0

		lastLineSawUglyKeyword = sawUglyKeyword
		lastLineCountOpenerBeforeUglyKeyword = countOpenerBeforeUglyKeyword
		sawUglyKeyword = false
		countOpenerBeforeUglyKeyword = 0
		lastLineParensKeywordFound = parensKeywordFound
		parensKeywordFound = false

		fmt.Printf("cl: %s\n", bytes.TrimSpace(currentLine.Bytes()))
		lastLine.Reset()
		_, _ = lastLine.Write(currentLine.Bytes())
	}

	for _, c := range input {
		// var charStr = string(c)
		// if charStr == "\n" {
		// 	charStr = "\\n"
		// }
		// if len(stack.s) > 0 {
		// 	fmt.Printf("%-17v %v\n%-2s -> ", state, stack.s, charStr)
		// } else {
		// 	fmt.Printf("%-17v\n%-2s -> ", state, charStr)
		// }

		if (state == stateStringSingle && c != '\'') || (state == stateStringDouble && c != '"') || state == stateStringEscape {
			currentLine.WriteByte(c)
			continue
		}

		if state == stateComment {
			if c == '\n' {
				state = stack.Pop()
			} else {
				currentLine.WriteByte(c)
				continue
			}
		}

		if state == stateCommentBlock && (c != '*' && c != '/') {
			currentLine.WriteByte(c)
			continue
		}

		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '$' || c == '.' || (c >= '0' && c <= '9') {
			if state != stateIdentifier {
				stack.Push(state)
			}
			state = stateIdentifier
			identBuf.WriteByte(c)
		} else if state == stateIdentifier {
			var identBufStr = identBuf.String()
			if identBufStr == "do" {
				foundDo = true
			}
			fmt.Printf("IdentBuf: %s\n", identBufStr)

			if _, ok := parensKeywordsMap[identBufStr]; ok {
				if identBufStr == "while" && foundDo {
					foundDo = false
				} else if currentOpenCounter == 0 {
					parensKeywordFound = true
					countOpenerBeforeUglyKeyword = currentOpenCounter
				}
			}

			if _, ok := keywordsMap[identBufStr]; ok {
				sawUglyKeyword = true
				parensKeywordFound = false
				countOpenerBeforeUglyKeyword = currentOpenCounter
			} else {
				sawUglyKeyword = false
			}

			fmt.Printf("sawUglyKeyword: %t, countOpenerBeforeUglyKeyword: %d, parensKeywordFound: %t, foundDo: %t\n", sawUglyKeyword, countOpenerBeforeUglyKeyword, parensKeywordFound, foundDo)
			state = stack.Pop()
			identBuf.Reset()
		}

		switch c {
		case '/':
			if state == stateCommentStart {
				state = stateComment
			} else if state == stateCommentEnd {
				state = stack.Pop()
			} else {
				stack.Push(state)
				state = stateCommentStart
			}
		case '*':
			if state == stateCommentStart {
				state = stateCommentBlock
			} else if state == stateCommentBlock {
				state = stateCommentEnd
			}
		case ' ':
		case '\n':
			if state == stateComment {
				state = stack.Pop()
			}
			newLineHandler()
			currentLine.Reset()
			fmt.Println()
			continue
		case '\\':
			if state != stateStringDouble && state != stateStringSingle {
				break
			}
			stack.Push(state)
			state = stateStringEscape
		case '\'':
			if state == stateStringSingle {
				state = stack.Pop()
				break
			}
			stack.Push(state)
			state = stateStringSingle
		case '"':
			if state == stateStringDouble {
				state = stack.Pop()
				break
			}
			stack.Push(state)
			state = stateStringDouble
		case '{':
			stack.Push(state)
			state = stateBrace
			currentOpenCounter++
			parensKeywordFound = false
		case '}':
			if state != stateBrace {
				panic("expected stateBrace")
			}
			state = stack.Pop()
			currentOpenCounter--
			parensKeywordFound = false
		case '(':
			stack.Push(state)
			state = stateParen
			currentOpenCounter++
		case ')':
			if state != stateParen {
				panic("expected stateParen")
			}
			state = stack.Pop()
			currentOpenCounter--
		case '[':
			stack.Push(state)
			state = stateArray
			currentOpenCounter++
		case ']':
			if state != stateArray {
				panic("expected stateArray")
			}
			state = stack.Pop()
			currentOpenCounter--
		}

		currentLine.WriteByte(c)
	}

	newLineHandler()
	_, _ = io.Copy(outBuf, currentLine)
	return outBuf.String()
}

func main() {
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Unable to read file:", err)
	}

	fmt.Print(doIndent(string(file)))
}

var (
	stringCharsMap    = map[string]struct{}{}
	openCharsMap      = map[string]struct{}{}
	closeCharsMap     = map[string]struct{}{}
	operatorsMap      = map[string]struct{}{}
	parensKeywordsMap = map[string]struct{}{}
	keywordsMap       = map[string]struct{}{}
)

var stringChars = []byte{'\'', '"'}
var openChars = []byte{'[', '(', ',', '{'}
var closeChars = []byte{']', ')', '}'}
var operators = []string{
	`+`,
	`:`,
	`=`,
	`-`,
	`*`,
	`/`,
	`%`,
	`&`,
	`|`,
	`!`,
	`++`,
	`--`,
	`==`,
	`!=`,
	`>`,
	`>=`,
	`<`,
	`<=`,
	`&&`,
	`||`,
	`^`,
	`~`,
	`<<`,
	`>>`,
	`>>>`,
	`+=`,
	`-=`,
	`*=`,
	`/=`,
	`%=`,
	`&=`,
	`^=`,
	`!=`,
	`<<=`,
	`>>=`,
	`>>>=`,
	`?:`,
}

var parensKeywords = []string{
	"if",
	"while",
	"for",
}

var keywords = []string{
	`do`,
	`instanceof`,
	`typeof`,
	`case`,
	`else`,
	`new`,
	`var`,
	`this`,
	`with`,
	`default`,
	`delete`,
	`in`,
	`try`,
}
