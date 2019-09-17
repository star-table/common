package pinyin

import "testing"

func TestConvertCode(t *testing.T) {

	str := "这是One IssueOfFirst任务"
	t.Log(ConvertCode(str))

	str = "One这是OfIssue First1123456"
	t.Log(ConvertCode(str))

	str = "OneOneTwoTwoThirdThirdFourFourFive"
	t.Log(ConvertCode(str))
	t.Log(ConvertCodeWithMaxLen(str, 8))
}
