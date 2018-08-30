package tick

import (
	"fmt"
	"testing"

	"github.com/influxdata/kapacitor/tick"
	"github.com/influxdata/kapacitor/tick/stateful"
	"github.com/stretchr/testify/assert"
	"gitlab.bearstech.com/bearstech/go-lepsius/tick/model"
)

func TestTick(t *testing.T) {
	script := `
var i = input
	|fromStdin()
		.parse(json)
	|stdout()

var fp = ['name', 'client']
var i2 = input
	|fromChan(chan)
	|grok()
		.source('message')
		.match('%{NUMBER:duration} %{IP:client}')
	|fingerprint()
		.source(fp)
		.target('uid')
	|stdout()

`
	scope := stateful.NewScope()
	input := NewInput()
	input.Test = true
	c := make(chan *model.Line, 1)
	scope.Set("input", input)
	scope.Set("json", JsonParser)
	scope.Set("chan", c)

	r, err := tick.Evaluate(script, scope, nil, false)
	assert.NoError(t, err)
	fmt.Println(r)
	i, err := scope.Get("i")
	assert.NoError(t, err)
	s, ok := i.(*Stdout)
	assert.True(t, ok)
	fmt.Println(s)
	fmt.Println(s.Input.Test)

	i2_, err := scope.Get("i2")
	assert.NoError(t, err)
	i2, ok := i2_.(*Stdout)
	assert.True(t, ok)
	assert.Len(t, i2.Input.Filters, 2)

	l, err := model.NewLine("beuha", "aussi")
	assert.NoError(t, err)
	c <- l
	fmt.Println(i2)
	l2, err := i2.Input.read()
	assert.NoError(t, err)
	fmt.Println(l2)
}
