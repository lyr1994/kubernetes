/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package content

import (
	"regexp"
	"unicode"
)

// This regex describes the interior of a label, which is slightly different
// than the rules for the first and last characters. For better errors, we
// handle them seperately.
const dns1123LabelInteriorFmt string = "[-a-z0-9]+"
const dns1123LabelMaxLength int = 63

// DNS1123LabelMaxLength is a DNS label's max length (RFC 1123).
const DNS1123LabelMaxLength int = dns1123LabelMaxLength

var dnsLabelRegexp = regexp.MustCompile("^" + dns1123LabelInteriorFmt + "$")

// IsDNS1123Label returns error messages if the specified value does not
// parse as per the definition of a label in DNS (approximately RFC 1123).
func IsDNS1123Label(value string) []string {
	var errs []string
	if len(value) > dns1123LabelMaxLength {
		errs = append(errs, MaxLenError(dns1123LabelMaxLength))
		return errs // Don't run further validation if we know it is too long.
	}
	if len(value) == 0 {
		errs = append(errs, "must contain at least 1 character")
		return errs // No point in going further.
	}

	isAlNum := func(r rune) bool {
		if r > unicode.MaxASCII {
			return false
		}
		if unicode.IsLetter(r) && unicode.IsLower(r) {
			return true
		}
		if unicode.IsDigit(r) {
			return true
		}
		return false
	}
	runes := []rune(value)
	if !isAlNum(runes[0]) || !isAlNum(runes[len(runes)-1]) {
		errs = append(errs, "must start and end with lower-case alphanumeric characters")
	}
	if len(runes) > 2 && !dnsLabelRegexp.MatchString(string(runes[1:len(runes)-1])) {
		errs = append(errs, "must contain only lower-case alphanumeric characters or '-'")
	}
	return errs
}
