package main

import (
	"bytes"
	"fmt"
	"io"
)

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

	foundCloserBeforeOpener := false
	lastFoundCloserBeforeOpener := false
	foundOpener := false

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
					currentIndentLevel++
					indentNextLine = false
				}

				lastLineEndsWithParen = bytes.HasSuffix(line, []byte{')'})
				if lastLineCountOpenerBeforeUglyKeyword == 0 && (lastLineSawUglyKeyword || (lastLineParensKeywordFound && lastLineEndsWithParen)) {
					indentNextLine = true
				}

				if lastFoundCloserBeforeOpener && lastOpenCounter >= 0 {
					currentIndentLevel--
				}

				for spaces := currentIndentLevel * spacesPerIndent; spaces > 0; spaces-- {
					outBuf.WriteByte(' ')
				}

				_, _ = outBuf.Write(line)
			}

			outBuf.WriteByte('\n')

			if currentOpenCounter < 0 {
				indentLevel--
			} else if lastOpenCounter > 0 {
				indentLevel++
			}

			if *flagDebug {
				fmt.Printf("lastLineOpens: %d, currentLineOpens: %d, indentLevel: %d, indentNextLine: %t, lastLineEndsWithParen: %t", lastOpenCounter, currentOpenCounter, indentLevel, indentNextLine, lastLineEndsWithParen)
				fmt.Printf("\nsawUglyKeyword: %t, lastLineSawUglyKeyword: %t, openBeforeKeyword: %d, lastLineOpenBeforeKeyword: %d\nlastLine: %s\n", sawUglyKeyword, lastLineSawUglyKeyword, countOpenerBeforeUglyKeyword, lastLineCountOpenerBeforeUglyKeyword, line)
			}
		}

		lastOpenCounter = currentOpenCounter
		currentOpenCounter = 0

		lastLineSawUglyKeyword = sawUglyKeyword
		lastLineCountOpenerBeforeUglyKeyword = countOpenerBeforeUglyKeyword
		sawUglyKeyword = false
		countOpenerBeforeUglyKeyword = 0
		lastLineParensKeywordFound = parensKeywordFound
		parensKeywordFound = false
		lastFoundCloserBeforeOpener = foundCloserBeforeOpener
		foundCloserBeforeOpener = false
		foundOpener = false

		if *flagDebug {
			fmt.Printf("cl: %s\n", bytes.TrimSpace(currentLine.Bytes()))
		}
		lastLine.Reset()
		_, _ = lastLine.Write(currentLine.Bytes())
	}

	for _, c := range input {
		if *flagDebug {
			var charStr = string(c)
			if charStr == "\n" {
				charStr = "\\n"
			}
			if len(stack.s) > 0 {
				fmt.Printf("%-17v %v\n%-2s -> ", state, stack.s, charStr)
			} else {
				fmt.Printf("%-17v\n%-2s -> ", state, charStr)
			}
		}

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

			if *flagDebug {
				fmt.Printf("sawUglyKeyword: %t, countOpenerBeforeUglyKeyword: %d, parensKeywordFound: %t, foundDo: %t\n", sawUglyKeyword, countOpenerBeforeUglyKeyword, parensKeywordFound, foundDo)
			}
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
			if *flagDebug {
				fmt.Println()
			}
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
			foundOpener = true
			currentOpenCounter++
			parensKeywordFound = false
		case '}':
			if state != stateBrace {
				panic("expected stateBrace")
			}
			state = stack.Pop()
			currentOpenCounter--
			if !foundOpener {
				foundCloserBeforeOpener = true
			}
			parensKeywordFound = false
		case '(':
			stack.Push(state)
			state = stateParen
			foundOpener = true
			currentOpenCounter++
		case ')':
			if state != stateParen {
				panic("expected stateParen")
			}
			state = stack.Pop()
			if !foundOpener {
				foundCloserBeforeOpener = true
			}
			currentOpenCounter--
		case '[':
			stack.Push(state)
			state = stateArray
			foundOpener = true
			currentOpenCounter++
		case ']':
			if state != stateArray {
				panic("expected stateArray")
			}
			state = stack.Pop()
			if !foundOpener {
				foundCloserBeforeOpener = true
			}
			currentOpenCounter--
		}

		currentLine.WriteByte(c)
	}

	newLineHandler()
	_, _ = io.Copy(outBuf, currentLine)
	return outBuf.String()
}
