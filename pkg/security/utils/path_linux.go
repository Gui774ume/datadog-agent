// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package utils

// PathPatternMatchOpts PathPatternMatch options
type PathPatternMatchOpts struct {
	WildcardLimit      int // max number of wildcard in the pattern
	PrefixNodeRequired int // number of prefix nodes required
	SuffixNodeRequired int // number of suffix nodes required
	NodeSizeLimit      int // min size required to substitute with a wildcard
}

// PathPatternMatch pattern builder for files
func PathPatternMatch(pattern string, path string, opts PathPatternMatchOpts) bool {
	return PatternMatch(pattern, path, opts, '/')
}

// PatternMatch pattern builder for files
func PatternMatch(pattern string, path string, opts PathPatternMatchOpts, separator uint8) bool {
	var (
		i, j, offsetPath                     = 0, 0, 0
		wildcardCount, nodeCount, suffixNode = 0, 0, 0
		patternLen, pathLen                  = len(pattern), len(path)
		wildcard                             bool

		computeNode = func() bool {
			if wildcard {
				wildcardCount++
				if wildcardCount > opts.WildcardLimit {
					return false
				}
				if nodeCount < opts.PrefixNodeRequired {
					return false
				}
				if opts.NodeSizeLimit != 0 && j-offsetPath < opts.NodeSizeLimit {
					return false
				}

				suffixNode = 0
			} else {
				suffixNode++
			}

			offsetPath = j

			if i > 0 {
				nodeCount++
			}
			return true
		}
	)

	if patternLen > 0 && pattern[0] != separator {
		return false
	}

	if pathLen > 0 && path[0] != separator {
		return false
	}

	for i < len(pattern) && j < len(path) {
		pn, ph := pattern[i], path[j]
		if pn == separator && ph == separator {
			if !computeNode() {
				return false
			}
			wildcard = false

			i++
			j++
			continue
		}

		if pn != ph {
			wildcard = true
		}
		if pn != separator {
			i++
		}
		if ph != separator {
			j++
		}
	}

	if patternLen != i || pathLen != j {
		wildcard = true
	}

	for i < patternLen {
		if pattern[i] == separator {
			return false
		}
		i++
	}

	for j < pathLen {
		if path[j] == separator {
			return false
		}
		j++
	}

	if !computeNode() {
		return false
	}

	if opts.SuffixNodeRequired == 0 || suffixNode >= opts.SuffixNodeRequired {
		return true
	}

	return false
}
