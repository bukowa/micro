/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package bolt_test

import (
	"bytes"
	. "github.com/bukowa/micro/storage/bolt"
	"os"
	"reflect"
	"testing"
)

var TestDatabase = func() (Storage, func()) {
	// create database
	os.Remove("test.dbx")
	db, err := NewStorage(nil, "test.dbx", &TestUser{})
	if err != nil {
		panic(err)
	}
	return db, func() {
		db.Bolt().Close()
		os.Remove("test.dbx")
	}
}

type TestUser struct {
	Login    []byte
	Password string
}

func (t *TestUser) SetKey(b []byte) {
	t.Login = b
}

func (t *TestUser) Key() []byte {
	return t.Login
}

func TestController_Delete(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()


	user := &TestUser{
		Login:    []byte("TestController_Delete"),
		Password: "",
	}

	// create
	EP(db.Create(user))

	exists, err := db.Exists(user)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error()
	}

	// delete
	EP(db.Delete(user))

	if err := db.Get(user); err != ErrorNotFound {
		t.Error(err)
	}

	exists, err = db.Exists(user)
	EP(err)

	// check exists
	if exists {
		t.Error()
	}
}

func TestController_Get(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()
	user := &TestUser{
		Login:    []byte("TestController_Get"),
		Password: "",
	}

	// empty
	if err := db.Get(user); err != ErrorNotFound {
		t.Error()
	}
}

func TestController_Exists(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()
	user := &TestUser{
		Login:    []byte("TestController_Exists"),
		Password: "",
	}
	//
	// create user & check
	EP(db.Create(user))
	exists, err := db.Exists(user)
	EP(err)
	if !exists {
		t.Error()
	}
	// check false
	user.Login = []byte("TestController_ExistsTestController_Exists")
	exists, err = db.Exists(user)
	EP(err)
	if exists {
		t.Error()
	}
}

func TestBoltDatabase_Init(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()


	user := &TestUser{
		Login:    []byte("TestBoltDatabase_Init"),
		Password: "TestBoltDatabase_Init",
	}
	user2 := &TestUser{
		Login: user.Login,
	}

	EP(db.Create(user))
	EP(db.Get(user2))

	if !reflect.DeepEqual(user, user2) {
		t.Error(user, user2)
	}

	if user.Password != user2.Password {
		t.Error()
	}
}

func EP(err error) {
	if err != nil {
		panic(err)
	}
}

func TestDB_GetAll(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()

	user := &TestUser{
		Login:    []byte("login"),
		Password: "passwd",
	}
	user2 := &TestUser{
		Login:    []byte("test"),
		Password: "test",
	}

	// insert user
	EP(db.Create(user))
	EP(db.Create(user2))

	// get all users
	var n int
	err := db.ForEach(user, func(m Model) {
		if n == 0 && string( user.Login) != "login" {
			t.Error()
		}
		if n == 1 && string(user.Login) != "test" {
			t.Error()
		}
		n ++
	})
	if err != nil {
		t.Error(err)
	}
	if n != 2 {
		t.Error()
	}
}

func Test_dbx_Pop(t *testing.T) {
	// setup test
	db, def := TestDatabase()
	defer def()

	user := &TestUser{
		Login:    []byte("login"),
		Password: "a",
	}
	user2 := &TestUser{
		Login:    []byte("login2"),
		Password: "b",
	}
	EP(db.Create(user))
	EP(db.Create(user2))

	userp := &TestUser{}
	err := db.Pop(userp)
	if err != nil {
		t.Error(err)
	}
	if userp.Password != user.Password {
		t.Error()
	}
	if !bytes.Equal(userp.Login, user.Login) {
		t.Error()
	}

	err = db.Pop(userp)
	if err != nil {
		t.Error(err)
	}
	if userp.Password != user2.Password {
		t.Error()
	}
	if !bytes.Equal(userp.Login, user2.Login) {
		t.Error()
	}

	err = db.Pop(userp)
	if err != ErrorNotFound {
		t.Error(err)
	}
}
