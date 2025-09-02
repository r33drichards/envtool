package shell

import (
	"io/ioutil"
	"os"
	"strings"
)

// FileManager handles reading and writing shell configuration files
type FileManager interface {
	AppendToFile(path, content string) error
	ReadFile(path string) (string, error)
	FileExists(path string) (bool, error)
	WriteFile(path, content string) error
}

// DefaultFileManager implements the FileManager interface
type DefaultFileManager struct{}

// AppendToFile appends content to a file
func (fm *DefaultFileManager) AppendToFile(path, content string) error {
	// Check if the content is already in the file to avoid duplicates
	exists, err := fm.ContainsContent(path, content)
	if err != nil {
		return err
	}
	if exists {
		// Content already exists, no need to append
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	if _, err = f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// ReadFile reads content from a file
func (fm *DefaultFileManager) ReadFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ContainsContent checks if a file contains specific content
func (fm *DefaultFileManager) ContainsContent(path, content string) (bool, error) {
	fileContent, err := fm.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	
	return strings.Contains(fileContent, content), nil
}

// FileExists checks if a file exists
func (fm *DefaultFileManager) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// WriteFile writes content to a file
func (fm *DefaultFileManager) WriteFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}