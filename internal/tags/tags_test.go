package tags

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestGetAll(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst)

	listMap := map[string]interface{}{
		HostResourceType: map[string]interface{}{
			"hostname1": map[string]interface{}{
				"tag1": struct{}{},
				"tag2": struct{}{},
				"tag3": struct{}{},
			},
			"hostname2": map[string]interface{}{
				"tag4": struct{}{},
				"tag5": struct{}{},
				"tag6": struct{}{},
			},
		},
		ClusterResourceType: map[string]interface{}{
			"cluster_id_1": map[string]interface{}{
				"tag1": struct{}{},
				"tag2": struct{}{},
				"tag3": struct{}{},
			},
			"cluster_id_2": map[string]interface{}{
				"tag4": struct{}{},
				"tag5": struct{}{},
				"tag6": struct{}{},
			},
		},
		SAPSystemResourceType: map[string]interface{}{
			"HA1": map[string]interface{}{
				"tag4": struct{}{},
				"tag5": struct{}{},
				"tag6": struct{}{},
			},
		},
	}
	kv.On("ListMap", consul.KvTagsPath, consul.KvTagsPath).Return(listMap, nil)

	resTags, _ := tags.GetAll()
	assert.Equal(t, []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"}, resTags)

	resTags, _ = tags.GetAll("sapsystems")
	assert.Equal(t, []string{"tag4", "tag5", "tag6"}, resTags)

	resTags, _ = tags.GetAll("systems", "hosts")
	assert.Equal(t, []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"}, resTags)
}

func TestGetAllByResource(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst)

	listMap := map[string]interface{}{
		"tag1": struct{}{},
		"tag2": struct{}{},
	}
	kv.On("ListMap", tags.getKvResourceTagsPath("res", "id"), tags.getKvResourceTagsPath("res", "id")).Return(listMap, nil)

	resTags, _ := tags.GetAllByResource("res", "id")

	assert.Equal(t, []string{"tag1", "tag2"}, resTags)
}

func TestCreate(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst)

	kv.On("PutMap", tags.getKvResourceTagsPath("res", "id")+"createme/", map[string]interface{}(nil)).Return(nil, nil)

	err := tags.Create("createme", "res", "id")

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	tags := NewTags(consulInst)

	kv.On("DeleteTree", tags.getKvResourceTagsPath("res", "id")+"deleteme/", (*consulApi.WriteOptions)(nil)).Return(nil, nil)

	err := tags.Delete("deleteme", "res", "id")

	assert.NoError(t, err)
	kv.AssertExpectations(t)
}
