// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package push

import (
	"errors"
	"regexp"
	"strings"
)

//ErrInvalidTopicEmptyString is the error returned when a topic string
//is passed in that is 0 length
var ErrInvalidTopicEmptyString = errors.New("Invalid Topic; empty string")

//ErrInvalidTopicMultilevel is the error returned when a topic string
//is passed in that has the multi level wildcard in any position but
//the last
var ErrInvalidTopicMultilevel = errors.New("Invalid Topic; multi-level wildcard must be last level")

// TopicID name
type TopicID string

// IsValid returns true if topic name is valid
func (t TopicID) IsValid() error {
	if len(t) == 0 {
		return ErrInvalidTopicEmptyString
	}

	levels := strings.Split(string(t), "/")

	for i, level := range levels {
		if level == "#" && i != len(levels)-1 {
			return ErrInvalidTopicMultilevel
		}
	}

	return nil
}

// IsSystemTopic returns true if topic name is system
func (t TopicID) IsSystemTopic() bool {
	return t[0] == '$'
}

// Match pattern filter
func (t TopicID) Match(pattern string) bool {
	levels := strings.Split(pattern, "/")
	size := len(levels)

	isReg := false

	for i, level := range levels {
		if level == "+" {
			levels[i] = `[a-zA-Z0-9]+`
			isReg = true
		} else if (level == "#") && (i == (size - 1)) {
			levels[i] = `[a-zA-Z0-9\/]+`
			isReg = true
		}
	}

	var patternReg string
	if isReg {
		patternReg = strings.Join(levels, "\\/")
	} else {
		patternReg = strings.Join(levels, "/")
	}

	// escape $
	if isReg && patternReg[0:1] == "$" {
		patternReg = "\\" + patternReg
	}

	if !isReg {
		return string(t) == pattern
	}

	patternReg = "^" + patternReg + "$"

	matched, err := regexp.MatchString(patternReg, string(t))
	if err != nil {
		return false
	}

	return matched
}

// String implements fmt.Stringer
func (t TopicID) String() string {
	return string(t)
}

// Match pattern filter
/*func (t TopicID) Match(pattern string) bool {
	topicParts := strings.Split(t, "/")
	topicPartsLen := len(topicParts)

	patternParts := strings.Split(pattern, "/")
	patternPartsLen := len(patternParts)

	// topic: foo == pattern: foo/bar => false
	if topicPartsLen < patternPartsLen {
		return false
	}

	// topic: foo == pattern: + => true
	// topic: $foo == pattern: + => false
	if patternPartsLen == 1 && topicPartsLen == patternPartsLen && patternParts[0] == "+" && !t.IsSystemTopic() {
		return true
	}

	for i, part := range patternParts {
		if part == "+" {
			continue
		} else if part == topicParts[i] {

		}
	}

	return false
}
*/
