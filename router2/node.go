package router2

// stolen from  https://raw.github.com/gocraft/web/master/tree.go and modified

import (
	"bytes"
	"fmt"
	"net/http"
	//"net/http"
	"strings"
)

type PathNode struct {
	// Given the next segment s, if edges[s] exists, then we'll look there first.
	edges map[string]*PathNode

	// If set, failure to match on edges will match on wildcard
	wildcard *PathNode

	// If set, and we have nothing left to match, then we match on this node
	leaf *PathLeaf

	// the parent node
	//parent *PathNode
}

// For the route /admin/forums/:forum_id:\d.*/suggestions/:suggestion_id:\d.*
// We'd have wildcards = ["forum_id", "suggestion_id"]
// For the route /admin/forums/:forum_id/suggestions/:suggestion_id:\d.*
// We'd have wildcards = ["forum_id", "suggestion_id"]
// For the route /admin/forums/:forum_id/suggestions/:suggestion_id
// We'd have wildcards = ["forum_id", "suggestion_id"]
type PathLeaf struct {
	// names of wildcards that lead to this leaf. eg, ["category_id"] for the wildcard ":category_id"
	wildcards []string

	// Pointer back to the route
	*route
}

func newPathNode() *PathNode {
	return &PathNode{edges: make(map[string]*PathNode)}
}

func (pn *PathNode) Inspect(indent int) string {
	var buf bytes.Buffer
	for p, edg := range pn.edges {
		fmt.Fprintf(&buf, "%s/%s\n%s\n", strings.Repeat("\t", indent), p, edg.Inspect(indent+1))
	}
	if pn.wildcard != nil {
		fmt.Fprintf(&buf, "%s*\n%s", strings.Repeat("\t", indent), pn.wildcard.Inspect(indent+1))
	}
	if pn.leaf != nil && pn.leaf.route != nil {
		fmt.Fprintf(&buf, "%s\n%s", strings.Repeat("\t", indent), pn.leaf.route.Inspect(indent))
	}

	return buf.String()
}

func (pn *PathNode) addX(path string, v verb, handler http.Handler) error {
	return pn.addInternalX(splitPath(path), v, handler, nil)
}

func (pn *PathNode) addInternalX(segments []string, v verb, handler http.Handler, wildcards []string) error {
	if len(segments) == 0 {
		if pn.leaf == nil {
			pn.leaf = &PathLeaf{route: NewRoute(), wildcards: wildcards}
		}
		return pn.leaf.route.AddHandlerX(handler, v)

	}
	// len(segments) >= 1
	seg := segments[0]
	wc, wcName := isWildcard(seg)
	if wc {
		if pn.wildcard == nil {
			pn.wildcard = newPathNode()
			//pn.wildcard.parent = pn
		}
		return pn.wildcard.addInternalX(segments[1:], v, handler, append(wildcards, wcName))
	}
	subPn, ok := pn.edges[seg]
	if !ok {
		subPn = newPathNode()
		//subPn.parent = pn
		pn.edges[seg] = subPn
	}
	return subPn.addInternalX(segments[1:], v, handler, wildcards)

}

func (pn *PathNode) Match(path string) (leaf *PathLeaf, wildcards map[string]string) {
	// Bail on invalid paths.
	if len(path) == 0 || path[0] != '/' {
		return nil, nil
	}

	return pn.match(splitPath(path), nil)
}

// Segments is like ["admin", "users"] representing "/admin/users"
// wildcardValues are the actual values accumulated when we match on a wildcard.
func (pn *PathNode) match(segments []string, wildcardValues []string) (leaf *PathLeaf, wildcardMap map[string]string) {
	// Handle leaf nodes:
	if len(segments) == 0 {
		return pn.leaf, makeWildcardMap(pn.leaf, wildcardValues)
	}

	var seg string
	seg, segments = segments[0], segments[1:]

	subPn, ok := pn.edges[seg]
	if ok {
		leaf, wildcardMap = subPn.match(segments, wildcardValues)
	}

	if leaf == nil && pn.wildcard != nil {
		leaf, wildcardMap = pn.wildcard.match(segments, append(wildcardValues, seg))
	}

	return leaf, wildcardMap
}

// key is a non-empty path segment like "admin" or ":category_id" or ":category_id:\d+"
// Returns true if it's a wildcard, and if it is, also returns it's name / regexp.
// Eg, (true, "category_id", "\d+")
func isWildcard(key string) (bool, string) {
	if key[0] == ':' {
		substrs := strings.SplitN(key[1:], ":", 2)
		if len(substrs) == 1 {
			return true, substrs[0]
		} else {
			return true, substrs[0]
		}
	} else {
		return false, ""
	}
}

// "/" -> []
// "/admin" -> ["admin"]
// "/admin/" -> ["admin"]
// "/admin/users" -> ["admin", "users"]
func splitPath(key string) []string {
	elements := strings.Split(key, "/")
	if elements[0] == "" {
		elements = elements[1:]
	}
	if elements[len(elements)-1] == "" {
		elements = elements[:len(elements)-1]
	}
	return elements
}

func makeWildcardMap(leaf *PathLeaf, wildcards []string) map[string]string {
	if leaf == nil {
		return nil
	}

	leafWildcards := leaf.wildcards

	if len(wildcards) == 0 || (len(leafWildcards) != len(wildcards)) {
		return nil
	}

	// At this point, we know that wildcards and leaf.wildcards match in length.
	assoc := make(map[string]string)
	for i, w := range wildcards {
		assoc[leafWildcards[i]] = w
	}

	return assoc
}
