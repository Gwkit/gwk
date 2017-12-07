package gwk

import (
	"honnef.co/go/js/dom"
	"math"
	"math/rand"
	"regexp"
	"strings"
)

type Layout struct {
}

func (*Layout) canBreakBefore(char rune) bool {
	chars := []rune{' ', '\t', '.', ']', ')', '}', ',', '?', ';', ':', '!', '"', '\'', '。', '？', '、', '”', '’', '】', '》', '）', '：', '；', '，'}
	for i := range chars {
		if chars[i] == char {
			return false
		}
	}

	return true
}

func (*Layout) justifyText(context *dom.CanvasRenderingContext2D, line string, maxWidth int) string {
	M := 0
	words := strings.Split(line, " ")
	N := len(words) - 1
	spaceWidth := context.MeasureText(" ").Width
	lastWidth := float64(maxWidth) - context.MeasureText(line).Width

	if lastWidth < spaceWidth || N == 0 {
		return line
	} else {
		M = int(math.Floor(lastWidth / spaceWidth))
	}

	residue := M % N
	num := int(math.Floor(float64(M / N)))
	positions := make([]int, N)

	for i := range positions {
		positions[i] = -1
	}

	realNum := 0
	for i := 0; i < residue; i++ {
		ran := int(math.Floor(rand.Float64() * float64(N)))
		if positions[ran] == -1 {
			positions[ran] = 0
		} else {
			positions[ran] += 1
		}
	}

	result := ""
	for i := 0; i < N; i++ {
		result += words[i] + " "
		realNum = 0

		if positions[i] == -1 {
			realNum = positions[i] + 1
		}

		realNum += num

		for j := 0; j < realNum; j++ {
			result += " "
		}
	}

	if N > 0 {
		result += words[N]
	}

	return result
}

func (layout *Layout) wrapBySpace(context *dom.CanvasRenderingContext2D, width int, content string, justity bool) []string {
	result := []string{}

	contents := strings.Split(content, "\n")
	for i := 0; i < len(contents); i++ {
		line := ""
		item := contents[i]
		words := strings.Split(item, " ")
		for j := 0; j < len(words); j++ {
			lineTest := line + words[j] + " "
			if int(context.MeasureText(lineTest).Width) > width {
				line = line[0 : len(line)-1]
				if justity {
					result = append(result, layout.justifyText(context, line, width))
				} else {
					result = append(result, line)
				}
			} else {
				line = lineTest
			}
		}

		if len(line) > 0 {
			if justity {
				result = append(result, layout.justifyText(context, strings.TrimSpace(line), width))
			} else {
				result = append(result, strings.TrimSpace(line))
			}
		}
	}

	return result
}

func (*Layout) checkCJK(char rune) bool {
	regex := regexp.MustCompile("[\u4E00-\u9FA5]")
	return regex.MatchString(string([]rune{char}))
}

func (layout *Layout) getWordEnd(text string, start int) int {
	if len(text) <= start {
		return -1
	}

	rtext := []rune(text)
	c := rtext[start]
	if layout.checkCJK(c) || c == ' ' {
		start++
		if start < len(text) && !layout.canBreakBefore(rtext[start]) {
			start++
		}
		return start
	} else {
		flag := false
		var i = start + 1
		for ; i < len(text); i++ {
			if rtext[i] == ' ' {
				flag = true
				break
			}
		}
		if flag {
			return i
		} else {
			return i - 1
		}
	}
}

func (layout *Layout) wrapByWord(context *dom.CanvasRenderingContext2D, width int, content string) []string {
	contents := strings.Split(content, "\n")
	result := make([]string, 0, len(contents))
	for _, it := range contents {
		var line, word, lineTest string
		if it == "" {
			result = append(result, it)
			continue
		}

		var startIndex int
		itr := []rune(it)
		for startIndex = 0; startIndex < len(itr); startIndex++ {
			endIndex := layout.getWordEnd(it, startIndex)
			if endIndex != -1 {
				word = it[startIndex:endIndex]
				lineTest = line + word
				startIndex = endIndex
			} else {
				break
			}

			if int(context.MeasureText(lineTest).Width) > width {
				if int(context.MeasureText(word).Width) > width {
					if line != "" {
						result = append(result, line)
					}
					var singleLine, singleLineTest string
					for c := range word {
						singleLineTest = singleLine + string(c)
						if int(context.MeasureText(singleLineTest).Width) > width {
							result = append(result, singleLine)
							singleLine = string(c)
						} else {
							singleLine = singleLineTest
						}
					}
					word = singleLine
				} else {
					result = append(result, line)
				}
				line = word
			} else {
				line = lineTest
			}
		}
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return result
}

var layoutInstance = &Layout{}

func layoutText(context *dom.CanvasRenderingContext2D, fontSize int, str string, width int, flexibleWidth int) []string {
	return layoutInstance.wrapByWord(context, width, str)
}
