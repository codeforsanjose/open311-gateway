# Engine and Adapters

Each supported backend has it's own Adapter.  The Adapter is a separate Go program, and the Engine (main program) uses RPC calls to communicate with the Adapter.  All logic specific to communicating with the backend (e.g. Citysourced) is encapsulated in the Adapter.  Communication between the Engine and Adapters uses "Native" data types - there is a Native type to Create a report, and a Native type for Searching for reports, etc.  

### Native Types Naming Convention

|Operation|Request|Response|
|---------|-------|--------|
|Create|NCreateRequest|NCreateResponse|
|Search|NSearchRequest|NSearchResponse|

### RPC Service List

|Service|Method|Name|Notes|
|-------|------|----|-------|
|Service|All|"Service.All"|Retrieves all Services|
|Service|City|"Service.Area"|Retrieves Services for the specified Area|
|Create|Report|"Create.Report"|Creates a new report|
|Search|DeviceID|"Search.DeviceID"|Search for the specified DeviceID|
|Search|Location|"Search.Location"|Search for reports near the specifed geoloc|
|Report|Comment|"Report.Comment"|Add a comment to the specified report|
|Report|Upvote|"Report.Upvote"|Add an upvote to the specified report|




## Implementation

The Engine's Router needs to know some information about the Providers, and hence about the Adapters:

* Providers servicing a location.  _This currently uses the City, but could be more granular in the future._ 
	* City -> Providers (Provider ID)
* All requests (Create, Search, Update, etc) will use the same pattern:
	* Unmarshal the request, and translate to "Native" format.
	* Build the list of Providers for the request. _Note: there can be multiple Providers serviced by a single Adapter._
	* Launch asynchronous calls to for each Adapter/Provider pair.  _Note: the same Adapter may be called more than once!_
		* The asnynchronous calls will return a channel.
	* Wait a limited time (5 seconds?) for all channels to return. _Use a "select"._
	* For each returned result, merge into a master Native struct.
	* Marshal the Native structs into the reply, and respond.



### Adapter

#### Data

Structure:

* serviceAreas
	* map[lower case city name]
		* id
		* name
		* providers[]
			* id
			* name
			* url
			* apiVersion
			* accessKey
			* services[]
				* id
				* name
				* catg[]



### ToDo.todo

* Load Config in Adapters.
* Move all Services to Adapters.



---
_Last updated: 5 Jan 2016_