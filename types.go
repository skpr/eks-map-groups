package main

// MapGroup works as a lookup template to build a list of MapUsers.
type MapGroup struct {
	Name     string   `yaml:"name"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

// MapUser is the resulting user who can access the cluster.
type MapUser struct {
	UserARN  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}
