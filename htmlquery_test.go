package htmlquery

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

var body string = `
<html>
<head></head>
<body>
  <div class="level_1">first</div>
  <div class="level_1">second
	<div class="level_2">inside second</div>
  </div>
</body
</html>
`

func TestFirstNode(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	node := FindFirstNode(root, HasTag("div"))
	assert.NotNil(t, node)
	assert.Equal(t, "first", FirstChildNodeText(node))
}

func TestFindAllRec(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	texts := []string{}
	for _, node := range FindAllNodesRec(root, HasTag("div")) {
		t := trim(FirstChildNodeText(node))
		texts = append(texts, t)
	}
	assert.Equal(t, []string{"first", "second", "inside second"}, texts)
}

func TestListChildren(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	bodyNode := FindFirstNode(root, HasTag("html"))
	texts := []string{}
	for _, node := range ListChildren(bodyNode, All()) {
		if node.Type == html.ElementNode {
			t := node.Data
			texts = append(texts, t)
		}
	}
	assert.Equal(t, []string{"head", "body"}, texts)
}

func TestSelectAllMatchers(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	texts := []string{}
	for _, node := range FindAllNodesRec(root, All(HasTag("div"), HasAttr("class", StringIs("level_1")))) {
		texts = append(texts, trim(FirstChildNodeText(node)))
	}
	assert.Equal(t, []string{"first", "second"}, texts)
}

func TestIsTextNode(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	texts := []string{}
	for _, node := range FindAllNodesRec(root, IsTextNode()) {
		t := trim(node.Data)
		if len(t) > 0 {
			texts = append(texts, t)
		}
	}
	assert.Equal(t, []string{"first", "second", "inside second"}, texts)
}

func TestStartingWith(t *testing.T) {
	root, err := html.Parse(strings.NewReader(body))
	assert.NoError(t, err)
	texts := []string{}
	for _, node := range FindAllNodesRec(root, HasAttr("class", StartingWith("level_"))) {
		t := GetAttrValue(node.Attr, "class")
		texts = append(texts, t)
	}
	assert.Equal(t, []string{"level_1", "level_1", "level_2"}, texts)
}

func trim(s string) string {
	return strings.Trim(s, "\n \t")
}
