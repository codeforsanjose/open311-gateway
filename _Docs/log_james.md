## To Do.todo

* Custom Type for ServiceID.
* CitySourced Adapter - load config.json data file. @done(2016-01-05)
* Get Create working again.
* Separate back-end interfaces as separate apps, using go-rpc calls. @done(2016-01-05)
* Change request/common.go: error check validate()... do not error check body or query parm conversion.
* Implement report searches:
	* Single ID
	* DeviceID
	* LatLng
	* Address
* Update RAML file with JSON specs for input and output payloads.
* Implement Gingko BDD tests. @done(2015-12-18)
* Modify report.Create() to use ServiceID and JID. @done(2015-12-28)
* Consolidate all Create functionality within the request.CreateReq type. @done(2015-12-28)
* Implement Services API endpoint. @done(2015-12-18)
	* Bring Google Maps API over. @done(2015-12-18)
* Modify CitySourced simulator to return the Request ID and Document ID. @done(2015-12-12)
* Outline the Displatch system. @done(2015-12-14)

## Log

### 2016.01.11 - Mon

* In integration/citysourced/data/data.go;
	* All indexes working now.
	* Deleted the "serviceProvider" index... it is no longer needed with the new Master Service IDs.
	* Added a Santa Clara section of test data in config.json.
* Moved the structs.go file into it's own package: integration/citysourced/structs/structs.go.  Had a circular import issue.
* Saved to GIT.
* RPC test working OK with Services API.
* Saved to GIT.
	
### 2016.01.07 - Thu

* Tried an alternative layout using the ID's as keys in the Adapter/config.json file.  This (theoretically) would have been convenient, but the JSON decoder doesn't support anything but strings for the "keys".  Not working well.   
* Revised config.json again... switched back to mostly lists.
* integration/citysourced/data/data.go working.  Data being loaded successfully.  Service list is of the NService type, so it will be very efficient to look up and return to the Engine.
* Test OK, for what it is.
* Saved to GIT.

### 2016.01.06 - Wed

* Worked out JSON marshalling for the Service ID custom type.
* Need to rearrange the citysourced/data/config.json file again... need to expand the Area (City) section.
* Saved to GIT.
* Rearranged citysourced/data/config.json.   

### 2016.01.05 - Tue

* Wrote up "EngineAndAdapters.md" file outlining design of the Engine / Adapter system.
* Saved to GIT.
* Changed _Test/CSclient.go to use asynchronous API calls.  Working OK.
* Saved to GIT.

### 2016.01.04 - Mon

* Created separate directories for the Engine and Adapters, and separated the existing source code.
* Saved to GIT.
* Worked on design of Engine/Adapter RPC system - responsibilities, design.
* Got the CitySourced client rudimentarily working.  
* Put all "Native" structs into a the request/structs.go file.
* 

### 2015.12.30 - Wed

* Renamed geo functions:
	* LatLngForAddr() -> LatLngForAddr()
	* AddrForLatLng() -> AddrForLatLng()
	* CityForLatLng() -> cityForLatLng()
	

### 2015.12.29 - Tue

* In integration/citysourced.go, created Search structs:
	* CSSearchLLReq - search by Lat/Lng
	* CSSearchDIDReq - search by Device ID
	* CSSearchZipReq - search by a Zip Code
* Created request/search.go:


---
#### How to Handle Dispatching Searches

**Search by LatLng**

* For now, find the city for the Lat/Lng coordinates, and send a search to all service providers for that city, with the specified coordinates and radius.
* *Limit the radius - 100m?*
* __Recipe__
	* Get City

**Search by Device ID**

* If the request includes a list of previous Service Provider ID's, then use it.
* Use the current location (or specified address), and query all Service Providers for that City for the Device ID.

**Search by Zip**

* Easy with CitySourced...

---

### 2015.12.28 - Mon

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
	* Removed JID from Create URL.  Also removed from use in Create struct, etc.
	* Moved the common code (cType and cIFace) from request/report.go to new file request/common.go.
	* In request/common.go, discontinued use of the "inputBody" and "inputQP" fields.  In cType.init(), will always attempt to decode payload, and parse query parms.  Simplifies the code.
	* Test OK.
	* Saved to GIT.

---
Thoughts on using the JID for most/all requests:

* The JID identifies the City ("Jurisdiction").  _This may become more finely grained using some type of non-overlapping jurisdictional map areas that do not necessarily match city boundaries._
* Requsts that return the Jurisdiction ID:
	* Service List
* Requests that would use the current Jurisdiction ID:
	* Create _(this can use Service ID to get the Jurisdiction ID)_
	* Search by DeviceID _(if the App cannot provide a list of previous Jurisdiction ID's)_
* Requests that would not necessarily use the current Jurisdiction:
	* Search by Current Location _(we need the exact lat/lng)_
	* Search for an Address
---

* Reorganized the Services request (get a list of available services for the specified location):
	* Created request/services.go.
	* Used Create as a template.
	* Created ServicesReq and ServicesResp types.
	* Encapsulated all services retrieval functionality inside the ServicesReq type.
* Modified request/common.go:
	* Changed init() to load().  This will remove the conflict with the parent init().
	* Test for error: "JSON payload is empty".  TODO: This needs to be cleaned up... all checking should fall on the validate() function, and the body and query parm parsing errors should be ignored.
* Renamed request/report.go to request.go, as this is the primary file in the request package.
* Test OK.
* Saved to GIT.
		

### 2015.12.23 - Wed

* Cleaned up router.ServiceProviderInterface().  Added test cases in data_test.go.
* Test OK.
* 

### 2015.12.19 - Sat

* In router/data.go:
	* Dropped "Getxxx" from the ServiceXXX() methods.
	* Added ServiceProviderInterface() go get the Service Provider interface type (currently on CitySourced).  
* in request/report.go, added a map lookup (beCreate) on the Create functions.  We will have a set of these maps to quickly route an incoming request (create, lookup, etc), to the correct backend interface.
* Test OK.
* Saved to GIT.

### 2015.12.18 - Fri

* Brought "geo" package over from CitySourced.
* Added "getCity()" function to mygeocode.go.  This scans through the Google response and retrieves the city.  We will need this for quickly mapping the Mobile Apps geoloc -> city -> Service Providers -> list of Services.
* Test OK.
* Saved to GIT.
* Added a CityForLatLng() func in the "geo" package.  This takes a latitude and longitude, and returns the City.
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


### 2015.12.17 - Thu

* Router data and indices are working.
* Test OK.
* Saved to GIT.
* Split the test data into two Service Providers for San Jose to test multiple providers.  
* Added providerService to the RouteData struct.  This implements a map of the ServiceID back to the Service Provider.  This provides a fast lookup of the appropriate provider for a New Report, based on the ServiceID.
* Test OK.
* Saved to GIT.

### 2015.12.16 - Wed

* Completely reformatted JSON file... made Services a direct child of Service Provider, and Service Provider is a child element of Service Areas.  Go will load this as a series of maps, which will automatically give us some fast indexing into the data.

### 2015.12.15 - Tue

* Updated ServiceRouter.md documentation.
* Created 

### 2015.12.14 - Mon

* Added design documentation for the Service Router capabilities.  
* Saved to GIT.

### 2015.12.10 - Thu

* Doing other things for a few days.  Coursework, etc.  Back on the case now!
* Current status: "Create" is working for CitySourced API. 
	* Current libraries
		* github.com/ant0ine/go-json-rest/rest - this is working well for the front-end, but does not support XML, and so will not work for the backend.
	* No routing/dispatching is in place yet...  This is probably the next task.

### 2015.12.04 - Fri

* Implemented "napping" REST client.  Not working due to XML.  Starting over using just the HTTP lib.
* Saved to GIT.
* 

### 2015.12.01 - Tue

* 

### 2015.11.30 - Mon

* Wrote up 311 Gateway Proposal.

### 2015.10.27 - Fri

* Reviewed RAML spec with Hassan.  Discussed the Gateway idea.

### 2015.10.24 - Tue

* Revised RAML spec.

### 2015.10.23 - Mon

* First draft of RAML spec.
* Posted to Slack.
