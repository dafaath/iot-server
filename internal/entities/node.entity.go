package entities

type Node struct {
	IdNode           int      `json:"id_node"`
	IdUser           int      `json:"id_user"`
	IdHardwareNode   int      `json:"id_hardware_node"`
	IdHardwareSensor []int    `json:"id_hardware_sensor"`
	Name             string   `json:"name"`
	Location         string   `json:"location"`
	FieldSensor      []string `json:"field_sensor"`
	IsPublic         bool     `json:"is_public"`
}

type NodeCreate struct {
	Name             string   `json:"name" validate:"required"`
	Location         string   `json:"location" validate:"required"`
	IdHardwareNode   int      `json:"id_hardware_node" validate:"required"`
	IdHardwareSensor []int    `json:"id_hardware_sensor" validate:"required"`
	FieldSensor      []string `json:"field_sensor" validate:"required"`
	IsPublic         bool     `json:"is_public"`
}

type NodeUpdate struct {
	Name             string   `json:"name"`
	Location         string   `json:"location"`
	IdHardwareNode   int      `json:"id_hardware_node"`
	IdHardwareSensor []int    `json:"id_hardware_sensor"`
	FieldSensor      []string `json:"field_sensor"`
	IsPublic         int      `json:"is_public"`
}

func (hu *NodeUpdate) ChangeSettedFieldOnly(node *Node) {
	if hu.Name == "" {
		hu.Name = node.Name
	}

	if hu.Location == "" {
		hu.Location = node.Location
	}

	if hu.IdHardwareNode == 0 {
		hu.IdHardwareNode = node.IdHardwareNode
	}

	if hu.IdHardwareSensor == nil || len(hu.IdHardwareSensor) == 0 {
		hu.IdHardwareSensor = node.IdHardwareSensor
	}

	if hu.FieldSensor == nil || len(hu.FieldSensor) == 0 {
		hu.FieldSensor = node.FieldSensor
	}

	if hu.IsPublic == 0 {
		hu.IsPublic = convertBoolToInt(node.IsPublic)
	}
}

type NodeWithChannel struct {
	Node
	Feed []Channel `json:"feed"`
}
