// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package environment

// Env type
type Env string

// String implements fmt.Stringer
func (e Env) String() string {
	return string(e)
}

// Env enums
const (
	Prod    Env = "prod"
	PreProd Env = "preprod"
	Dev     Env = "dev"
)

// FromString returns Env
func FromString(env string) Env {
	switch env {
	case "prod":
		return Prod
	case "preprod":
		return PreProd
	default:
		return Dev
	}
}
