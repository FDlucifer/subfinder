// 
// certspotter.go : A Golang based client for Certspotter Parsing
// Written By : @ice3man (Nizamul Rana)
// 
// Distributed Under MIT License
// Copyrights (C) 2018 Ice3man
//

package certspotter

import (
	"io/ioutil"
	"encoding/json"
	"strings"
	"fmt"

	"subfinder/libsubfinder/helper"
)

// Structure of a single dictionary of output by crt.sh
type certspotter_object struct {
	Dns_names	[]string `json:"dns_names"`
}

// array of all results returned
var certspotter_data []certspotter_object

// all subdomains found
var subdomains []string

// 
// Query : Queries awesome Certspotter service for subdomains
// @param state : current application state, holds all information found
//
func Query(state *helper.State, ch chan helper.Result) {

	// Create a result object 
	var result helper.Result
	result.Subdomains = subdomains

	// Make a http request to Certspotter
	resp, err := helper.GetHTTPResponse("https://certspotter.com/api/v0/certs?domain="+state.Domain, 3000)
	if err != nil {
		// Set values and return
		result.Error = err
		ch <- result
	}

	// Get the response body
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.Error = err
		ch <- result
	}

	// Decode the json format
	err = json.Unmarshal([]byte(resp_body), &certspotter_data)
	if err != nil {
		result.Error = err
		ch <- result
	}

	// Append each subdomain found to subdomains array
	for _, block := range certspotter_data {
		for _, dns_name := range block.Dns_names {

			// Fix Wildcard subdomains containg asterisk before them
			if strings.Contains(dns_name, "*.") {
				dns_name = strings.Split(dns_name, "*.")[1]
			}

			if state.Verbose == true {
				if state.Color == true {
					fmt.Printf("\n[%sCERTSPOTTER%s] %s", helper.Red, helper.Reset, dns_name)
				} else {
					fmt.Printf("\n[CERTSPOTTER] %s", dns_name)
				}
			}

			subdomains = append(subdomains, dns_name)
		}	
	}

	result.Subdomains = subdomains
	result.Error = nil
	ch <-result
}
