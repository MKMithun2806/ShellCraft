package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

type Encoder struct {
	Name string
	Wrap func(payload string) string
}

func GetEncoders() []Encoder {
	return []Encoder{
		{
			Name: "Raw",
			Wrap: func(p string) string { return p },
		},
		{
			Name: "URL Encoded",
			Wrap: func(p string) string { return url.QueryEscape(p) },
		},
		{
			Name: "Base64/Pipeline",
			Wrap: func(p string) string {
				b64 := base64.StdEncoding.EncodeToString([]byte(p))
				return fmt.Sprintf("echo %s | base64 -d | sh", b64)
			},
		},
	}
}
