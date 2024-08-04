package session

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var sessionDirPath string
var sessionFilePath string
var metadataFilePath string

type Session struct {
	Uuid     string        `json:"uuid"`
	Messages []Interaction `json:"history"`
}

type Metadata struct {
	LatestSessionUuid string `json:"latest_session_uuid"`
}

type Interaction struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		panic(err)
	}

	metadataFilePath = filepath.Join(configDir, "llm-cli", "metadata.gob")
	sessionDirPath = filepath.Join(configDir, "llm-cli", "sessions")
	sessionFilePath = filepath.Join(sessionDirPath, "session-%s.gob")

	if err := os.MkdirAll(sessionDirPath, 0755); err != nil {
		fmt.Println("Error creating session directory:", err)
		panic(err)
	}

}

func NewSession() *Session {
	return &Session{
		Uuid:     uuid.New().String(),
		Messages: []Interaction{},
	}
}

func (s *Session) AddMessage(role string, content string) {
	s.Messages = append(s.Messages, Interaction{
		Role:    role,
		Content: content,
	})
}

func createSessionFile(formattedFilePath string) (*os.File, error) {
	file, err := os.Create(formattedFilePath)
	if err != nil {
		fmt.Println("Error creating session file:", err)
		return nil, err
	}
	return file, nil
}

func updateMetadata(latestSessionUuid string) error {
	metadata := Metadata{
		LatestSessionUuid: latestSessionUuid,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(metadata)
	if err != nil {
		fmt.Println("Error encoding metadata:", err)
		return err
	}

	err = os.WriteFile(metadataFilePath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing metadata file:", err)
		return err
	}
	return nil
}

func (s *Session) Save() {
	formattedFilePath := fmt.Sprintf(sessionFilePath, s.Uuid)
	file, err := createSessionFile(formattedFilePath)
	if err != nil {
		fmt.Println("Error creating session file:", err)
		panic(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(file)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(s)
	if err != nil {
		fmt.Println("Error encoding session:", err)
		panic(err)
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error writing session file:", err)
		panic(err)
	}

	err = updateMetadata(s.Uuid)
	if err != nil {
		fmt.Println("Error updating metadata:", err)
		panic(err)
	}

}

func getLatestSessionUuid() (string, error) {
	metadata := Metadata{}
	byteValue, err := os.ReadFile(metadataFilePath)
	if err != nil {
		fmt.Println("Error reading metadata file:", err)
		return "", err
	}
	buf := bytes.NewBuffer(byteValue)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&metadata)
	return metadata.LatestSessionUuid, nil
}

func LoadLatest() *Session {
	if _, err := os.Stat(metadataFilePath); os.IsNotExist(err) {
		fmt.Println("Could not find metadata file. Creating a new session.")
		s := NewSession()
		s.Save()
		return s
	}

	latestSessionUuid, err := getLatestSessionUuid()
	if err != nil {
		fmt.Println("Error getting latest session UUID:", err)
		panic(err)
	}
	latestFilepath := fmt.Sprintf(sessionFilePath, latestSessionUuid)
	s, err := loadFromFilepath(latestFilepath)
	if err != nil {
		fmt.Println("Error loading latest session:", err)
		panic(err)
	}
	return s
}

func loadFromFilepath(filepath string) (*Session, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Could not find session file:", err)
		return nil, err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading session file:", err)
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	s := &Session{}
	err = dec.Decode(s)
	if err != nil {
		fmt.Println("Error decoding session file:", err)
		return nil, err
	}
	return s, nil
}

func ClearSession() *Session {
	newSess := NewSession()
	newSess.Save()
	return newSess
}

func ListSessions() {
	filePaths, err := filepath.Glob(sessionDirPath + "/session-*.gob")
	if err != nil {
		fmt.Println("Error listing session files:", err)
		panic(err)
	}
	latest, _ := getLatestSessionUuid()

	for _, filePath := range filePaths {
		s, err := loadFromFilepath(filePath)
		if err != nil {
			fmt.Println("Error loading session, skipping:", err)
			continue
		}
		if s.Uuid == latest {
			fmt.Printf("* ")
		}
		fmt.Printf("Session %s\nMessage count: %d\n-------------\n", s.Uuid, len(s.Messages))
	}
}

func InspectSession(uuid string, yml bool) {
	session, err := loadFromFilepath(fmt.Sprintf(sessionFilePath, uuid))
	if err != nil {
		fmt.Println("Error loading session:", err)
		panic(err)
	}

	var marshalled []byte
	if yml {
		marshalled, _ = yaml.Marshal(session)
	} else {
		marshalled, _ = json.MarshalIndent(session, "", "  ")
	}
	fmt.Println(string(marshalled))
}

func SwitchSession(uuid string) {
	latestSessionUuid, err := getLatestSessionUuid()
	if err != nil {
		fmt.Println("Error getting latest session UUID:", err)
		panic(err)
	}
	if latestSessionUuid == uuid {
		return
	}
	err = updateMetadata(uuid)
	if err != nil {
		fmt.Println("Could not switch sessions:", err)
		panic(err)
	}
	fmt.Println("Switched to session:", uuid)
	return
}
