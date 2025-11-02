package unmarshaller

import (
	"encoding/json"
	"fmt"
)

type Meta struct {
	Size int    `json:"size"`
	HREF string `json:"href"`
	Type string `json:"type"`
}

type Agent struct {
	Meta Meta `json:"meta"`
}

type Attributes struct {
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type Order struct {
	HREF                  string                 `json:"href"`
	Meta                  Meta                   `json:"meta"`
	Name                  string                 `json:"name"`
	Sum                   float64                `json:"sum"`
	Agent                 Agent                  `json:"agent"`
	Attributes            []Attributes           `json:"attributes"`
	Description           string                 `json:"description"`
	Positions             Positions              `json:"positions"`
	DeliveryPlannedMoment string                 `json:"deliveryPlannedMoment"`
	AttributesMap         map[string]interface{} `json:"-"`
	ShipmentAddress       string                 `json:"shipmentAddress"`
}

func (o *Order) UnmarshallOrderAttributes() error {
	o.AttributesMap = make(map[string]interface{})
	for _, attribute := range o.Attributes {
		var value interface{}
		var err error

		switch attribute.Type {
		case "string":
			var s string
			err = json.Unmarshal(attribute.Value, &s)
			if err != nil {
				return fmt.Errorf("failed to parse attribute %s: %v", attribute.Name, err)
			}
			value = s
		case "customentity":
			var ce struct {
				Name string `json:"name"`
			}
			err = json.Unmarshal(attribute.Value, &ce)
			value = ce.Name
			if err != nil {
				return fmt.Errorf("failed to parse attribute %s: %v", attribute.Name, err)
			}
		case "employee":
			var Emp struct {
				Name string `json:"name"`
			}
			err = json.Unmarshal(attribute.Value, &Emp)
			value = Emp.Name
			if err != nil {
				return fmt.Errorf("failed to parse attribute %s: %v", attribute.Name, err)
			}

		}
		o.AttributesMap[attribute.Name] = value
	}
	return nil
}

type UnmarshalledMoySkladResponse struct {
	Meta struct {
		Size int `json:"size"`
	} `json:"meta"`
	Rows []Order `json:"rows"`
}

func BasicMoySkladResponseUnmarshalling(body []byte) (*UnmarshalledMoySkladResponse, error) {
	Response := UnmarshalledMoySkladResponse{}

	err := json.Unmarshal(body, &Response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i := range Response.Rows {
		if err := Response.Rows[i].UnmarshallOrderAttributes(); err != nil {
			return nil, fmt.Errorf("failed to parse attributes: %w", err)
		}
	}

	return &Response, nil
}

type UnmarshalledPositionsHREFS struct {
	Rows []Position `json:"rows"`
}

type Position struct {
	Meta struct {
		HREF string `json:"href"`
	} `json:"meta"`
	Quantity float64 `json:"quantity"`
	// Price float64 in copecks
	Price      float64 `json:"price"`
	Assortment struct {
		Meta Meta `json:"meta"`
	} `json:"assortment"`
	PositionType string `json:"type"`
}

type Positions struct {
	Meta Meta `json:"meta"`
}

func PositionsUnmarshalling(body []byte) (*UnmarshalledPositionsHREFS, error) {
	Response := UnmarshalledPositionsHREFS{}
	err := json.Unmarshal(body, &Response)
	if err != nil {
		return nil, err
	}
	return &Response, nil
}

type Reciever struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func AgentUnmarshalling(body []byte) (*Reciever, error) {
	Response := Reciever{}
	err := json.Unmarshal(body, &Response)
	if err != nil {
		return nil, err
	}
	return &Response, nil
}

type ProductInfo struct {
	Code string `json:"code"`
	// Weight is fetched in float64(kilograms) so in future calculations it should be converted to gramms
	Weight float64 `json:"weight"`
}

func ProductInfoUnmarshalling(body []byte) (*ProductInfo, error) {
	Response := ProductInfo{}
	err := json.Unmarshal(body, &Response)
	if err != nil {
		return nil, err
	}
	return &Response, nil
}
