package tags

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/trento-project/trento/internal/consul"
)

type Tags struct {
	client   consul.Client
	resource string
	id       string
}

func NewTags(client consul.Client, resource string, id string) *Tags {
	return &Tags{
		client:   client,
		resource: resource,
		id:       id,
	}
}

func (t *Tags) getKvTagsPath() string {
	return fmt.Sprintf(consul.KvTagsPath, t.resource, t.id)
}

// GetAll returns all the tags
func (t *Tags) GetAll() ([]string, error) {
	path := t.getKvTagsPath()

	tagsMap, err := t.client.KV().ListMap(path, path)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving tags")
	}

	var tags []string
	for tag := range tagsMap {
		tags = append(tags, tag)
	}

	return tags, nil
}

// Creeate creates a new tag
// The tag is the key of the KV pair to indicate that the tag is present
// The value of the KV pair empty since it is not used
// This simplifies the access to the tags, avoiding the need of loops
func (t *Tags) Create(tag string) error {
	path := fmt.Sprintf("%s%s/", t.getKvTagsPath(), tag)

	if err := t.client.KV().PutMap(path, nil); err != nil {
		return errors.Wrap(err, "Error storing a tag")
	}

	return nil
}

// Delete deletes a tag
func (t *Tags) Delete(tag string) error {
	path := fmt.Sprintf("%s%s/", t.getKvTagsPath(), tag)

	_, err := t.client.KV().DeleteTree(path, nil)
	if err != nil {
		return err
	}
	return nil
}
