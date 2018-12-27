package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	TestCoin
	Name string `json:"name"`
}

type TestCoin struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func TestSelectFields(t *testing.T) {
	coin := TestCoin{
		123,
		"BTC",
	}

	user := TestUser{
		TestCoin: coin,
		Name:     "yiplee",
	}

	fileds := SelectFields(coin, "id")

	assert.Len(t, fileds, 1)

	// if id := fileds["id"]; id.(uint) != 123 {
	// 	t.Error("id should be selected")
	// }

	// if _, found := fileds["name"]; found {
	// 	t.Error("name should be unselected")
	// }

	fileds = UnselectFields(user)
	assert.NotNil(t, fileds["id"])
}

func TestExportTime(t *testing.T) {
	type A struct {
		Create time.Time `json:"create,update,unix"`
	}

	a := A{Create: time.Now()}
	out := UnselectFields(a)
	assert.Empty(t, out)
}

func TestUnselectedFields(t *testing.T) {
	coin := TestCoin{
		123,
		"BTC",
	}

	fileds := UnselectFields(coin, "id")

	if len := len(fileds); len != 1 {
		t.Error("number of selected fileds should be 1")
	}

	if _, found := fileds["id"]; found {
		t.Error("id should be unselected")
	}

	if name := fileds["name"].(string); name != "BTC" {
		t.Errorf("name : %s should be selected", name)
	}
}

func TestJsonDict(t *testing.T) {
	var objects []interface{}
	objects = append(objects, "key")

	var json = Map(objects...)
	if len := len(json); len != 0 {
		t.Error("json should be empty")
	}

	objects = append(objects, TestCoin{})
	json = Map(objects...)
	if len := len(json); len != 1 {
		t.Error("json should not have one entries")
	}
}
