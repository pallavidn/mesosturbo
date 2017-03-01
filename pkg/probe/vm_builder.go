package probe

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"

	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/golang/glog"
)


const PROXY_VM_IP string = "Proxy_VM_IP"
// Builder for creating VM Entities to represent the Mesos Agents or Slaves in Turbo server
// This will create a proxy VM in the server. Hypervisor probes in the server will discover and manage the Agent VMs
type VMEntityBuilder struct {
	MasterState   *util.MesosAPIResponse
}

// Build VM EntityDTO using the agent listed in the 'state' json returned from the Mesos Master
func (nb *VMEntityBuilder) BuildEntities() ([]*proto.EntityDTO, error) {
	glog.V(2).Infof("[BuildEntities] ...... ")
	nodes := nb.MasterState.Slaves
	result := []*proto.EntityDTO{}
	// For each agent
	for _, slave := range nodes {
		commoditiesSold, _ := nb.vmCommSold(&slave)
		slaveIP := nb.MasterState.SlaveIdIpMap[slave.Id]
		entityDTO, _ := nb.vmEntity(slaveIP, slave.Id, slave.Name, commoditiesSold)
		result = append(result, entityDTO)
	}
	//glog.Infof("[BuildEntities] Entity DTOs :", result)
	return result, nil
}

// Build VM EntityDTO
func (nb *VMEntityBuilder) vmEntity(slaveIP, nodeID, displayName string,
					commoditiesSold []*proto.CommodityDTO) (*proto.EntityDTO, error) {
	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE, nodeID).
					DisplayName(displayName).
					SellsCommodities(commoditiesSold)
	// Stitching and proxy metadata
	ipAddress := slaveIP
	ipPropName := PROXY_VM_IP	// We create a different property for ip address, so the IP object in the server entity
					// is not deleted during reconciliation
					// TODO: create a builder for proxy VMs
	ipProp := &proto.EntityDTO_EntityProperty{
		Namespace: &DEFAULT_NAMESPACE,
		Name: &ipPropName,
		Value: &ipAddress,
	}
	entityDTOBuilder = entityDTOBuilder.WithProperty(ipProp)
	glog.Infof("[NodeBuilder] Parse node: The ip of vm to be reconcile with is %s", ipAddress)
	metaData := generateReconciliationMetaData()

	entityDTOBuilder = entityDTOBuilder.ReplacedBy(metaData)
	entityDto, _ := entityDTOBuilder.Create()
	return entityDto, nil
}

func (nb *VMEntityBuilder) vmCommSold(slaveInfo *util.Slave) ([]*proto.CommodityDTO, error) {
	var commoditiesSold []*proto.CommodityDTO
	// MemProv
	memProvComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM_PROVISIONED).
				Capacity(1.0).
				Create()
	commoditiesSold = append(commoditiesSold, memProvComm)
	// CpuProv
	cpuProvComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU_PROVISIONED).
				Capacity(1.0).
				Create()
	commoditiesSold = append(commoditiesSold, cpuProvComm)
	// VMem
	vMemComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).
				Capacity(1.0).
				Create()
	commoditiesSold = append(commoditiesSold, vMemComm)
	// VCpu
	vCpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).
				Capacity(1.0).
				Create()
	commoditiesSold = append(commoditiesSold, vCpuComm)

	// Access Commodities
	// ClusterComm
	clusterComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CLUSTER).
					Key(nb.MasterState.ClusterName).
					Create()
	commoditiesSold = append(commoditiesSold, clusterComm)

	// TODO add port commodity sold

	return commoditiesSold, nil
}


func generateReconciliationMetaData() *proto.EntityDTO_ReplacementEntityMetaData {
	replacementEntityMetaDataBuilder := builder.NewReplacementEntityMetaDataBuilder()
	replacementEntityMetaDataBuilder.Matching(PROXY_VM_IP)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_CPU_PROVISIONED).
					PatchSelling(proto.CommodityDTO_MEM_PROVISIONED).
					PatchSelling(proto.CommodityDTO_CLUSTER).
					PatchSelling(proto.CommodityDTO_VCPU).
					PatchSelling(proto.CommodityDTO_VMEM).
					PatchSelling(proto.CommodityDTO_VMPM_ACCESS)
	metaData := replacementEntityMetaDataBuilder.Build()
	return metaData
}
