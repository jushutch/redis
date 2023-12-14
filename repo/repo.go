package repo

import "fmt"

type Repo struct {
	data map[string]string
}

func NewRepo() *Repo {
	return &Repo{
		data: make(map[string]string),
	}
}

func (r *Repo) Set(key, value string) error {
	fmt.Printf("Set Key: %q\nSet Value: %q\n", key, value)
	r.data[key] = value
	return nil
}

func (r *Repo) Get(key string) (string, error) {
	value, ok := r.data[key]
	if !ok {
		return "", fmt.Errorf("the key %s does not exist", key)
	}
	return value, nil
}
