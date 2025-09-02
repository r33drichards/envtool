package envfile

import (
	"bufio"
	"os"
	"strings"
)

// Parser handles reading and parsing .env files
type Parser interface {
	Parse(path string) (map[string]string, error)
}

// DefaultParser implements the Parser interface
type DefaultParser struct{}

// Parse reads and parses a .env file at the given path
func (p *DefaultParser) Parse(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Split on first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		
		envVars[key] = value
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return envVars, nil
}