package server_usecase

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type leafRoute struct {
	Method string
	Path   string
}

type RouteList struct {
	app *fiber.App
}

func NewRouteList(app *fiber.App) *RouteList {
	return &RouteList{
		app: app,
	}
}

func (u *RouteList) Get() []string {
	var routes []string

	for _, route := range u.getLeafRoutes() {
		routes = append(routes, route.Method+" "+route.Path)
	}

	return routes
}

func (u *RouteList) getLeafRoutes() []leafRoute {
	var all []leafRoute
	var leafRoutes []leafRoute

	for _, routes := range u.app.Stack() {
		for _, r := range routes {
			if r.Path == "/" || r.Path == "/api/" || r.Method == "HEAD" {
				continue
			}

			all = append(all, leafRoute{
				Method: r.Method,
				Path:   r.Path,
			})
		}
	}

	seen := make(map[string]struct{})

	for _, r := range all {
		if u.isLeaf(all, r.Path) {
			key := r.Method + "|" + r.Path
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				leafRoutes = append(leafRoutes, r)
			}
		}
	}

	return leafRoutes
}

func (u *RouteList) isLeaf(all []leafRoute, path string) bool {
	for _, other := range all {
		isSubPath := path != other.Path &&
			strings.HasPrefix(other.Path, path) &&
			len(other.Path) > len(path) &&
			other.Path[len(path)] == '/'

		if isSubPath {
			return false
		}
	}

	return true
}
