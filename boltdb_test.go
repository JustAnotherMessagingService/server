package main

import (
	"testing"

	"github.com/twinj/uuid"
)

var ID = uuid.NewV1().String()
var ID2 = uuid.NewV1().String()

var USER = "test1"
var PASS = "hunter2"

var user = &User{Id: ID, Username: USER, Password: PASS}
var userDupeName = &User{Id: ID2, Username: USER, Password: PASS}

func TestBoltDBOpen(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Try and open another.
	_, err = BoltDBOpen(DBFILE)
	if err == nil {
		t.Errorf("boltdb: error expected with opening dir")
	}
}

func TestBoltSaveUser(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	var nilUser *User
	nilUser = nil
	// Test saving a nil user.
	err = db.SaveUser(nilUser)
	if err != ErrUserObjectNil {
		t.Errorf("Saving user with nil object did not return ErrUserObjectNil")
	}

	// Test regular ol' save.
	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test a user with a duplicate username.
	err = db.SaveUser(userDupeName)
	if err != ErrUsernameAlreadyExists {
		t.Errorf("Saving user with duplicate username did not return ErrUsernameAlreadyExists")
	}
}

func TestBoltGetUserById(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Test getting user with blank Id
	u, err := db.GetUserById("")
	if err != ErrIdCannotBeEmpty {
		t.Errorf("Getting user with blank id did not return ErrIdCannotBeEmpty")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting user with unknown Id
	u, err = db.GetUserById("1234-1234-1234-1234")
	if err != ErrUserNotFound {
		t.Errorf("Getting user with no known Id did not return ErrUserNotFound")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting user by a real Id.
	u, err = db.GetUserById(user.Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if u == nil {
		t.Errorf("GetUserById returned nil user")
	} else {
		testUsersEqual(u, user, t)
	}
}

func TestBoltGetUserByUsername(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// Test getting user with blank username
	u, err := db.GetUserByUsername("")
	if err != ErrUsernameCannotBeEmpty {
		t.Errorf("Getting user with blank username did not return ErrIdCannotBeEmpty")
	}
	if u != nil {
		t.Errorf("Getting user with blank username did not return nil user")
	}

	// Test getting username by unknown username
	u, err = db.GetUserByUsername("1234-1234-1234-1234")
	if err != ErrUserNotFound {
		t.Errorf("Getting user with no known Username did not return ErrUserNotFound")
	}
	if u != nil {
		t.Errorf("Getting user with unknown username did not return nil user")
	}

	// Test getting user by a real username
	u, err = db.GetUserByUsername(user.Username)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if u == nil {
			t.Errorf("GetUserById returned nil user")
		} else {
			testUsersEqual(u, user, t)
		}
	}
}

func TestBoltGetAllUsers(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	users, err := db.GetAllUsers()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(users) == 0 {
		t.Errorf("GetAllUsers had an unexpected length of 0")
	}
}

func TestBoltDeleteUser(t *testing.T) {
	// Setup DB
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	var nilUser *User
	nilUser = nil
	// Test saving a nil user.
	err = db.DeleteUser(nilUser)
	if err != ErrUserObjectNil {
		t.Errorf("Saving user with nil object did not return ErrUserObjectNil")
	}

	// Delete this test user form the DB.
	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBoltMessageLifecycle(t *testing.T) {
	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	// TODO: Create some sort of init function to handle the calling of any
	// necessary init requirements throughout server.
	mes := NewMessage("test message body")
	err = db.SaveMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestBoltMessageSend(t *testing.T) {
	user := &User{
		Id:       ID,
		Username: USER,
		Password: PASS,
	}

	db, err := BoltDBOpen(DBFILE)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer db.Conn.Close()

	err = db.SaveUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	// TODO: Create some sort of init function to handle the calling of any
	// necessary init requirements throughout server.
	mes := NewMessage("test message body")
	err = db.SaveMessage(mes)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testUsersEqual(u1, u2 *User, t *testing.T) {
	if u2 == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if u1.Id != u2.Id {
		t.Errorf("Id of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
	if u1.Username != u2.Username {
		t.Errorf("Username of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
	if u1.Password != u2.Password {
		t.Errorf("Password of retrieved user does not match stored user: %s vs %s", u1.Id, u2.Id)
	}
}

func testGetUserByUsername(username string, t *testing.T) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("User should not be nil.")
		return
	}
	if user.Id != ID {
		t.Errorf("Id of retrieved user does not match stored user: %s vs %s", user.Id, ID)
	}
	if user.Username != USER {
		t.Errorf("Username of retrieved user does not match stored user")
	}
	if user.Password != PASS {
		t.Errorf("Password of retrieved user does not match stored user")
	}
}
