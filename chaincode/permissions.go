package main

// Action is the the type of an action
type Action string

// Enum the different action types available
const (
	Process  Action = "process"
	Download Action = "download"
)

// Privilege represents one permission based on an action type
type Privilege struct {
	// IsRestricted is true if this permission is given to the asset's owner only and
	// the nodes listed in AuthorizedNodes (open to all nodes if false)
	IsRestricted bool `json:"isRestricted"`
	// AuthorizedNodes list all authorised nodes other than the asset's owner
	AuthorizedNodes []string `json:"authorizedNodes"`
}

// Permissions represents all privileges associated with an asset
type Permissions struct {
	// Download define if a given node can allow its nodes to download the asset
	Download Privilege `json:"download"`
	// Process define if a given node can process the sset
	Process Privilege `json:"process"`
}

// CanProcess checks if a node can process the asset with the current permissions
func (perms Permissions) CanProcess(owner, node string) bool {
	if owner == node {
		return true
	}

	if !perms.Process.IsRestricted {
		return true
	}

	for _, authorizedNode := range perms.Process.AuthorizedNodes {
		if node == authorizedNode {
			return true
		}
	}
	return false
}

// NewPermissions the Permissions Privilege according to the arg received
func NewPermissions(in inputPermissions) Permissions {
	perms := Permissions{}
	process := newPrivilege(in.Process)
	perms.Process = process
	// Download privilege is not implemented in the node server, so it is set to the process privilege
	perms.Download = process
	return perms
}

func newPrivilege(in inputPrivilege) Privilege {
	return Privilege{
		IsRestricted:    in.IsRestricted,
		AuthorizedNodes: in.AuthorizedNodes,
	}
}

func (priv Privilege) include(other Privilege) bool {
	if !priv.IsRestricted {
		return true
	}
	if !other.IsRestricted {
		return false
	}
	for _, node := range other.AuthorizedNodes {
		if !stringInSlice(node, priv.AuthorizedNodes) {
			return false
		}
	}
	return true
}

// MergePermissions returns the intersection of input permissions
func MergePermissions(x, y Permissions) Permissions {
	perm := Permissions{}
	perm.Process = mergePrivileges(x.Process, y.Process)
	perm.Download = mergePrivileges(x.Download, y.Download)
	return perm
}

func mergePrivileges(x, y Privilege) Privilege {
	priv := Privilege{}
	priv.IsRestricted = x.IsRestricted || y.IsRestricted

	if x.IsRestricted && !y.IsRestricted {
		priv.AuthorizedNodes = x.AuthorizedNodes
	} else if !x.IsRestricted && y.IsRestricted {
		priv.AuthorizedNodes = y.AuthorizedNodes
	} else {
		priv.AuthorizedNodes = x.getNodesIntersection(y)
	}
	return priv
}

func (priv Privilege) getNodesIntersection(p Privilege) []string {
	nodes := []string{}
	for _, i := range priv.AuthorizedNodes {
		for _, j := range p.AuthorizedNodes {
			if i == j {
				nodes = append(nodes, i)
				break
			}
		}
	}
	return nodes
}
