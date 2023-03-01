package entities

type HardwareCreate struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required,oneof='microcontroller unit' 'single-board computer' 'sensor'"`
	Description string `json:"description" validate:"required"`
}

type HardwareUpdate struct {
	Name        string `json:"name"`
	Type        string `json:"type" validate:"oneof='microcontroller unit' 'single-board computer' 'sensor'"`
	Description string `json:"description"`
}

func (hu *HardwareUpdate) ChangeSettedFieldOnly(hardware *Hardware) {
	if hu.Name == "" {
		hu.Name = hardware.Name
	}

	if hu.Type == "" {
		hu.Type = hardware.Type
	}

	if hu.Description == "" {
		hu.Description = hardware.Description
	}
}

type Hardware struct {
	IdHardware int `json:"id_hardware" validate:"required"`
	HardwareCreate
}

type HardwareWithNode struct {
	Hardware
	Nodes []Node `json:"nodes"`
}

type HardwareWithSensor struct {
	Hardware
}
