package datapipeline

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewClustersProjector(db *gorm.DB) *projector {
	clusterProjector := NewProjector("clusters", db)
	clusterProjector.AddHandler(ClusterDiscovery, clustersProjector_ClusterDiscoveryHandler)

	return clusterProjector
}

// TODO: this is a temporary solution, this code needs to be abstracted in the projector.Project() method
func clustersProjector_ClusterDiscoveryHandler(event *DataCollectedEvent, db *gorm.DB) error {
	data, _ := event.Payload.MarshalJSON()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var cluster cluster.Cluster
	if err := dec.Decode(&cluster); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	clusterListReadModel, err := transformClusterData(&cluster)
	if err != nil {
		log.Errorf("can't transform data: %s", err)
		return err
	}

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(clusterListReadModel).Error
}

// transformClusterData transforms the cluster data into the read model
func transformClusterData(cluster *cluster.Cluster) (*models.Cluster, error) {
	if cluster.Id == "" {
		return nil, fmt.Errorf("no cluster ID found")
	}

	return &models.Cluster{
		ID:          cluster.Id,
		Name:        cluster.Name,
		ClusterType: detectClusterType(cluster),
		// TODO: Cost-optimized has multiple SIDs, we will need to implement this in the future
		SIDs:            getHanaSIDs(cluster),
		ResourcesNumber: cluster.Crmmon.Summary.Resources.Number,
		HostsNumber:     cluster.Crmmon.Summary.Nodes.Number,
	}, nil
}

// detectClusterType returns the cluster type based on the cluster resources
func detectClusterType(cluster *cluster.Cluster) string {
	var hasSapHanaTopology, hasSAPHanaController, hasSAPHana bool

	for _, c := range cluster.Crmmon.Clones {
		for _, r := range c.Resources {
			switch r.Agent {
			case "ocf::suse:SAPHanaTopology":
				hasSapHanaTopology = true
			case "ocf::suse:SAPHana":
				hasSAPHana = true
			case "ocf::suse:SAPHanaController":
				hasSAPHanaController = true
			}
		}
	}

	switch {
	case hasSapHanaTopology && hasSAPHana:
		return models.ClusterTypeScaleUp
	case hasSapHanaTopology && hasSAPHanaController:
		return models.ClusterTypeScaleOut
	default:
		return models.ClusterTypeUnknown
	}
}

// getHanaSIDs returns the SIDs of the cluster
// TODO: HANA scale-out has multiple SIDs, we will need to implement this in the future
func getHanaSIDs(c *cluster.Cluster) []string {
	var sids []string
	for _, r := range c.Cib.Configuration.Resources.Clones {
		if r.Primitive.Type == "SAPHanaTopology" {
			for _, a := range r.Primitive.InstanceAttributes {
				if a.Name == "SID" && a.Value != "" {
					sids = append(sids, a.Value)
				}
			}
		}
	}

	return sids
}
