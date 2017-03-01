package probe

import (

	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/golang/glog"
)

// Builder for creating Application Entities to represent the Mesos Tasks in Turbo server
type AppEntityBuilder struct {
	MasterState   *util.MesosAPIResponse
}

// Build Application EntityDTO using the tasks listed in the 'state' json returned from the Mesos Master
func (tb *AppEntityBuilder) BuildEntities() ([]*proto.EntityDTO, error) {
	glog.V(2).Infof("[BuildEntities] ...... ")
	result := []*proto.EntityDTO{}
	taskList := tb.MasterState.TaskMasterAPI.Tasks

	// builder.BuildConstraintMap(m.MApps.Apps)	//TODO: what is this for? uses framework apps
	for _, task := range taskList {
		// fmt.Println("[AppEntityBuilder] Task ====> %+v", task)
		//fmt.Println("[AppEntityBuilder] task is %s and state %s ", task.Name, task.State)
		if task.State != "TASK_RUNNING" {
			continue
		}

		// builder.SetTaskConstraints(taskProbe)	//TODO:
		// taskResource, err := taskProbe.GetTaskResourceStat(m.MapTaskStatistics, taskProbe.Task, taskUseMap) //TODO:
		// Commodities sold
		commoditiesSoldApp := tb.appCommsSold(&task)
		// Application Entity for the task
		entityDTOBuilder := tb.appEntityDTO(&task, commoditiesSoldApp)
		// Commodities bought
		entityDTOBuilder = tb.appCommoditiesBought(entityDTOBuilder, &task)
		// Entity DTO
		entityDTO, _ := entityDTOBuilder.Create()
		result = append(result, entityDTO)
	}
	//glog.Infof("[BuildEntities] Task entity DTOs :", result)
	return result, nil
}

// Build Application DTO
func (tb *AppEntityBuilder) appEntityDTO(task *util.Task, commoditiesSold []*proto.CommodityDTO) *builder.EntityDTOBuilder {
	appEntityType := proto.EntityDTO_APPLICATION
	id := "APP:" + task.Name + "-" + task.Id	//task.Name + "::" + "APP:" + task.Id
	dispName := "APP:" + task.Name
	entityDTOBuilder := builder.NewEntityDTOBuilder(appEntityType, id).	//+"foo").
					DisplayName(dispName)

	return entityDTOBuilder
}

// Build commodityDTOs for commodity sold by the app
func  (tb *AppEntityBuilder) appCommsSold(task *util.Task) []*proto.CommodityDTO {
	var commoditiesSold []*proto.CommodityDTO
	transactionComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).
					Key(task.Name).
					Create()
	commoditiesSold = append(commoditiesSold, transactionComm)
	return commoditiesSold
}

// Build commodityDTOs for commodity bought by the app
func (tb *AppEntityBuilder) appCommoditiesBought(appDto *builder.EntityDTOBuilder, task *util.Task) *builder.EntityDTOBuilder {
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

	applicationComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).
				Key(task.Id).
				Create()
	commoditiesBought = append(commoditiesBought, applicationComm)

	// From Container
	containerName := task.Id
	containerProvider := builder.CreateProvider(proto.EntityDTO_CONTAINER, containerName)
	appDto.Provider(containerProvider)
	appDto.BuysCommodities(commoditiesBought)

	return appDto
}
