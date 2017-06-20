package gohtml

import "testing"

func TestParseArticleLink(t *testing.T) {
	a, err := ParseArticleLink("http://www.jianshu.com/p/8f37c12a854f")
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	t.Log(a.Title)
	t.Log(a.Content)
}
