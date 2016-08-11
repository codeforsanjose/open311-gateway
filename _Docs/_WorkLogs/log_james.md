## To Do.todo

* Convert ServiceID.ID to string.
* In Adapter config files, change "type" to "responseType".
* In Adapter data, convert "serviceCategories" to "serviceGroup" for consistency.
* SeeClickFix adapter.
* Quick Start guide.
* Allow logs to be configured to go to the System Logger (MacOS & Linux).
* Test running the Engine on HTTPS.
* Spin up a test implementation. @done(2016-04-12)
* Auto start & stop adapters.  @done(2016-03-31)
* Finish Service Definitions, including optional attributes on a POST Service Request.   @done(2016-04-12)
* Stress tester.
*  Improve error handling. @done(2016-03-31)
* Fix dates on Searches. @done(2016-03-31)
* Run race detector.  _Ran OK!  Will start running with "-race" flag for all testing._ @done(2016-03-18)
* Rework engine/router/rpcroute.go - the map idea seemed like a good idea at the time, but it's difficult to understand and debug. @done(2016-03-18)
* Implement Report Upvote
* Implement Report Comment
* Implement report searches: @done(2016-02-29)
	* Single ID (SearchID) @done(2016-02-29)
	* DeviceID  (SearchDID) @done(2016-02-29)
	* LatLng    (SearchLL) @done(2016-02-23)
	* Address (converted to LatLng in engine... submitted to Adapter as SearchLL)
* Update RAML file for all current Query Parms.
* On the engine, search needs an index of CityCode -\> Providers.  When the ServiceList is loaded, we need to create that index. @done(2016-02-04)
* * Get Create working again. @done(2016-01-25)
* Update Services cache on Engine. @done(2016-01-22)
	* Be able to safely and easily refresh the services list. @done(2016-01-22)
	* Quickly lookup location -\> AreaID -\> Service list. @done(2016-01-22)
	* Quickly lookup location -\> AreaID -\> Adapter list. @done(2016-01-22)
* Implement RPC Dispatch System @done(2016-01-19)
* Custom Type for ServiceID. @done(2016-01-11)
* CitySourced Adapter - load config.json data file. @done(2016-01-05)
* Separate back-end interfaces as separate apps, using go-rpc calls. @done(2016-01-05)
* Implement Gingko BDD tests. @done(2015-12-18)
* Modify report.Create() to use ServiceID and JID. @done(2015-12-28)
* Consolidate all Create functionality within the request.CreateReq type. @done(2015-12-28)
* Implement Services API endpoint. @done(2015-12-18)
	* Bring Google Maps API over. @done(2015-12-18)
* Modify CitySourced simulator to return the Request ID and Document ID. @done(2015-12-12)
* Outline the Displatch system. @done(2015-12-14)

---- 

## Log

### 2016.08.11 - Thu

* Moved all code be under "github.com", rather than a local path.
* Cleaned up all imports to work properly under github.com.
* Compiles OK with "make maccompile".

### 2016.06.09 - Thu

* Ari found a CORS issue.  See [here](http://stackoverflow.com/questions/22972066/how-to-handle-preflight-cors-requests-on-a-go-server) for a possible solution.
* Installed CORS, recompiled and redeployed to AWS.

### 2016.04.30 - Sat

* Fixed Services query.  Was returning the same name, group, etc. for all services.  Broken when making the empty values return JSON null… 
* Fixed all instances of “long” to “lng” in the engine, as per Ari.

### 2016.04.19 - Tue

* Organizing documentation.
* Moved SearchRadius from a static constant to a configuration parameter.
* Saved to GIT.
* Added doc for Engine Configuration.
\* 

### 2016.04.13 - Wed

* Saved to GIT.
* Engine: added config settings for network, including (untested!) HTTPS support.
* Saved to GIT.

### 2016.04.12 - Tue

* Created JSON Schema for Engine and Adapter config files, to document and allow easy verification of their structure.
* Created a copy of the CitySourced adapter code as a starting point for SeeClickFix adapter.

### 2016.04.11 - Mon

* Docker build working.
* Saved to GIT.
* Docker container running on AWS!

### 2016.04.08 - Fri

* Make file for building Linux deployment image (or Mac for testing) complete.
* Saved to GIT.

### 2016.04.07 - Thu

* Deployment cleanup.
* Converted hardcoded "monitor" address to config files.
* In engine, added startup for auxiliary programs.  This can be used to start the CitySourced simulator for testing purposes.
* Saved to GIT.

### 2016.04.06 - Wed

* Converted logging in Adapters to jeffizhungry version of logrus.
* Cleanup
* Started Deployment / Docker.
* Saved to GIT.

### 2016.03.31 - Thu

* Modified Services response to return empty strings as null.
* Saved to GIT.
* Modified Services response to return empty strings as null.
* Saved to GIT.
* Replaced logging system (replaced go-logging with jeffizhungry version of logrus).
* Removed some debug prints.
* Saved to GIT.
* Added Adapter auto-start system.
* Saved to GIT.

### 2016.03.30 - Wed

* Modified go-json-rest lib so that errors are returned in the Open311 format.
* Modified Search to return "null" JSON values instead of empty strings.
* Saved to GIT.

### 2016.03.29 - Tue

* Address validation (in engine/common) working well, including ServiceArea validation.
* Saved to GIT.
* Create requests working with all validations, both CitySourced and email.
* Added a service list to the engine/services system.  This allows quick checks on ServiceID's.
* Saved to GIT.

### 2016.03.28 - Mon

* Updated Address validation.

### 2016.03.21 - Mon

* Revised ServiceList to closely match Open311-GeoReportV2 spec.
* Code cleanup (metalinter).
* Saved to GIT.
* Modified Create to match Open311/GeoReport2 spec RE: JSON names.
* Saved to GIT.
* Modified Error returns to match Open311/GeoReport2.
* Modified Service query to more closely match Open311/GeoReport2.  Full address string and addr, city, state query parms all working.
* Removed multiple image resolution ImageURLxxx.
* Saved to GIT.

### 2016.03.17 - Thu

* Email adapter working end to end.  
* Saved to GIT.
* In email/mail/send.go, change Send() function to launch a go routine to send email (speed and performance).
* Cleanup.
* Saved to GIT.

### 2016.03.14 - Mon

* Email adapter working at the unit test level.
* Saved to GIT.

### 2016.03.10 - Thu

* Started code for Email Adapter.  
* Email send via Gmail working.
* Create email template working with Send().
* Saved to GIT.


### 2016.03.09 - Wed

* Saved to GIT.
* Cleaned up engine/services to use new "Mgr" style requests, with new RPC system.
* Cleaned up "common" files.
* Added timeout to http.Post calls in Adapters.
* Saved to GIT.
* New, cleaner, lighter RPC/Mgr method in place.  Old RPC code deleted.
* Saved to GIT.
* Cleaned up Service Cache refresh().  Converted to "Mgr" style.
* Saved to GIT.

### 2016.03.08 - Tue

* RPC system cleaned up.  First draft - needs testing.  Old code is still mixed in.
* Saved to GIT.
* Deleted old code from engine/request/search.go
* Cleaned up engine/request/create.go, and converted to new RPC calls.

### 2016.03.07 - Mon

* Started clean up of engine/router RPC system.

### 2016.03.05 - Sat

* In engine/request/search.go:
	* In searchMgr type:
		* Added reqType structs.NRequestType
		* Removed srchType
		* This will make the request managers more uniform.
* Saved to GIT.

### 2016.03.04 - Fri

* Cleaned up Search().  It's looking pretty good now.
* Modified engine/router/data.go:
	* In the routeData struct, changed "routes" to "indArea" - this is a more descriptive name for it's purpose.
	* Added an "all" member that holds all active routes.
	* Modified the routeData.update(), and the function cache.sendRoutes() in engine/requests/services.go.  This function now simply builds a list of unique routes and sends it to router/data.go, rather than building a map of routes indexed by AreaID.  
	* router/data.go saves the full, unique route list, and then builds the indexed list of routes by AreaID.
* Added new file router/routes.go.  This file returns route lists (all, by Area, and by ReportID).
* Saved to GIT.
\* 

### 2016.03.03 - Thu

* Create request cleaned up, and using the "Manager" pattern.
* Built Validation and Conversion objects.
* Replaced existing conversion routines with a single function (engine/request/common.go - conversion.convert()).
* Saved to GIT.
* Removed GetSID() call from engine/request/request.go, and moved to Create Mgr.  Also moved telemetry calls to CreatMgr.  This keeps the request.go concise and clean, and pushes all request processing to the Create Mgr where it belongs.  
* Saved to GIT.
* First round cleanup to Search complete.  Created Search Mgr.
* Saved to GIT.

### 2016.03.02 - Wed

* Started cleaning up engine/request.  It will more closely follow the model used in the Adapters of using "Request Managers" to process requests.  
	\*

### 2016.03.01 - Tue

* Cleaned up RequestID - all messages, including Service List Requests, now have a message ID.
* Moved the Message ID (SID) mechanism to new file engine/router/sid.go. 
* Minor fixes in the Monitor program.
* Saved to GIT.

### 2016.02.29 - Mon

* Search by ReportID working.
* Added RID (ReportID) to structs.go.
* Saved to GIT.
\* 

### 2016.02.24 - Wed

* Search by DeviceID working.  Needs thorough testing.
* Saved to GIT.


### 2016.02.23 - Tue

* Added panic/recover functions to the Services(), Create() and Search() functions in engine/request/request.go.
*  Modified Engine Search - in engine/request/search.go:
	* Added srchType field to SearchRequest struct.  
	* Added constants for the valid srchType's (ReportID, DeviceID, and LatLng).
	* Cleaned up JSON/XML tags for the SearchRequest struct.
	* SearchRequest.validate() calls the new function setSearchType().  This function ascertains the type of search based on which query parms were used, in the following priority: ReportID, DeviceID, LatLng.  In other words, if the ReportID is present, then that will be search type.  If ReportID has not been set, or appears invalid, and there is a valid DeviceID, then that will will be the search used.  If neither ReportID or DeviceID searches look ok, then LatLng will be checked.  If the Lat & Lng don't look valid, then the search request is rejected.
* Saved to GIT.
\* 

### 2016.02.22 - Mon

* In structs, changed the ID in responses from an int to a ReportID (new struct).  The ReportID includes the Route, so an upvote or comment to a specific report can be properly routed.
* Saved to GIT.
* Deleted engine/request/\_old code.
* Saved to GIT.
* Fixed Search function on the Engine - it was not consolidating multiple returns from Adapter Routes. 
* Saved to GIT.

### 2016.02.21 - Sun

* Telemetry and Message ID's working pretty well now for Create and Search.  
* Added result count to telemetry.
* Deleted old engine router code.
* Saved to GIT.

### 2016.02.19 - Fri

* Added Adapter monitoring messages and calls.
* Saved to GIT.
* Adapter RPC monitoring starting to work.
* Saved to GIT.


### 2016.02.18 - Thu

* Request monitoring working.
* Starting to make changes to support IDs in RPC requests.
* Saved to GIT.

### 2016.02.17 - Wed

* Saved current work on the Monitor program to git.
* Fleshed out messages and data:
	* Engine Status
	* Engine Requests
	* Engine Status
* Saved to GIT.
* Created a new package "comm".  This will contain all the basic communication - messages, udp/network, etc.
* Moved all Message handling to "comm".
* Test OK.
* Saved to GIT.
* Renamed "comm" package to "telemetry".  This will match with "engine" and "adapter" usages.
* Created hard link of monitor/telemetry/message.go in engine/telemetry and adapter/citysourced/telemetry.
* Saved to GIT.
* In Monitor:
	* Cleaned up network connection (telemetry/network.go).
	* Changed "AdpEngRequest" to "EngRPC" (more concise and descriptive).
	* Changed MsgTypeEA to MsgTypeERPC to match above.
	* Moved all network related code into telemetry/network.go.
* In Engine:
	* Update telemetry - added SendEngRequest() and SendEngRPC() calls.
* Saved to GIT.
* Cleaned up Monitor:
	* Deleted old test and development code.
	* Added telemetry.Start() and display.Start() calls.
	* Initialization and shutdown is cleaner now.
* Saved to GIT.
* In monitor/display/display.go: added go func() to process message chan (incoming message queue).
* Engine starting to communicate with Monitor...
* Saved to GIT.

### 2016.02.16 - Tue

* Building Monitor program using uiTerm.

### 2016.02.12 - Fri

* Started building monitoring app. This will make it easier to verify routing is working properly, as well as monitoring critical processes in the Engine and Adapters.
* Saved to GIT.


### 2016.02.11 - Thu

* Started termui monitor interface.

### 2016.02.08 - Mon

* Saved to GIT.

### 2016.02.05 - Fri

* In engine/router, changed "routeMap" to "serviceMap".  This is more descriptive.
* Test OK.
* Saved to GIT..
* Changed the "map[structs.NRoute]\*rpcAdapterStatus" to a custom type in engine/router/rpc.go and rpcroute.go.  
* Test OK.
* Saved to GIT.
* Cleaned up debug prints.
* Saved to GIT.
* Cleaned up routing.
* Reversed accidental name refactor.

### 2016.02.04 - Thu

* Cleaned up debug display of routes.
* Created NRouteType enumeration.
* Added RouteType() to NRouter interface.
* In engine/request/create.go, merged CreateReqBase into CreateReq, and dropped the CreateReqBase struct. Unnecessary complication.
* Cleaned up engine/request/search.go.
* Changed many method receivers to be "r" for consistency and ease of readability.
* Revised engine/router/rpc.go to put the route properly into each outbound RPC request - this involved making a copy of each struct.
* Test OK.
* Saved to GIT.
* Created "routes" in the Cache in engine/services/services.go.  This is a map of Routes by AreaID.  This will be used by Search and any other functions requiring routing to an Adapter/Provider for an Area.
* Removed all references to "bkend" in any of the data structures.  This has been replaced by structs.NRoute(s).
* Test OK.
* Saved to GIT.
* In engine/router:
	* Added routeData to data.go.  Contains all routing data, and is updated everytime the Service List is updated.  A channel is used to pipe the update from the Services package to the Router package (necessary to avoid an import cycle).
	* Reworked engine/router/rpcroute.go to support the use of routeData.
* Removed route data from engine/services/services.go.
* In structs.go, added NResponseCommon (analogous to NRequestCommon).  Also NResponser interface and NResponseType enumerated type.
* Saved to GIT.
* SearchLL rudimentarily working!
* Saved to GIT.

### 2016.02.03 - Wed

* Creating a "wrapper" for the rpc calls (structs.NRequestPkg) is not working... the rpc/gob system on the Adapter (client) is rejecting the RPC call seemingly because it doesn't know what to do with the Request interface values - and registering those types in gob isn't helping.  So... backing up to a previous commit, and will put the common request type as an anonymous struct into NServiceRequest, NCreateRequest, etc. Keep things simple...
* All is good again - Services and Create working again.   Now back to Search.
* Saved to GIT.
* Route is now being engraved in each RPC request.
* Test OK.
* Saved to GIT.

### 2016.02.02 - Tue

* Changed structs.NRoute to always be a slice, so that multiple routes can be supported.  This will be needed for Search by Device ID. 
* Started organizing the Search in adapters/citysourced.
* Test OK.
* Saved to GIT.
* Reverted the above change.  This pushes routing logic into the Adapters, which is not a good road.  The Adapters should be very simple, mediating a single request.  There should not be any routing logic in the Adapters...
* First draft of SearchLL() ready for testing in adapters/citysourced.
* Saved to GIT.
\* 

### 2016.02.01 - Mon

* Reworked the Create code in the CitySourced Adapter.  There is now a very clean division between the Normal structs, the native structs, and process managers.  Create working very well with this new division.
* The native request and response for Create are now in the create directory package.
* Test OK.
* Saved to GIT.
* In adapters/citysourced/request.go, created an interface "processer", and runRequest() for the main processing steps for all requests.
* Test OK.
* Saved to GIT.
\* 

### 2016.01.28 - Thu

* Added NSearchType to structs.
* Added "go generate" with stringer tool to automatically generate the String() method for NSearchType.

### 2016.01.27 - Wed

* Revised RPC management in engine/router.  The router is much more self contained, and will properly route requests to the proper Adapter based solely on information within the RPCCall struct.  
* Test OK.
* Saved to GIT.
* Changed IFID to AdpID everywhere for clarity.
* Deleted unused, commented code.
* Removed some Debug logs.
* Test OK.
* Saved to GIT.
* More debug print cleanup.
* Saved to GIT.
* More cleanup:
	* Removed "arith" sample code from Adapter code.
	* Cleaned up all import statements.
	* In engine/router/rpc.go, renamed RPCCall.listIF to RPCCall.adpList.  This is in line with the terminology and naming elsewhere in the project.
* Added color coding capability to common.LogString and log.LogString.
* Test OK
* Saved to GIT.


### 2016.01.25 - Mon

* Create is working.
* Saved to GIT.

### 2016.01.22 - Fri

* Modified services/services.go to create a map of AdapterIDs by AreaID, and send it back to router/adapters.go for updating.
* See new doc: "channels.md" for details.
* The adapters are frequently referenced by AdapterID, so I changed engine/data/config.json to use objects for each Area, rather than a list.
* Changed Adapters.Name to Adapters.ID in router/adapters.go, to match the current usage and docs.
* Test OK.
* Saved to GIT.


### 2016.01.21 - Thu

* Modified ServicesList refresh to go through a channel, so that only one update can possibly be running simultaneously.  To start a refresh, call router.RefreshServicesList().
* Saved to GIT.
* Moved services.go from the engine/router package into it's own, new package "services".  All code related to the ServiceList cache system will be in this package.
* The services package needs access to the RPC call system, so revised engine/router/rpc.go - made the following exported:
	* newRPCCall -\> NewRPCCall
	* rpcCalls -\> RPCCall
	* rpcCalls.run() -\> RPCCall.Run()
* Restored engine/request/services.go, and updated to match all the changes in the Engine.
* Test OK.
* Saved to GIT.
* Renamed engine/router/adapters.go -\> data.go.  There is more in the config.json file than just adapter data.
* Fixed Services response.
* Test OK from Paw client, with location in San Jose.
* Saved to GIT.
\* 


### 2016.01.20 - Wed

* Service List update working better.  Needs better locking... multiple accesses are possible now.

### 2016.01.19 - Tue

* Changed all variables and types name "city" to "area".  For now, "areas" correspond to cities, but there is a good chance that will change in the future.  
* Got Humpty Dumpty put back together.  RPC dispatch system is now working.
* RPC Dispatch System is in new file engine/router/rpc.go.
* Saved to GIT.

### 2016.01.15 - Fri

* Services List working from Engine and Adapter (Citysourced)!.
* Set up 2 Adapters (CS1 and CS2) for testing.
	* CS1 = San Jose
	* CS2 = Santa Clara
* Removed the LogString (for boxed struct printing) from logs.  It was already in common.  Updated the code in common with comments, and removed the LogPrinter stuff (not used).
* Added command line args (flags) to the Engine, as well as SignalHandlers.
* Saved to GIT.

### 2016.01.14 - Thu

* Implemented go-logging in adapters/citysourced.
	* Converted all fmt.Print statements to log statements.
* Command line options working now.
* Signal handler working.
* Saved to GIT.
* Renamed functions in adapters/citysourced/data.go to show they are about Services:
	* City() -\> ServicesCity()
	* All() -\> ServicesAll()
* Removed debug from config file.  It's redundant and unnecessary.
* Saved to GIT.


### 2016.01.13 - Wed

* Renamed directories:
	* gateway -\> engine
	* integration -\> adapters
* Started rebuilding the Engine.
* Saved to GIT.
\* 


### 2016.01.12 - Tue

* In integrations/citysourced/request/services.go:
	* Change Service.ServicesForCity() to City().
	* Added Service.All() - returns a list of ALL services.

### 2016.01.11 - Mon

* In integration/citysourced/data/data.go;
	* All indexes working now.
	* Deleted the "serviceProvider" index... it is no longer needed with the new Master Service IDs.
	* Added a Santa Clara section of test data in config.json.
* Moved the structs.go file into it's own package: integration/citysourced/structs/structs.go.  Had a circular import issue.
* Saved to GIT.
* RPC test working OK with Services API.
* Saved to GIT.
* Added second CitySourced integration for testing purpose.
* Saved to GIT.
* Changed gateway/request/\*.go:
	* All fields need to be exported in the Create structs, so change "longitude", "latitude", etc. to "LongitudeV", "LatitudeV", etc.
* Moved gateway/request/structs.go to gateway/structs/structs.go
	* Added gateway/structs/methods.go.
\* 
 
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
* Changed \_Test/CSclient.go to use asynchronous API calls.  Working OK.
* Saved to GIT.

### 2016.01.04 - Mon

* Created separate directories for the Engine and Adapters, and separated the existing source code.
* Saved to GIT.
* Worked on design of Engine/Adapter RPC system - responsibilities, design.
* Got the CitySourced client rudimentarily working.  
* Put all "Native" structs into a the request/structs.go file.
\* 

### 2015.12.30 - Wed

* Renamed geo functions:
	* LatLngForAddr() -\> LatLngForAddr()
	* AddrForLatLng() -\> AddrForLatLng()
	* CityForLatLng() -\> cityForLatLng()

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
		* CreateReport type -\> CreateReq
		* CreateReportResp type -\> CreateResp
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

* Cleaned up router.ServiceProviderInterface().  Added test cases in data\_test.go.
* Test OK.
\* 

### 2015.12.19 - Sat

* In router/data.go:
	* Dropped "Getxxx" from the ServiceXXX() methods.
	* Added ServiceProviderInterface() go get the Service Provider interface type (currently on CitySourced).  
* in request/report.go, added a map lookup (beCreate) on the Create functions.  We will have a set of these maps to quickly route an incoming request (create, lookup, etc), to the correct backend interface.
* Test OK.
* Saved to GIT.

### 2015.12.18 - Fri

* Brought "geo" package over from CitySourced.
* Added "getCity()" function to mygeocode.go.  This scans through the Google response and retrieves the city.  We will need this for quickly mapping the Mobile Apps geoloc -\> city -\> Service Providers -\> list of Services.
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
\* 

### 2015.12.01 - Tue

\* 

### 2015.11.30 - Mon

* Wrote up 311 Gateway Proposal.

### 2015.10.27 - Fri

* Reviewed RAML spec with Hassan.  Discussed the Gateway idea.

### 2015.10.24 - Tue

* Revised RAML spec.

### 2015.10.23 - Mon

* First draft of RAML spec.
* Posted to Slack.