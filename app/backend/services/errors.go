package services

import (
	"errors"
	"fmt"
)

// INFO: YAMLフォーマット読み込みに失敗した場合に返されるエラー
var ErrReadFailure = errors.New("read failure")

// INFO: YAMLフォーマットが不正な場合に返されるエラー
type FormatError struct {
	Detail string
}

// FUNCTION: FormatErrorのエラーメッセージを返す関数
func (e *FormatError) Error() string {
	return fmt.Sprintf("invalid format: %s", e.Detail)
}

// INFO: サポートされていないメジャーバージョンのYAMLフォーマットが読み込まれた場合に返されるエラー
type UnsupportedMajorVersionError struct {
	Major int
}

// FUNCTION: UnsupportedMajorVersionErrorのエラーメッセージを返す関数
func (e *UnsupportedMajorVersionError) Error() string {
	return fmt.Sprintf("unsupported major version: %d", e.Major)
}

// INFO: サポートされていないマイナーバージョンのYAMLフォーマットが読み込まれた場合に返されるエラー
type UnsupportedMinorVersionError struct {
	Major int
	Minor int
}

// FUNCTION: UnsupportedMinorVersionErrorのエラーメッセージを返す関数
func (e *UnsupportedMinorVersionError) Error() string {
	return fmt.Sprintf(
		"unsupported minor version: %d.%d",
		e.Major,
		e.Minor,
	)
}
