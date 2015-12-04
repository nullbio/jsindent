package main

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
