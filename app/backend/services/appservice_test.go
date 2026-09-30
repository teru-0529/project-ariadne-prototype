package services_test

import (
	"ariadne/backend/services"
	"testing"
)

func Test_Display_App_Name(t *testing.T) {
	// INFO: アプリ名称の応答
	service := &services.AppService{}

	got := service.GetAppName()
	want := "Project ARIADNE Prototype"

	if got != want {
		t.Errorf("GetAppName() = %q, want %q", got, want)
	}
}

func Test_Small_Logic(t *testing.T) {
	// INFO: ダミーテスト（D：Logicテスト相当）
	if !*logic {
		t.SkipNow()
	}
}
