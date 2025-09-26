package router

import (
	"bytes"

	router_context "github.com/aomirun/content-router/context"
)

// PrefixMatcher 创建一个前缀匹配器
func PrefixMatcher(prefix string) Matcher {
	prefixBytes := []byte(prefix)
	return MatcherFunc(func(ctx router_context.Context) bool {
		data := ctx.Buffer().Get()
		return len(data) >= len(prefixBytes) && bytes.HasPrefix(data, prefixBytes)
	})
}

// SuffixMatcher 创建一个后缀匹配器
func SuffixMatcher(suffix string) Matcher {
	suffixBytes := []byte(suffix)
	return MatcherFunc(func(ctx router_context.Context) bool {
		data := ctx.Buffer().Get()
		return len(data) >= len(suffixBytes) && bytes.HasSuffix(data, suffixBytes)
	})
}

// ContainsMatcher 创建一个包含匹配器
func ContainsMatcher(substring string) Matcher {
	substringBytes := []byte(substring)
	return MatcherFunc(func(ctx router_context.Context) bool {
		data := ctx.Buffer().Get()
		return len(data) >= len(substringBytes) && bytes.Contains(data, substringBytes)
	})
}
