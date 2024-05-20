package api

import (
	"encoding/json"
	"io"
	"os"
)

type Session struct {
	History []Interaction `json:"history"`
}

type Interaction struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSession() *Session {
	return &Session{
		History: []Interaction{},
	}
}

func (s *Session) AddInteraction(role string, content string) {
	s.History = append(s.History, Interaction{
		Role:    role,
		Content: content,
	})
}

func (s *Session) Save(filename string) error {
	data, err := json.Marshal(filename)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(data))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Load(filename string) (*Session, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return NewSession(), nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
