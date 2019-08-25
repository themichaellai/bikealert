package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/themichaellai/bikealert/jump"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	latitude, err := getEnvFloat("LAT")
	if err != nil {
		return err
	}
	longitude, err := getEnvFloat("LNG")
	if err != nil {
		return err
	}

	jumpClient := jump.NewClient(jump.NetworkSanFrancisco)
	bikes, err := jumpClient.Bikes()
	if err != nil {
		return err
	}

	sort.Slice(bikes, func(i, j int) bool {
		iLocation := bikes[i].CurrentPosition.Coordinates
		jLocation := bikes[j].CurrentPosition.Coordinates
		iDistance := distance(latitude, longitude, iLocation[1], iLocation[0])
		jDistance := distance(latitude, longitude, jLocation[1], jLocation[0])
		return iDistance < jDistance
	})

	fmt.Println("Bikes")
	for _, bike := range bikes[:5] {
		location := bike.CurrentPosition.Coordinates
		dist := distance(latitude, longitude, location[1], location[0])
		fmt.Printf("Bike %s %s (%0.2f miles)\n", bike.Name, bike.Address, dist)
	}
	fmt.Println("")

	hubs, err := jumpClient.Hubs()
	if err != nil {
		return err
	}

	sort.Slice(hubs, func(i, j int) bool {
		iLocation := hubs[i].MiddlePoint.Coordinates
		jLocation := hubs[j].MiddlePoint.Coordinates
		iDistance := distance(latitude, longitude, iLocation[1], iLocation[0])
		jDistance := distance(latitude, longitude, jLocation[1], jLocation[0])
		return iDistance < jDistance
	})

	fmt.Println("Hubs")
	for _, hub := range hubs[:5] {
		location := hub.MiddlePoint.Coordinates
		dist := distance(latitude, longitude, location[1], location[0])
		fmt.Printf("Hub %s %s (%0.2f miles)\n", hub.Name, hub.Address, dist)
	}
	return nil
}

func getEnvFloat(name string) (float64, error) {
	val, set := os.LookupEnv(name)
	if !set {
		return 0, fmt.Errorf("envvar \"%s\" not set", name)
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return f, errors.Wrap(err, "error parsing env var \"%s\" as float")
	}
	return f, nil
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// distance returns distance between two coordinates in miles.
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, r float64
	// convert to radians
	// must cast radius as float to multiply later
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	r = 3958.756

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
