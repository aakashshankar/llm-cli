package session

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var tempDir = os.TempDir()
var sessionFilePath = tempDir + "/session.json"

type Session struct {
	Messages []Interaction `json:"history"`
}

type Interaction struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSession() *Session {
	return &Session{
		Messages: []Interaction{},
	}
}

func (s *Session) AddMessage(role string, content string) {
	s.Messages = append(s.Messages, Interaction{
		Role:    role,
		Content: content,
	})
}

func (s *Session) Save() error {
	file, err := os.Create(sessionFilePath)
	if err != nil {
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) LoadLatest() error {
	if _, err := os.Stat(sessionFilePath); os.IsNotExist(err) {
		// create a new session
		return nil
	}

	f, err := os.Open(sessionFilePath)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, s)
	if err != nil {
		return err
	}

	return nil
}

func ClearSession() {
	err := os.Remove(sessionFilePath)
	if err != nil {
		return
	}
}

func ListSession() {
	if _, err := os.Stat(sessionFilePath); os.IsNotExist(err) {
		// create a new session
		return
	}

	f, err := os.Open(sessionFilePath)
	if err != nil {
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	byteValue, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading session file:", err)
		os.Exit(1)
	}
	var session Session
	err = json.Unmarshal(byteValue, &session)
	if err != nil {
		fmt.Println("Error unmarshalling session file:", err)
		os.Exit(1)
	}

	byteValue, err = json.MarshalIndent(session, "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(byteValue))
}
