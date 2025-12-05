package portfolio

import (
	"gopkg.in/yaml.v3"
)

// Data represents the root structure of data.yaml
type Data struct {
	Name           string          `yaml:"name"`
	Title          string          `yaml:"title"`
	Location       string          `yaml:"location"`
	Phone          string          `yaml:"phone"`
	Email          string          `yaml:"email"`
	LinkedIn       string          `yaml:"linkedin"`
	Portfolio      string          `yaml:"portfolio"`
	Summary        string          `yaml:"summary"`
	Experience     []Experience    `yaml:"experience"`
	Skills         Skills          `yaml:"skills"`
	Education      []Education     `yaml:"education"`
	Certifications []Certification `yaml:"certifications"`
}

// Experience represents a single work experience entry
type Experience struct {
	Title        string   `yaml:"title"`
	Company      string   `yaml:"company"`
	Location     string   `yaml:"location"`
	StartDate    string   `yaml:"start_date"`
	EndDate      string   `yaml:"end_date"`
	Current      bool     `yaml:"current"`
	Highlights   []string `yaml:"highlights"`
	Technologies string   `yaml:"technologies"`
}

// Skills represents all skill categories
type Skills struct {
	Languages   []string `yaml:"languages"`
	Frontend    []string `yaml:"frontend"`
	Backend     []string `yaml:"backend"`
	CloudDevOps []string `yaml:"cloud_devops"`
	Databases   []string `yaml:"databases"`
	Tools       []string `yaml:"tools"`
}

// Education represents an education entry
type Education struct {
	School     string `yaml:"school"`
	Degree     string `yaml:"degree"`
	Status     string `yaml:"status"`
	Graduation string `yaml:"graduation"`
}

// Certification represents a certification entry
type Certification struct {
	Name         string `yaml:"name"`
	Issued       string `yaml:"issued"`
	Expires      string `yaml:"expires"`
	CredentialID string `yaml:"credential_id"`
}

// LoadData unmarshals the provided YAML bytes into a Data struct
func LoadData(content []byte) (*Data, error) {
	var data Data
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
