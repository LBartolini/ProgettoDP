package internal

import (
	"context"
	"errors"
	"testing"

	pb "auth/proto"
)

type MockDB struct {
	credentials map[string]string
}

func NewMockDB(creds map[string]string) *MockDB {
	if creds == nil {
		creds = map[string]string{
			"test":  "12345",
			"test2": "abcde",
		}
	}

	return &MockDB{credentials: creds}
}

func (m *MockDB) Login(username string, password string) (bool, error) {
	pass, ok := m.credentials[username]

	if !ok {
		return false, errors.New("user not registered")
	}

	return password == pass, nil
}

func (m *MockDB) Register(username string, password string, email string, phone string) (bool, error) {
	_, ok := m.credentials[username]

	if ok {
		return false, errors.New("user already registered")
	}

	m.credentials[username] = password
	return true, nil
}

func TestLoginCorrect(t *testing.T) {
	server := NewServer(NewMockDB(nil))
	res, _ := server.Login(context.TODO(), &pb.PlayerCredentials{Username: "test", Password: "12345"})

	if res.Result {
		return
	}

	t.Errorf("Login not accepted but should be (correct credentials provided)")
}

func TestLoginUnregisteredUser(t *testing.T) {
	server := NewServer(NewMockDB(nil))
	res, _ := server.Login(context.TODO(), &pb.PlayerCredentials{Username: "foo", Password: "ggggg"})
	if !res.Result {
		return
	}

	t.Errorf("Login accepted but should not be (user not registered)")
}

func TestLoginWrongPassword(t *testing.T) {
	server := NewServer(NewMockDB(nil))
	res, _ := server.Login(context.TODO(), &pb.PlayerCredentials{Username: "test", Password: "not_correct"})

	if !res.Result {
		return
	}

	t.Errorf("Login accepted but should not be (wrong password provided)")
}

func TestRegisterCorrect(t *testing.T) {
	server := NewServer(NewMockDB(nil))
	res, _ := server.Register(context.TODO(), &pb.PlayerDetails{Username: "foo", Password: "correct", Phone: "123456789", Email: "foo@test.com"})

	if res.Result {
		return
	}

	t.Errorf("Register not accepted but should be (correct details provided)")
}

func TestRegisterUserAlreadyRegistered(t *testing.T) {
	server := NewServer(NewMockDB(nil))
	res, _ := server.Register(context.TODO(), &pb.PlayerDetails{Username: "test", Password: "12345", Phone: "123456789", Email: "test@test.com"})

	if !res.Result {
		return
	}

	t.Errorf("Register accepted but should not be (user already registered)")
}
