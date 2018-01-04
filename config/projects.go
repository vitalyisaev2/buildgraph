package config

import (
	"fmt"

	"github.com/vitalyisaev2/buildgraph/graph"
)

// ProjectsConfig describes projects and their relations
type ProjectsConfig struct {
	Descriptions []*Description      `yaml:"descriptions"`
	Relations    map[string][]string `yaml:"relations"`
}

// Description contains basic information about the project
type Description struct {
	ID        string `yaml:"id"`        // project's unique identifier (used to describe project relations)
	Namespace string `yaml:"namespace"` // namespace that project belongs to
	Name      string `yaml:"name"`      // project's name
}

func (c *ProjectsConfig) validate() error {
	if len(c.Descriptions) == 0 {
		return fmt.Errorf("empty project descriptions")
	}

	for _, d := range c.Descriptions {
		if d.ID == "" || d.Name == "" || d.Namespace == "" {
			return fmt.Errorf("invalid project description: %v", d)
		}
	}

	if len(c.Relations) == 0 {
		return fmt.Errorf("empty project relations")
	}

	g, err := graph.NewGraphFromAdjacencyMap(c.Relations)
	if err != nil {
		return err
	}

	cyclic, cycle, err := g.Cyclic()
	if cyclic {
		return fmt.Errorf("project dependency graph contains cycle: %s", cycle.String())
	}
	if err != nil {
		return err
	}

	return nil
}
