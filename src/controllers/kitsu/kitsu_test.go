package kitsu

import (
	basicauth "bot/src/controllers/basicauth"
	config "bot/src/controllers/config"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	os.Setenv("TEST", "true")

	// Conf
	conf := config.Read()

	// Basic auth
	JWTToken := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("JWTToken", JWTToken)
}

func TestGetPerson(t *testing.T) {
	id := "2bc7bfa1-3a66-41ed-ba4d-a14d4db61126"

	resp := GetPerson(id)

	if resp.ID != "2bc7bfa1-3a66-41ed-ba4d-a14d4db61126" {
		t.Error("ID")
	}

	/*
		if resp.CreatedAt != "2021-03-22T07:39:34" {
			t.Error("CreatedAt")
		}
	*/

	if resp.FirstName != "Innokentiy" {
		t.Error("FirstName")
	}

	if resp.LastName != "Sokolov" {
		t.Error("LastName")
	}

	/*if resp.Email != "sv.krekov@gmail.com" {
		t.Error("Email")
	}
	*/

	if resp.Phone != "@keshon" {
		t.Error("@keshon")
	}
}
