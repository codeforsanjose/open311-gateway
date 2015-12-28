## To Do.todo

* Implement Gingko BDD tests. @done(2015-12-18)
* Modify report.Create() to use ServiceID and JID.
* Update RAML file with JSON specs for input and output payloads.
* Implement report searches:
	* Single ID
	* DeviceID
	* LatLng
	* Address
* Implement Services API endpoint. @done(2015-12-18)
	* Bring Google Maps API over. @done(2015-12-18)
* Modify CitySourced simulator to return the Request ID and Document ID. @done(2015-12-12)
* Outline the Displatch system. @done(2015-12-14)

## Log

[2015.12.28 - Mon]

* Greatly reorganized request processing to Create a report (in the "request" package):
	* Renamed:
		* CreateReport type -> CreateReq
		* CreateReportResp type -> CreateResp
	* Moved all Create functionality under the CreateReq type - current methods:
		1. validate
		2. init
		3. parseQP
		4. run
		5. ProcessCS (processes against CitySourced backend)
		6. String
		7. toCS (converts CreateReq struct to citysourced.CSReport type)
	* Moved all of the above Create types into the new file: "request/create.go".
	* Added "apiVersion" to config.json in the Provider section, and also added to the Provider type in "router/data.go".
	* Saved to GIT.
	* Test OK
	* Saved to GIT.

[2015.12.23 - Wed]

* Cleaned up router.ServiceProviderInterface().  Added test cases in data_test.go.
* Test OK.
* 

[2015.12.19 - Sat]

* In router/data.go:
	* Dropped "Getxxx" from the ServiceXXX() methods.
	* Added ServiceProviderInterface() go get the Service Provider interface type (currently on CitySourced).  
* in request/report.go, added a map lookup (beCreate) on the Create functions.  We will have a set of these maps to quickly route an incoming request (create, lookup, etc), to the correct backend interface.
* Test OK.
* Saved to GIT.

[2015.12.18 - Fri]

* Brought "geo" package over from CitySourced.
* Added "getCity()" function to mygeocode.go.  This scans through the Google response and retrieves the city.  We will need this for quickly mapping the Mobile Apps geoloc -> city -> Service Providers -> list of Services.
* Test OK.
* Saved to GIT.
* Added a GetCity() func in the "geo" package.  This takes a latitude and longitude, and returns the City.
* Test OK.
* Saved to GIT.
* "/services" endpoint working - returning the list of services for San Jose.  Returns 500 for city: "Morgan Hill", with error: The city "Morgan Hill" is not serviced by this Gateway.
* Test OK.
* Saved to GIT.
* Added Gingko test suite.
* Revised router/data.go:
	* Moved Services() and ServiceProviders() to be methods of RouteData.
	* ServiceProviders() now returns the Provider list, and an error.
* Test OK.
* Saved to GIT.


[2015.12.17 - Thu]

* Router data and indices are working.
* Test OK.
* Saved to GIT.
* Split the test data into two Service Providers for San Jose to test multiple providers.  
* Added providerService to the RouteData struct.  This implements a map of the ServiceID back to the Service Provider.  This provides a fast lookup of the appropriate provider for a New Report, based on the ServiceID.
* Test OK.
* Saved to GIT.

[2015.12.16 - Wed]

* Completely reformatted JSON file... made Services a direct child of Service Provider, and Service Provider is a child element of Service Areas.  Go will load this as a series of maps, which will automatically give us some fast indexing into the data.

[2015.12.15 - Tue]

* Updated ServiceRouter.md documentation.
* Created 

[2015.12.14 - Mon]

* Added design documentation for the Service Router capabilities.  
* Saved to GIT.

[2015.12.10 - Thu]

* Doing other things for a few days.  Coursework, etc.  Back on the case now!
* Current status: "Create" is working for CitySourced API. 
	* Current libraries
		* github.com/ant0ine/go-json-rest/rest - this is working well for the front-end, but does not support XML, and so will not work for the backend.
	* No routing/dispatching is in place yet...  This is probably the next task.

[2015.12.04 - Fri]

* Implemented "napping" REST client.  Not working due to XML.  Starting over using just the HTTP lib.
* Saved to GIT.
* 

[2015.12.01 - Tue]

* 

[2015.11.30 - Mon]

* Wrote up 311 Gateway Proposal.

[2015.10.27 - Fri]

* Reviewed RAML spec with Hassan.  Discussed the Gateway idea.

[2015.10.24 - Tue]

* Revised RAML spec.

[2015.10.23 - Mon]

* First draft of RAML spec.
* Posted to Slack.
