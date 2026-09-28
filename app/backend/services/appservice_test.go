package services_test

import (
	"ariadne/backend/services"
	"testing"
)

func TestB_Display_App_Name(t *testing.T) {
	service := &services.AppService{}

	got := service.GetAppName()
	want := "Project ARIADNE Prototype"
	// want := "WRONG"

	if got != want {
		t.Errorf("GetAppName() = %q, want %q", got, want)
	}
}

// FIXME:
func TestUnclassified(t *testing.T) {
}
