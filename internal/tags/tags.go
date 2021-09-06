package tags

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/trento-project/trento/internal/consul"
)

const (
	HostResourceType      = "hosts"
	ClusterResourceType   = "clusters"
	SAPSystemResourceType = "sapsystems"
)

type Tags struct {
	client consul.Client
}

func NewTags(client consul.Client) *Tags {
	return &Tags{
		client: client,
	}
}

func (t *Tags) getKvResourceTagsPath(resourceType string, resourceId string) string {
	return fmt.Sprintf(consul.KvResourceTagsPath, resourceType, resourceId)
}

// GetAllByResource returns all the tags for a given resource
func (t *Tags) GetAllByResource(resourceType string, resourceId string) ([]string, error) {
	path := t.getKvResourceTagsPath(resourceType, resourceId)

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

// Create creates a new tag for a given resource
// The tag is the key of the KV pair to indicate that the tag is present
// The value of the KV pair empty since it is not used
// This simplifies the access to the tags, avoiding the need of loops
func (t *Tags) Create(tag string, resourceType string, resourceId string) error {
	path := fmt.Sprintf("%s%s/", t.getKvResourceTagsPath(resourceType, resourceId), tag)

	if err := t.client.KV().PutMap(path, nil); err != nil {
		return errors.Wrap(err, "Error storing a tag")
	}

	return nil
}

// Delete deletes a tag for a given resource
func (t *Tags) Delete(tag string, resourceType string, resourceId string) error {
	path := fmt.Sprintf("%s%s/", t.getKvResourceTagsPath(resourceType, resourceId), tag)

	_, err := t.client.KV().DeleteTree(path, nil)
	if err != nil {
		return err
	}
	return nil
}
