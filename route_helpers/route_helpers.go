package route_helpers

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/receptor"
)

const AppRouter = "cf-router"

//const AppRouter = "lattice-router"

type AppRoutes []AppRoute

type AppRoute struct {
	Hostnames []string `json:"hostnames"`
	Port      uint16   `json:"port"`
}

func (l AppRoutes) RoutingInfo() receptor.RoutingInfo {
	data, _ := json.Marshal(l)
	routingInfo := json.RawMessage(data)
	return receptor.RoutingInfo{
		AppRouter: &routingInfo,
	}
}

func AppRoutesFromRoutingInfo(routingInfo receptor.RoutingInfo) (AppRoutes, error) {
	if routingInfo == nil {
		return nil, nil
	}

	data, found := routingInfo[AppRouter]
	if !found {
		return nil, nil
	}

	if data == nil {
		return nil, nil
	}

	routes := AppRoutes{}
	err := json.Unmarshal(*data, &routes)

	return routes, err
}
