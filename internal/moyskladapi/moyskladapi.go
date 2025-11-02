package moyskladapi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"mstorefgo/internal/config"
)

type MoySkladProcessor struct {
	Ratelimiter *Ratelimiter
	Config      *config.Moyskladapiconfig
}

func NewMoySkladProcessor(r *Ratelimiter, c *config.Moyskladapiconfig) *MoySkladProcessor {
	return &MoySkladProcessor{Ratelimiter: r, Config: c}
}

func (m *MoySkladProcessor) FetchDeliverableOrders() []byte {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	dayaftertomorrow := tomorrow.AddDate(0, 0, 1)
	tomorrowstart := ">=" + tomorrow.Format("2006-01-02") + " 00:00:00"
	fmt.Println(tomorrowstart)

	dayaftertomorrowend := "<=" + dayaftertomorrow.Format("2006-01-02") + " 23:59:59"

	baseURL, err := url.Parse(m.Config.URLstart)
	if err != nil {
		panic(err)
	}
	baseURL.Path = path.Join(baseURL.Path, "entity/customerorder")

	filterValue := fmt.Sprintf(
		"deliveryPlannedMoment%s;deliveryPlannedMoment%s;state=%s", tomorrowstart, dayaftertomorrowend, m.Config.Readystatehref)

	q := baseURL.Query()
	q.Set("filter", filterValue)
	baseURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+m.Config.APIKEY)
	req.Header.Set("Accept-Encoding", "gzip")
	log.Println("Waiting for RateLimiter")
	m.Ratelimiter.Wait()
	log.Println("Done waiting")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	return body
}

func (m *MoySkladProcessor) FetchEntityByHREF(href string) []byte {
	url := href
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+m.Config.APIKEY)
	req.Header.Set("Accept-Encoding", "gzip")

	m.Ratelimiter.Wait()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	return body
}

type orderShipmentUpdate struct {
	State *State `json:"state"`
}

type State struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Href      string `json:"href"`
	Type      string `json:"type"`
	MediaType string `json:"metadataHref"`
}

func (m *MoySkladProcessor) SetOrderShipped(href string) error {
	url := href
	orderShipmentUpdate := orderShipmentUpdate{
		State: &State{
			Meta: Meta{
				Href:      m.Config.Shipedstatehref,
				Type:      "state",
				MediaType: "application/json",
			},
		},
	}

	orderShipmentUpdateJSON, err := json.Marshal(orderShipmentUpdate)
	if err != nil {
		fmt.Printf("failed to marshal orderShipmentUpdate: %s", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(orderShipmentUpdateJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.Config.APIKEY)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	m.Ratelimiter.Wait()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil
}

type Organization struct {
	Meta *Meta `json:"meta"`
}

type Agent struct {
	Meta Meta `json:"meta"`
}

type Store struct {
	Meta Meta `json:"meta"`
}

type CustomerOrder struct {
	Meta Meta `json:"meta"`
}

type Positions struct {
	Meta Meta `json:"meta"`
}
type ShipmentInfo struct {
	Organization  Organization  `json:"organization"`
	Agent         Agent         `json:"agent"`
	Store         Store         `json:"store"`
	CustomerOrder CustomerOrder `json:"customerOrder"`
	Positions     Positions     `json:"positions"`
}

func (m *MoySkladProcessor) ShipOrder(orderhref string, counterhref string) error {
	url := "https://api.moysklad.ru/api/remap/1.2/entity/demand"
	positionshref := strings.Trim(orderhref, `"`)
	positionshref = `"` + positionshref + "/positions" + `"`
	fmt.Println(positionshref)
	shipmentinfo := ShipmentInfo{
		Organization: Organization{
			Meta: &Meta{
				Href:      m.Config.Orghref,
				Type:      "organization",
				MediaType: "application/json",
			},
		},
		Agent: Agent{
			Meta: Meta{
				Href:      counterhref,
				Type:      "counterparty",
				MediaType: "application/json",
			},
		},
		Store: Store{
			Meta: Meta{
				Href:      m.Config.Storehref,
				Type:      "store",
				MediaType: "application/json",
			},
		},
		CustomerOrder: CustomerOrder{
			Meta: Meta{
				Href:      orderhref,
				Type:      "customerorder",
				MediaType: "application/json",
			},
		},
		Positions: Positions{
			Meta: Meta{
				Href:      positionshref,
				Type:      "demandposition",
				MediaType: "application/json",
			},
		},
	}

	shipmentinfoJSON, err := json.Marshal(shipmentinfo)
	if err != nil {
		fmt.Printf("failed to marshal orderShipmentUpdate: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(shipmentinfoJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.Config.APIKEY)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	m.Ratelimiter.Wait()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil

}
