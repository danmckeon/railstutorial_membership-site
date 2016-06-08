package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// BingMapResponse represents the internal structure for API responses, omitting data that we aren't using right now
type BingMapResponse struct {
	ResourceSets []struct {
		Resources []struct {
			Point struct {
				Type        string        `json:"type"`
				Coordinates []json.Number `json:"coordinates"`
			} `json:"point"`
		} `json:"resources"`
	} `json:"resourceSets"`
	statusCode uint16
}

// Coordinate stores coordinates as numeric strings
type Coordinate struct {
	Latitude  string
	Longitude string
}

// Location stores info and coordinates for a place
type Location struct {
	Name        string
	Description string
	Geodata     Coordinate
}

// data will carry a pointer to a custom type (struct) that determines how to parse the JSON
func getJSON(url string, data interface{}) error {
	// prepare the new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// make the request and store the response
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	// retrieve the body of the response as a raw byte array
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	// TODO: Test for empty response

	// Parse the JSON response
	json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	// we succeeded
	return nil
}

func geocode(loc string) Coordinate {
	// Storage for API response (make sure that $BING_KEY is set)
	var data BingMapResponse
	url := fmt.Sprintf("http://dev.virtualearth.net/REST/v1/Locations/%s?maxResults=1&key=%s",
		loc, os.Getenv("BING_KEY"))

	// Call the BING API and store JSON
	err := getJSON(url, &data)
	if err != nil {
		panic(err)
	}

	point := data.ResourceSets[0].Resources[0].Point
	location := Coordinate{Latitude: point.Coordinates[0].String(), Longitude: point.Coordinates[1].String()}

	return location
}

// keeping it simple we'll hard code the query
func findNearby() []Location {
	// Store the results
	var locations []Location

	// Commented code would load from a database
	/*
			// For simplicity the search area is hardcoded rather than coming from the user
			neighborhood := "Ballard"


			rows, err := Db.Query("select name, descripition, latitude, longitude from place where neighborhood = $1", neighborhood)
			if err != nil {
		    return
			}

			for rows.Next() {
				// for each row, we create an empty Location object
				var loc Location

				// go can scan the columns returned from the select directly into the properties from our object
				// we need &loc.xxx so that scan can update the properties in memory (&loc.Name means address of the Name property for this instance of loc)
				err = rows.Scan(&loc.Name, &loc.Description, &loc.Geodata.Latitude, &loc.Geodata.Longitude)
				if err != nil {
		    	return
		    }
				// append each intermediate loc to our array
		    locations = append(locations, loc)
		  }
			rows.Close()
	*/

	locations = []Location{
		{
			Name:        "Ballard Coffee Works",
			Description: "Coffee & Pastries",
			Geodata:     geocode("2060 NW Market St, 98107"),
		},
		{
			Name:        "Kiss Cafe",
			Description: "Beer & Community",
			Geodata:     geocode("2817 NW Market St, 98107"),
		},
	}

	return locations
}

// This handler will take care of GET and POST requests
func indexHandler(w http.ResponseWriter, r *http.Request) {
	form, _ := template.ParseFiles("./index.html")
	switch r.Method {
	case "GET": // No form data to retrieve
		form.Execute(w, "nothing")
	// case "POST": // Handle form data and update the page
	// 	r.ParseForm()
	// 	locations := findNearby()
	// 	current := geocode(r.FormValue("location"))

	// 	// This struct holds all of the context that my index template needs to render
	// 	context := struct {
	// 		Mapbox    string
	// 		Current   Coordinate
	// 		Locations []Location
	// 	}{
	// 		os.Getenv("MAPBOX_KEY"),
	// 		current,
	// 		locations,
	// 	}
	// {{ range .Locations }}
	// 	{{ .Name }}
	//  {{ .Description }}
	// {{ end }}

	// UNCOMMENT THE FOLLOWING LINES TO RETURN JSON INSTEAD OF BUILDING HTML
	// *********************************************************************
	// jsonified, err := json.Marshal(context)
	// if err != nil {
	// 	panic(err)
	// }
	// w.Write(jsonified)
	// *********************************************************************

	// pass contextual data to the template engine
	// form.Execute(w, context)
	default: // we'll treat the default case as a GET
		form.Execute(w, "nothing")
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	server := http.Server{
		Addr: ":" + port,
	}
	http.HandleFunc("/", indexHandler)
	server.ListenAndServe()
}
