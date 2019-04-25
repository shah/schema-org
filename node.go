package schemamd

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

type NodeMap map[string]*Node
type IDRefs []string

func (refs IDRefs) Contains(id string) (bool, int) {
	for i, s := range refs {
		if id == s {
			return true, i
		}
	}
	return false, -1
}

type Node struct {
	index    int
	defn     FlexMap
	dataType bool
	class    bool
	property bool

	id       string
	name     string
	label    string
	comments string

	domainRefs        IDRefs
	rangeRefs         IDRefs
	partOfRefs        IDRefs
	subPropertyOfRefs IDRefs
}

func NewDataTypeNode(md *MetaData, id string, index int, defn FlexMap, typeRefs []interface{}) *Node {
	node := &Node{index: index, defn: defn, dataType: true, class: true, property: false, id: id, name: md.nameFromID(id)}
	var isClass, isDataType bool
	for _, tr := range typeRefs {
		switch tr {
		case "rdfs:Class":
			isClass = true
		case "http://schema.org/DataType":
			isDataType = true
		}
	}
	if isClass && isDataType {
		node.init(md)
	} else {
		md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("%s (%d)", id, index), GraphNodeIDNotHandled, "Node @type was a slice but not a DataType?", false))
	}
	return node
}

func NewClassNode(md *MetaData, id string, index int, defn FlexMap) *Node {
	node := &Node{index: index, defn: defn, dataType: false, class: true, property: false, id: id, name: md.nameFromID(id)}
	node.init(md)
	return node
}

func NewPropertyNode(md *MetaData, id string, index int, defn FlexMap) *Node {
	node := &Node{index: index, defn: defn, dataType: false, class: false, property: true, id: id, name: md.nameFromID(id)}
	node.init(md)
	return node
}

func (node *Node) init(md *MetaData) {
	for key, value := range node.defn {
		switch key {
		case "@type":
			continue
		case "@id":
			continue
		case "http://schema.org/domainIncludes":
			node.domainRefs = md.flattenIDRefs(node.id, node.index, value)
		case "http://schema.org/rangeIncludes":
			node.rangeRefs = md.flattenIDRefs(node.id, node.index, value)
		case "http://schema.org/isPartOf":
			node.partOfRefs = md.flattenIDRefs(node.id, node.index, value)
		case "rdfs:comment":
			node.comments = value.(string)
			//fmt.Printf("_%T_ %+v\n", value, value)
		case "rdfs:label":
			switch label := value.(type) {
			case string:
				node.label = label
			case map[string]interface{}:
				node.label = label["@value"].(string)
			}
		case "rdfs:subPropertyOf":
			node.subPropertyOfRefs = md.flattenIDRefs(node.id, node.index, value)
		case "http://purl.org/dc/terms/source":
		case "http://schema.org/inverseOf":
		case "http://schema.org/supersededBy":
		case "http://www.w3.org/2002/07/owl#equivalentProperty":
		case "http://www.w3.org/2004/02/skos/core#closeMatch":
		case "http://www.w3.org/2004/02/skos/core#exactMatch":
		case "http://www.w3.org/2002/07/owl#equivalentClass":
		case "rdfs:subClassOf":
		case "http://schema.org/sameAs":
		case "http://schema.org/category":
			// TODO: figureo out what to do with these, ignoring for now
		default:
			md.issuesFound = append(md.issuesFound, newIssue(fmt.Sprintf("%s (%d)", node.id, node.index), GraphPropertyNodeKeyNotHandled, fmt.Sprintf("Property node key %q unknown", key), false))
		}
	}
}

func (node Node) ClassName() string {
	return node.name
}

func (node Node) ClassMembers(md *MetaData) NodeMap {
	result := make(NodeMap)
	for _, pnode := range md.properties {
		found, _ := pnode.domainRefs.Contains(node.id)
		if found {
			result[pnode.id] = pnode
		}
	}
	return result
}

func (node Node) PropertyName() string {
	return strcase.ToLowerCamel(node.name)
}
