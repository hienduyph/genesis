package coders

import "github.com/gorilla/schema"

var Query = func() *schema.Decoder {
	c := schema.NewDecoder()
	c.SetAliasTag("json")
	return c
}()
