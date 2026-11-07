package services_test

import (
	"ariadne/backend/tests/createtask"
	"ariadne/backend/tests/gettask"
	"ariadne/backend/tests/listtasks"
	"ariadne/backend/tests/modifytask"
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

// TEST:
func Test_GetTask_Success(t *testing.T) {
	gettask.Success(t)
}

// TEST:
func Test_GetTask_NotFound(t *testing.T) {
	gettask.NotFound(t)
}

// TEST:
func Test_GetTask_ReadFailure(t *testing.T) {
	gettask.ReadFailure(t)
}

// TEST:
func Test_CreateTask_Success(t *testing.T) {
	createtask.Success(t)
}

// TEST:
func Test_CreateTask_ReadFailure(t *testing.T) {
	createtask.ReadFailure(t)
}

// TEST:
func Test_CreateTask_UnexpectedError(t *testing.T) {
	createtask.UnexpectedError(t)
}

// TEST:
func Test_ModifyTask_Success(t *testing.T) {
	modifytask.Success(t)
}

// TEST:
func Test_ModifyTask_NotFound(t *testing.T) {
	modifytask.NotFound(t)
}

// TEST:
func Test_ModifyTask_UnexpectedError(t *testing.T) {
	modifytask.UnexpectedError(t)
}
