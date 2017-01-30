package probe

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"

type EntityBuilder interface {

	BuildEntities(*proto.EntityDTO_EntityType) ([]*proto.EntityDTO, error)
}
