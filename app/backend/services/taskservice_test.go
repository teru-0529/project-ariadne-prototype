package services_test

import (
	"ariadne/backend/tests/listtasks"
	"testing"
)

// TEST:
func Test_ListTasks_Success(t *testing.T) {
	listtasks.Success(t)
}

// TEST:
func Test_ListTasks_empty(t *testing.T) {
	listtasks.Empty(t)
}

// TEST:
func Test_ListTasks_ReadFailure_file(t *testing.T) {
	listtasks.ReadFailureFile(t)
}

// TEST:
func Test_ListTasks_ReadFailure_format_Header(t *testing.T) {
	listtasks.ReadFailureFormatHeader(t)
}

// TEST:
func Test_ListTasks_Unsupported_major(t *testing.T) {
	listtasks.UnsupportedMajor(t)
}

// TEST:
func Test_ListTasks_ReadFailure_Unsupported_Minor(t *testing.T) {
	listtasks.ReadFailureUnsupportedMinor(t)
}

// TEST:
func Test_ListTasks_ReadFailure_Parser(t *testing.T) {
	listtasks.ReadFailureParser(t)
}

// TEST:
func Test_ListTasks_ReadFailure_format_Dataset(t *testing.T) {
	listtasks.ReadFailureFormatDataset(t)
}

// TEST:
func Test_ListTasks_ReadFailure_invalid_model(t *testing.T) {
	listtasks.ReadFailureInvalidModel(t)
}
