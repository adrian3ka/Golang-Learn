package main

import (
	"testing"
	"fmt"
)

func TestNewUserNoUsername(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	_, err := NewUser("", "user@example.com", "password")
	fmt.Println("err >>",err)
	if err != errNoUsername {
		t.Error("Expected err to be errNoUsername")
	}
}

func TestNewUserNoPassword(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	_, err := NewUser("user", "user@example.com", "")
	fmt.Println("err >>",err)
	if err != errNoPassword {
		t.Error("Expected err to be errNoUsername")
	}
}

type MockUserStore struct {
	findUser         *User
	Kalau_FindByEmail_Dipanggil_seolah_olah_balikin_user_tipedatanya_aja_PointerUser    *User
	findUsernameUser *User
	saveUser         *User
}

func (store *MockUserStore) Find(string) (*User, error) {
	//fmt.Println("Mock Find Called")
	//fmt.Println("store.findUser >>",store.findUser)
	return store.findUser, nil
}

func (store *MockUserStore) FindByEmail(string) (*User, error) {
	//fmt.Println("Mock FindEmail Called")
	return store.Kalau_FindByEmail_Dipanggil_seolah_olah_balikin_user_tipedatanya_aja_PointerUser, nil
}

func (store *MockUserStore) FindByUsername(string) (*User, error) {
	//fmt.Println("Mock FindUsername Called")
	//fmt.Println("store.findUsernameUser >>",store.findUsernameUser)
	return store.findUsernameUser, nil
}

func (store *MockUserStore) Save(user User) error {
	store.saveUser = &user
	return nil
}

func TestNewUserExistingUsername(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	globalUserStore = &MockUserStore{
		findUsernameUser: &User{},
	}
	//fmt.Println("globalUserStore >>",globalUserStore)
	_, err := NewUser("userzxcasdqwf", "user@example.com", "somepassword")
	fmt.Println("err >>",err)
	if err != errUsernameExists {
		t.Error("Expected err to be errUsernameExists")
	}
}

func TestNewUserExistingEmail(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	globalUserStore = &MockUserStore{
		Kalau_FindByEmail_Dipanggil_seolah_olah_balikin_user_tipedatanya_aja_PointerUser: &User{},
	}
	//fmt.Println("globalUserStore >>",globalUserStore)
	//Kan di NewUser dia panggil fungsi findByEmail, nah karena di awal udah di set &MockUserStorenya, maka bakalan
	//balikin nilai palsu dan seolah-olah ada yang emailnya kembar
	_, err := NewUser("user", "user@example.com", "somepassword")
	fmt.Println("err >>",err)
	if err != errEmailExists {
		t.Error("Expected err to be errEmailExists")
	}
}
