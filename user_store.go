package main
import(
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
)
type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}
type FileUserStore struct {
	filename string
	Users map[string]User
}
func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	contents, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}
	//0660 itu permission yang memberbolehkan write ke filenya
	err = ioutil.WriteFile(store.filename, contents, 0660)
	if err != nil {
		return err
	}
	return nil
}
func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}
func (store FileUserStore) FindByUsername(username string) (*User,error) {
	if username == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(username) == strings.ToLower(user.Username) {
			return &user, nil
		}
	}
	return nil, nil
}
func (store FileUserStore) FindByEmail(email string) (*User,error) {
	if email == "" {
		return nil, nil
	}
	for _, user := range store.Users {
		if strings.ToLower(email) == strings.ToLower(user.Email) {
			return &user, nil
		}
	}
	return nil, nil
}
func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		Users: map[string]User{},
		filename: filename,
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	
	if err != nil {
		return nil, err
	}
	return store, nil
}
var globalUserStore UserStore
func init() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store
}