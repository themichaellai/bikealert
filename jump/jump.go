package jump

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client has methods for accessing JUMP data.
type Client struct {
	networkID string

	httpClient *http.Client
}

const httpTimeout = 5 * time.Second

// NewClient creates a new JUMP client. It will make requests with
// respect to the given JUMP network ID.
func NewClient(networkID string) *Client {
	return &Client{
		networkID: networkID,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// Position contains coordinates for a bike or hub.
type Position struct {
	// Coordinates is a two-element list.
	Coordinates []float64 `json:"coordinates"`
}

// Bike has information about a bike and its location.
type Bike struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	NetworkID            int64    `json:"network_id"`
	StatsLastPor         string   `json:"stats_last_por"`
	BatteryLevel         int64    `json:"battery_level"`
	VehicleType          string   `json:"vehicle_type"`
	UnlockingMethods     []string `json:"unlocking_methods"`
	Sponsored            bool     `json:"sponsored"`
	EbikeBatteryLevel    int64    `json:"ebike_battery_level"`
	EbikeBatteryDistance float64  `json:"ebike_battery_distance"`
	InsideArea           bool     `json:"inside_area"`
	Address              string   `json:"address"`
	CurrentPosition      Position `json:"current_position"`
}

// bikesResponse is the response from /bikes.
type bikesResponse struct {
	CurrentPage  int64  `json:"current_position"`
	PerPage      int64  `json:"per_page"`
	TotalEntries int64  `json:"total_entries"`
	Items        []Bike `json:"items"`
}

// Bikes retrieves all of the bikes for the network.
func (c *Client) Bikes() ([]Bike, error) {
	errPrefix := "jump.Bikes"

	url := fmt.Sprintf(
		"https://app.jumpbikes.com/api/networks/%s/bikes?collapsed=false&per_page=999",
		c.networkID)
	req, err := c.newRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, errPrefix)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errPrefix)
	} else if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		var body string
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			body = fmt.Sprintf("could not parse body (%s)", err.Error())
		} else {
			body = string(bodyBytes)
		}
		return nil, errors.Wrap(
			fmt.Errorf("got status code %d: %s", res.StatusCode, body),
			errPrefix)
	}

	var parsedBody bikesResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&parsedBody); err != nil {
		return nil, errors.Wrap(err, errPrefix)
	}
	return parsedBody.Items, nil
}

// Hub has information about a hub and its location.
type Hub struct {
	ID          float64 `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`

	Address                   string           `json:"address"`
	AreaID                    int64            `json:"area_id"`
	AvailableBikes            int64            `json:"available_bikes"`
	AvailableEbikes           int64            `json:"available_ebikes"`
	AvailableScooters         int64            `json:"available_scooters"`
	AvailableVehicles         int64            `json:"available_vehicles"`
	CollapseBikes             bool             `json:"collapse_bikes"`
	CurrentBikes              int64            `json:"current_bikes"`
	DisplayMethod             string           `json:"display_method"`
	FreeRacks                 int64            `json:"free_racks"`
	HasChargingInfrastructure bool             `json:"has_charging_infrastructure"`
	HasKiosk                  bool             `json:"has_kiosk"`
	LowChargeBountyHub        bool             `json:"low_charge_bounty_hub"`
	MiddlePoint               Position         `json:"middle_point"`
	NetworkID                 float64          `json:"network_id"`
	Polygon                   *json.RawMessage `json:"polygon"`
	Priority                  bool             `json:"priority"`
	Public                    bool             `json:"public"`
	RacksAmount               float64          `json:"racks_amount"`
	RebalanceBountyHub        bool             `json:"rebalance_bounty_hub"`
	Sponsored                 bool             `json:"sponsored"`
	SponsoredBikes            int64            `json:"sponsored_bikes"`
	Visible                   bool             `json:"visible"`
	Warehouse                 bool             `json:"warehouse"`
}

type hubResponse struct {
	CurrentPage  int64 `json:"current_position"`
	PerPage      int64 `json:"per_page"`
	TotalEntries int64 `json:"total_entries"`
	Items        []Hub `json:"items"`
}

// Hubs retrieves all of the hubs for the network.
func (c *Client) Hubs() ([]Hub, error) {
	errPrefix := "jump.Hubs"

	url := fmt.Sprintf(
		"https://app.jumpbikes.com/api/networks/%s/hubs?collapsed=false&per_page=999",
		c.networkID)
	req, err := c.newRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, errPrefix)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errPrefix)
	} else if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		var body string
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			body = fmt.Sprintf("could not parse body (%s)", err.Error())
		} else {
			body = string(bodyBytes)
		}
		return nil, errors.Wrap(
			fmt.Errorf("got status code %d: %s", res.StatusCode, body),
			errPrefix)
	}

	var parsedBody hubResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&parsedBody); err != nil {
		return nil, errors.Wrap(err, errPrefix)
	}
	return parsedBody.Items, nil
}

func (c *Client) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Sorry.
	req.Header.Add("Referer", fmt.Sprintf("https://map.jump.com/?network_id=%s&theme=jump", c.networkID))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	return req, nil
}
