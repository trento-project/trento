package tags

import (
	"sort"
	"testing"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestGetAll(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst, "res", "id")

	listMap := map[string]interface{}{
		"tag1": struct{}{},
		"tag2": struct{}{},
	}
	kv.On("ListMap", tags.getKvTagsPath(), tags.getKvTagsPath()).Return(listMap, nil)

	resTags, _ := tags.GetAll()
	sort.Strings(resTags)

	assert.Equal(t, []string{"tag1", "tag2"}, resTags)
}

func TestCreate(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst, "res", "id")

	kv.On("PutMap", tags.getKvTagsPath()+"createme/", map[string]interface{}(nil)).Return(nil, nil)

	err := tags.Create("createme")

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst, "res", "id")

	kv.On("DeleteTree", tags.getKvTagsPath()+"deleteme/", (*consulApi.WriteOptions)(nil)).Return(nil, nil)

	err := tags.Delete("deleteme")

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}
