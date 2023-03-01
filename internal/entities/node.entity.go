package entities

type Node struct {
	IdNode int `json:"id_node" validate:"required"`
	NodeCreate
	IdUser int `json:"id_user" validate:"required"`
}

type NodeCreate struct {
	Name       string `json:"name" validate:"required"`
	Location   string `json:"location" validate:"required"`
	IdHardware int    `json:"id_hardware" validate:"required"`
}

type NodeUpdate struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (hu *NodeUpdate) ChangeSettedFieldOnly(node *Node) {
	if hu.Name == "" {
		hu.Name = node.Name
	}

	if hu.Location == "" {
		hu.Location = node.Location
	}
}

type NodeWithHardwareAndSensors struct {
	Node
	Hardware Hardware `json:"hardware"`
	Sensor   []Sensor `json:"sensor"`
}
