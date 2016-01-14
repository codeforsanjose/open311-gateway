# Implementation

## Search by Location
There are two searches by location:
1. Search by lat/lng.
2. Search by address.

**Recipe**  
1. Get the City.
2. Get the list of service Providers for the City.
3. Send a "search by location" to each Provider.
4. Merge the results.
5. Return results.


## Search by Device ID

We have two scenarios:  
1. The user's app has kept a list of Provider IDs for all previously created reports.  This Provider ID list is submitted with the search request.
2. The above list is not provided.

**Recipe - Previous Provider List**
1. Send a "search by Device ID" to each Provider.
2. Merge the results.
3. Return the results.

**Recipe - Unknown Previous Providers**
1. Get the City for the current location (or their "home" location if something like that exists [future]).
2. Get the list of service Providers for the City.
3. Send a "search by location" to each Provider.
4. Merge the results.
5. Return results.

## Helper Functions

__Services by Location__  
1. Get City from Geolocation.
2. Look up Service List for that city (`map[city][services]`).  _Data returned includes: service name, category, id._

__Service Providers by Location__  
1. Get City from Geolocation.
2. Look up Service Providers for that city (`map[city][service providers]`). _Data returned includes: provider name.

__Service Provider by Service ID__  
When a report is created, it will contain the Service ID.  We need to quickly get the Service Provider data from the ServiceID.
1. Look up Service Provider from the ServiceID (map[serviceID]*Provider).


__Service Provider Config__
Get the configuration for a Service Provider using their ID.
1. URL
2. Access Credentials


## Data

* The Service endpoint and credentials define a "Service Provider".  That is, if CitySourced was servicing both San Jose and Cupertino, but through a different URL, then we would have 2 Service Provider records.  It is expected that, in general, a Service Provider would be unique to Service Area.
* Most of the detail information below will be on the Adapters. 

---
### Engine
The Engine has the following capabilities:

1. Provide a list of Services available for the specified location.  _Consolidated list queried from all Adapters and cached._
2. Route a Create, Comment or UpVote request to the appropriate Adapter based on the Service MID.
3. Route Search requests:
	1. Search by DeviceID, with IFID/Area list.
     2. Search by DeviceID, using current location.
     3. Search by Current Location.

#### Data

##### City List
|Name|Type|Description|
|----|----|-----------|
|City Code|string|City Code|
|City|[]string|City name|
|AdapterList|[]string|List of RPC destinations|

_The above supports aliases._

__Indexes__
`map[CityCode]CityList`
`map[lower(City)]CityList`

##### Services
|Name|Type|Description|
|----|----|-----------|
|MID|string|Service MID|
|CityCode|string|CityCode|
|Name|string|Description of the Service|
|Categories|[]string|List of categories|

__Indexes__
`map[CityCode][]Services`

##### Request Routing

|Operation|Key|Routing|
|---------|----|-------|
|Services|CityCode|Retrieve cached list of Services for the specified City|
|Create|Service MID|Route to Adapter using IFID in ServiceMID|
|Search DeviceID with previous list|IFID & AreaID from previous Service MIDs|Route to Adapter using IFID; submit DeviceID search to each Adapter and consolidate results.
|
|Search DeviceID / current location|Lookup City from lat/lng|Use the CityList.City index to get the AdapterList; submit DeviceID search to each Adapter and consolidate results.|
|Search by Location|Lookup City from lat/lng|Use the CityList.City index to get the AdapterList; submit Location search to each Adapter and consolidate results.|

---
### Adapters 

#### Service Areas
|Name|Type|Description|
|----|----|-----------|
|ID|int|UniqueID|
|City|string|Unique ID|

#### Service Providers
|Name|Type|Description|
|----|----|-----------|
|ID|int|Unique ID|
|AreaID|int|Mandatory reference to a Service Area.  This is essentially a mandatory foreign key to Service Areas.|
|Name|string|Service Provider's name|
|InterfaceType|string|Specifies the interface subsystem used to communicate with the Provider (e.g. CitySourced, SeeClickFix, etc)|
|URL|string|The endpoint for accessing the Service Provider|
|Key|string|The access key for the endpoint|

#### Services
|Name|Type|Description|
|----|----|-----------|
|ID|int|Unique ID|
|ProviderID|int|The ID of the Service Provider for this service|
|Name|string|Service name|
|Categories|[string]|List of categories. _Used to group services into major categories (see below)._|

#### Service Categories
Used to place the Services into more manageable, user friendly major groups.  __NOTE__: these categories may change, and will be editable on the Gateway Admin page.

|Name|Comments|
|----|---|
|Graffiti||
|Abandoned|Abandoned cars, bicycles, shopping carts, etc.|
|Trash|Trash and illegal dumping.|
|Street|Streets, sidewalks, lighting, etc.|
|Eyesore|Dilapidated or foreclosed homes, illegal signage, etc.|

## Notes

* Initially use a JSON config file.  The file will have the above information in hierarchical order.  Once the system is up and running, and if there is sufficient interest, we can transition to perhaps Postgres or MongoDB.  Postgres has good GeoLoc capabilities baked in.
* Use maps with pointers to structs (`*struct`) to create the "indexes" we need for fast lookups.
