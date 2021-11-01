package coders

import "github.com/gorilla/schema"

var queryDecoder = func() *schema.Decoder {
	c := schema.NewDecoder()
	c.SetAliasTag("json")
	return c
}()
var queryEncode = func() *schema.Encoder {
	c := schema.NewEncoder()
	c.SetAliasTag("json")
	return c
}()

var DecodeQuery = queryDecoder.Decode
var EncodeQuery = queryEncode.Encode
