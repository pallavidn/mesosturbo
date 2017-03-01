package probe

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

// Builder for creating Container Entities to represent the default container Mesos Tasks in Turbo server
type ContainerEntityBuilder struct {
	MasterState   *util.MesosAPIResponse
}

// Build Container EntityDTO using the tasks listed in 'state' json returned from the Mesos Master
func (cb *ContainerEntityBuilder) BuildEntities() ([]*proto.EntityDTO, error) {
	glog.V(2).Infof("[BuildEntities] ....")
	result := []*proto.EntityDTO{}
	taskList := cb.MasterState.TaskMasterAPI.Tasks
	for _, task := range taskList {
		// fmt.Println("[ContainerEntityBuilder] task is %s and state %s ", task.Name, task.State)
		// skip non running tasks
		if task.State != "TASK_RUNNING" {
			continue
		}

		//taskResource, err := taskProbe.GetTaskResourceStat(m.MapTaskStatistics, taskProbe.Task, taskUseMap) // TODO: monitoring related
		commoditiesSoldContainer := cb.containerCommsSold(&task)
		entityDTOBuilder, _ := cb.buildContainerEntityDTO(&task, commoditiesSoldContainer)
		entityDTOBuilder = cb.containerCommoditiesBought(entityDTOBuilder, &task)
		entityDTO, _ := entityDTOBuilder.Create()
		result = append(result, entityDTO)
	}
	//glog.Infof("[BuildEntities] Container entity DTOs :", result)
	return result, nil
}


// Build Container entityDTO
func (cb *ContainerEntityBuilder) buildContainerEntityDTO(task *util.Task, commoditiesSold []*proto.CommodityDTO) (*builder.EntityDTOBuilder, error) {
	id := task.Id
	dispName := task.Name

	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_CONTAINER, id).
					DisplayName(dispName)

	slaveId := task.SlaveId
	if slaveId == "" {
		return nil, fmt.Errorf("Cannot find the hosting slave for task %s", dispName)
	}
	glog.V(2).Infof("Pod %s is hosted on %s", dispName, slaveId)

	entityDTOBuilder.SellsCommodities(commoditiesSold)

	////	providerUid := nodeUidTranslationMap[slaveId]
	slaveIdIpMap := cb.MasterState.SlaveIdIpMap
	ipAddress := slaveIdIpMap[task.SlaveId]
	ipPropName := supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS
	ipProp := &proto.EntityDTO_EntityProperty {	// TODO: create Property Builder
		Namespace: &DEFAULT_NAMESPACE,
		Name: &ipPropName,
		Value: &ipAddress,
	}
	entityDTOBuilder = entityDTOBuilder.WithProperty(ipProp)
	glog.V(2).Infof("Pod %s will be stitched to VM with IP %s", dispName, ipAddress)

	return entityDTOBuilder, nil
}

// Build commodityDTOs for commodity sold by the pod
func  (tb *ContainerEntityBuilder) containerCommsSold(task *util.Task) []*proto.CommodityDTO {
	var commoditiesSold []*proto.CommodityDTO
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
	// Application with task id as the key
	applicationComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).
					Key(task.Id).
					Create()
	commoditiesSold = append(commoditiesSold, applicationComm)
	return commoditiesSold
}


// Build commodityDTOs for commodity sold by the pod
func (tb *ContainerEntityBuilder) containerCommoditiesBought(containerDto *builder.EntityDTOBuilder,
								task *util.Task) *builder.EntityDTOBuilder {
	var commoditiesBought []*proto.CommodityDTO
	// VMem
	vMemComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).
				Used(1.0).
				Create()
	commoditiesBought = append(commoditiesBought, vMemComm)
	// VCpu
	vCpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).
				Used(1.0).
				Create()
	commoditiesBought = append(commoditiesBought, vCpuComm)
	// MemProv
	memProvComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM_PROVISIONED).
				Used(1.0).
				Create()
	commoditiesBought = append(commoditiesBought, memProvComm)
	// CpuProv
	cpuProvComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU_PROVISIONED).
				Used(1.0).
				Create()
	commoditiesBought = append(commoditiesBought, cpuProvComm)

	clusterCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CLUSTER).
					Key(tb.MasterState.Cluster.ClusterName).
					Create()
	commoditiesBought = append(commoditiesBought, clusterCommBought)

	// TODO other constraint operator types
	//taskProbe.getPortsBought()	// TODO:
	////glog.V(3).Infof("\n\n\n")
	//for k, _ := range taskProbe.PortsUsed {
	//	//glog.V(3).Infof(" -------->>>> ports used by task %+v  and  %+v \n", k, v)
	//	networkCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_NETWORK).
	//		Key(k).
	//		Create()
	//	commoditiesBought = append(commoditiesBought, networkCommBought)
	//
	//}

	providerDto := builder.CreateProvider(proto.EntityDTO_VIRTUAL_MACHINE, task.SlaveId)
	containerDto.Provider(providerDto)
	containerDto.BuysCommodities(commoditiesBought)

	return containerDto
}