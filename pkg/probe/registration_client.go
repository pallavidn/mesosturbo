package probe

import (
	"fmt"

	// turbo sdk imports
	supplychain "github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
	proto "github.com/turbonomic/turbo-go-sdk/pkg/proto"
	builder "github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/mesosturbo/pkg/conf"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
)


// Registration Client for the Mesos Probe
// Implements the TurboRegistrationClient interface
type MesosRegistrationClient struct {
	mesosMasterType 	 conf.MesosMasterType
}

func NewRegistrationClient (mesosMasterType conf.MesosMasterType) probe.TurboRegistrationClient {
	client := &MesosRegistrationClient{
		mesosMasterType: mesosMasterType,
	}
	return client
}

var (
	vCpuType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_VCPU
	vMemType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_VMEM
	cpuAllocType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_CPU_ALLOCATION
	memAllocType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_MEM_ALLOCATION
	appType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_APPLICATION
	clusterType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_CLUSTER
	diskAllocationType proto.CommodityDTO_CommodityType 	= proto.CommodityDTO_STORAGE_ALLOCATION
	networkType proto.CommodityDTO_CommodityType 		= proto.CommodityDTO_NETWORK
	transactionType proto.CommodityDTO_CommodityType 	= proto.CommodityDTO_TRANSACTION

	//Commodity key is optional, when key is set, it serves as a constraint between seller and buyer
	//for example, the buyer can only go to a seller that sells the commodity with the required key
	fakeKey string = "fake"
	emptyKey string = ""

	vCpuTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vCpuType}
	vMemTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vMemType}
	vCpuTemplateCommWithEmptyKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vCpuType, Key: &emptyKey}
	vMemTemplateCommWithEmptyKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vMemType, Key: &emptyKey}
	cpuAllocTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &cpuAllocType, Key: &fakeKey}
	memAllocTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &memAllocType, Key: &fakeKey}

	appTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &appType,  Key: &fakeKey}
	clusterTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &clusterType, Key: &fakeKey}
	transactionTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &transactionType, Key: &fakeKey}
)

func (registrationClient *MesosRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	fmt.Println("[MesosRegistrationClient] .......... Now use builder to create a supply chain ..........")

	// VM Node
	slaveSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_MACHINE).
						Sells(cpuAllocTemplateCommWithKey).
						Sells(memAllocTemplateCommWithKey).
						//		Selling(sdk.CommodityDTO_STORAGE_ALLOCATION, fakeKey).
						Sells(vCpuTemplateCommWithEmptyKey).
						Sells(vMemTemplateCommWithEmptyKey).
						Sells(appTemplateCommWithKey).
						Sells(clusterTemplateCommWithKey)

	// Container Node
	containerSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_CONTAINER).
							//		Selling(sdk.CommodityDTO_STORAGE_ALLOCATION, fakeKey).
						Sells(cpuAllocTemplateCommWithKey).
						Sells(memAllocTemplateCommWithKey)

	// Container Node to VM Link
	containerSupplyChainNodeBuilder = containerSupplyChainNodeBuilder.
						Provider(proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_LAYERED_OVER).
						Buys(cpuAllocTemplateCommWithKey).
						Buys(memAllocTemplateCommWithKey).
						//		Buys(*diskAllocationTemplateComm).
						Buys(clusterTemplateCommWithKey)

	// Application Node
	appSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_APPLICATION).
						Sells(transactionTemplateCommWithKey)

	// Application Node to Container Link
	appSupplyChainNodeBuilder = appSupplyChainNodeBuilder.
						Provider(proto.EntityDTO_CONTAINER, proto.Provider_LAYERED_OVER).
						//		Buys(*appDiskAllocationTemplateComm).
						Buys(cpuAllocTemplateCommWithKey).
						Buys(memAllocTemplateCommWithKey)

	// Application Node to VM Link
	appSupplyChainNodeBuilder = appSupplyChainNodeBuilder.
						Provider(proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING).
						Buys(vCpuTemplateComm).
						Buys(vMemTemplateComm).
						Buys(appTemplateCommWithKey)

	// External Link from Container (Pod) to VM
	vmContainerExtLinkBuilder := supplychain.NewExternalEntityLinkBuilder().
			Link(proto.EntityDTO_CONTAINER, proto.EntityDTO_VIRTUAL_MACHINE,
					proto.Provider_LAYERED_OVER).
			Commodity(cpuAllocType, true).
			Commodity(memAllocType, true).
			Commodity(diskAllocationType, true).
			Commodity(clusterType, true).
			Commodity(networkType, true).
			ProbeEntityPropertyDef(supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS,
						"IP Address where the Container is running").
			ExternalEntityPropertyDef(supplychain.VM_IP)

	vmContainerExternalLink, _ := vmContainerExtLinkBuilder.Build()
	slaveSupplyChainNodeBuilder.ConnectsTo(vmContainerExternalLink)

	// Link from Application to VM
	vmAppExtLinkBuilder := supplychain.NewExternalEntityLinkBuilder().
			Link(proto.EntityDTO_APPLICATION, proto.EntityDTO_VIRTUAL_MACHINE,
					proto.Provider_HOSTING).
			Commodity(vCpuType, false).
			Commodity(vMemType, false).
			Commodity(appType, true).
			ProbeEntityPropertyDef(supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS,
							"IP Address where the Application is running").
			ExternalEntityPropertyDef(supplychain.VM_IP)

	vmAppExternalLink, _ := vmAppExtLinkBuilder.Build()
	slaveSupplyChainNodeBuilder.ConnectsTo(vmAppExternalLink)

	appNode, _ := appSupplyChainNodeBuilder.Create()
	containerNode, _ := containerSupplyChainNodeBuilder.Create()
	vmNode, _ := slaveSupplyChainNodeBuilder.Create()

	supplyChainBuilder := supplychain.NewSupplyChainBuilder()
	supplyChainBuilder.
			Top(appNode).
			Entity(containerNode).
			Entity(vmNode)

	supplychain, _ := supplyChainBuilder.Create()
	return supplychain
}

func (registrationClient *MesosRegistrationClient) GetIdentifyingFields() string {
	return string(MasterIP)
}

// The return type is a list of ProbeInfo_AccountDefProp.
// For a valid definition, targetNameIdentifier, username and password should be contained.
// Account Definition for Mesos Probe
func (registrationClient *MesosRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry

	// master ip
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(string(MasterIP), string(MasterIP),	 //"MasterIP",
		"IP of the mesos master", ".*",
		true, false).
		Create()
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)
	// master port
	masterPortAcctDefEntry := builder.NewAccountDefEntryBuilder(string(MasterPort), string(MasterPort),	//"MasterPort",
		"Port of the mesos master", ".*",
		false, false).
		Create()
	acctDefProps = append(acctDefProps, masterPortAcctDefEntry)

	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder( string(Username), string(Username),		//"Username",
									"Username of the mesos master", ".*",
									false, false).
					Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)

	// password
	passwdAcctDefEntry := builder.NewAccountDefEntryBuilder(string(Password), string(Password),		//"Password",
								"Password of the mesos master", ".*",
								false, true).
					Create()
	acctDefProps = append(acctDefProps, passwdAcctDefEntry)

	// TODO: Should contain the fields required to connect to a Mesos Agent, same as the ones in the MesosTargetConf
	if registrationClient.mesosMasterType == conf.Apache {
		// framework id
		frameworkIpAcctDefEntry := builder.NewAccountDefEntryBuilder(string(FrameworkIP), string(FrameworkIP), //"FrameworkIP",
			"IP for the Framework", ".*",
			false, false).
			Create()
		acctDefProps = append(acctDefProps, frameworkIpAcctDefEntry)

		// framework port
		frameworkPortAcctDefEntry := builder.NewAccountDefEntryBuilder(string(FrameworkPort), string(FrameworkPort),	//"FrameworkPort",
			"Port for the Framework", ".*",
			false, false).
			Create()
		acctDefProps = append(acctDefProps, frameworkPortAcctDefEntry)


		// username
		frameworkUserAcctDefEntry := builder.NewAccountDefEntryBuilder(string(FrameworkUsername), string(FrameworkUsername),	//"FrameworkUsername",
			"Username for the framework", ".*",
			false, false).
			Create()
		acctDefProps = append(acctDefProps, frameworkUserAcctDefEntry)

		// password
		frameworkPwdAcctDefEntry := builder.NewAccountDefEntryBuilder(string(FrameworkPassword), string(FrameworkPassword),	//"FrameworkPassword",
			"Password for the framework", ".*",
			false, true).
			Create()
		acctDefProps = append(acctDefProps, frameworkPwdAcctDefEntry)
	}

	// action ip
	actionIPAcctDefEntry := builder.NewAccountDefEntryBuilder(string(ActionIP), string(ActionIP),
		"IP of the action executor framework", ".*",
		false, false).
		Create()
	acctDefProps = append(acctDefProps, actionIPAcctDefEntry)
	// action port
	actionPortAcctDefEntry := builder.NewAccountDefEntryBuilder(string(ActionPort), string(ActionPort),
		"Port of the action executor framework", ".*",
		false, false).
		Create()
	acctDefProps = append(acctDefProps, actionPortAcctDefEntry)

	return acctDefProps
}

// TODO: change the acct def depending on the mesos or dcos target