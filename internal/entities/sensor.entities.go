package entities

type Sensor struct {
	IdSensor int `json:"id_sensor" validate:"required"`
	SensorCreate
}

type SensorCreate struct {
	Name       string `json:"name" validate:"required"`
	Unit       string `json:"unit" validate:"required"`
	IdNode     int    `json:"id_node" validate:"required"`
	IdHardware int    `json:"id_hardware" validate:"required"`
}

type SensorUpdate struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

func (su *SensorUpdate) ChangeSettedFieldOnly(sensor *Sensor) {
	if su.Name == "" {
		su.Name = sensor.Name
	}

	if su.Unit == "" {
		su.Unit = sensor.Unit
	}
}

type SensorWithChannel struct {
	Sensor
	Channel []Channel `json:"channel"`
}
