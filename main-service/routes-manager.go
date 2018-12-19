package mainservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Route type declaration
type Route struct {
	Route_name        string         `json:"route_name"`
	Related_app       string         `json:"related_app"`
	Route_description string         `json:"route_description"`
	Route_body        interface{}    `json:"route_body"`
	Route_params      []RouteParams  `json:"route_params"`
	Route_queries     []RouteQueries `json:"route_queries"`
}

type RouteParams struct {
	Param       string
	Description string
}

type RouteQueries struct {
	Query       string
	Description string
}

// Returns filtered routes from json file
func GetRoutesData(appType string) *[]Route {
	var routes []Route
	var filteredRoutes []Route

	bFile, err := ioutil.ReadFile("routes-info.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(bFile, &routes)
	if err != nil {
		fmt.Println(err)
	}

	if len(appType) == 0 {
		return &routes
	}

	for _, route := range routes {
		if route.Related_app == appType {
			filteredRoutes = append(filteredRoutes, route)
		}
	}

	return &filteredRoutes
}
