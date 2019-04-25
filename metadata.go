package schemamd

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TextMap map[string]string
type FlexMap map[string]interface{}
type Graph []FlexMap

type MetaData struct {
	Context TextMap `json:"@context"`
	Graph   Graph   `json:"@graph"`

	apiEndpoint      string
	issuesFound      []Issue
	dataTypes        NodeMap
	classes          NodeMap
	properties       NodeMap
	handledGraphNode []bool
}

func NewMetaData() *MetaData {
	result := new(MetaData)
	result.apiEndpoint = "https://schema.org/version/3.5/all-layers.jsonld"
	httpRes, issue := getHTTPResult(result.apiEndpoint, HTTPUserAgent, HTTPTimeout)
	if issue != nil {
		result.issuesFound = append(result.issuesFound, issue)
		return result
	}
	result.apiEndpoint = httpRes.apiEndpoint

	json.Unmarshal(*httpRes.body, result)

	result.handledGraphNode = make([]bool, len(result.Graph))
	result.dataTypes = make(NodeMap)
	result.classes = make(NodeMap)
	result.properties = make(NodeMap)

	result.index()

	for i, handled := range result.handledGraphNode {
		if !handled {
			result.issuesFound = append(result.issuesFound, newIssue(fmt.Sprintf("node %d %+v", i, result.Graph[i]), GraphNodeIDNotHandled, "Node was not handled by the indexer", true))
		}
	}

	return result
}

func (md MetaData) nameFromID(id string) string {
	return strings.Title(strings.Split(id, "/")[len(strings.Split(id, "/"))-1])
}

func (md *MetaData) graphNode(index int, node FlexMap) (id string, typ interface{}, handled, ok bool) {
	handled = md.handledGraphNode[index]
	if handled {
		return
	}

	id, ok = node["@id"].(string)
	if !ok {
		md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("node %d %+v", index, node), GraphNodeIDNotFound, "Node @id was not found", true))
		return
	}

	typ, ok = node["@type"]
	if !ok {
		md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("%s (%d) %+v", id, index, node), GraphNodeTypeNotFound, "Node @type was not found", true))
		return
	}

	return
}

func (md *MetaData) flattenIDRefs(id string, index int, node interface{}) IDRefs {
	var idRefs IDRefs
	switch v := node.(type) {
	case map[string]interface{}:
		idRefs = append(idRefs, v["@id"].(string))
	case []interface{}:
		for _, m := range v {
			item := m.(map[string]interface{})
			idRefs = append(idRefs, item["@id"].(string))
		}
	default:
		md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("flattenIDRefs %s (%d) %T", id, index, v), GraphPropertyNodeKeyNotHandled, fmt.Sprintf("Unable to flatten IDs in %+v", v), false))
	}
	return idRefs
}

func (md *MetaData) index() {
	for i, defn := range md.Graph {
		id, typ, handled, ok := md.graphNode(i, defn)
		if !ok || handled {
			md.handledGraphNode[i] = true
			continue
		}

		switch v := typ.(type) {
		case string:
			switch v {
			case "rdfs:Class":
				md.classes[id] = NewClassNode(md, id, i, defn)
			case "rdf:Property":
				md.properties[id] = NewPropertyNode(md, id, i, defn)
			default:
				md.classes[id] = NewClassNode(md, id, i, defn)
			}
		case []interface{}:
			node := NewDataTypeNode(md, id, i, defn, v)
			md.classes[id] = node
			md.dataTypes[id] = node
		default:
			md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("%s (%d) %T", id, i, v), GraphNodeTypeNotFound, "Node @type was not properly handled", true))
			continue
		}

		md.handledGraphNode[i] = true
	}
}
