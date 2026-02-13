// Package stubs provides embedded migration stub templates.
package stubs

import (
	_ "embed"
)

// Blank is the blank migration stub template.
//
//go:embed migration.stub
var Blank string

// Create is the create table migration stub template.
//
//go:embed migration.create.stub
var Create string

// Update is the update table migration stub template.
//
//go:embed migration.update.stub
var Update string
