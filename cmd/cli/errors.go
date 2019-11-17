package main

import "errors"

var (
	VersionConflict = errors.New("version conflict")
	NotFound        = errors.New("not found")
	//RequiredArgument = errors.New("")
)
