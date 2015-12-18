# Service Router
The Service Router subsystem is responsible for:  

* Maintaining the list of Service Providers and their Services.
* Retrieving the Services list for the user's location.
* Routing user Requests and Queries to the appropriate Service Provider.

## New Report
The user application should identify the appropriate location as soon as possible, and relay this to the Gateway.  For a mobile app, the current geolocation will typically be used, but the user could also specify a different address.  The  location is submitted to the Gateway using the "/services" endpoint.  The list of Services available for the current location will be returned.  The Services list will contain an ID and a description.  It may also contain our own tags (TODO).  
The user will select the appropriate Service, and define the issue.  The issue, with the ServiceID, is submitted using a POST to the "/requests" endpoint. The ServiceID will be used to route the new request to the appropriate Service Provider (i.e. Jurisdiction or Authority).
The Service Provider ID will be returned to the App in the response JSON data.  The Service Provider ID should be persisted in a list of unique Service Provide ID's (see "Search / Device ID" below).

## Search

### Device ID
If the user wishes to review all of their issues, then we need a list of Service Provider ID's for all of their previously created issues.  This will be passed in the Search request GET at the "/requests" endpoint.  The Gateway will send appropriate search requests to each Service Provider for the specified DeviceID, consolidate the returned data, and send it back in the GET response.  If the list of previous Service Providers is not provided, then the current City will be 


### Location
A search by location will build a list of all Service Providers for the specified location, send search requests to each, consolidate the results, and return that in the GET response.


## App Requirements
Applications using the Gateway must:  
1. Provide access credentials for each request to the Gateway.

To provide the best user experience, applications should also:
1. Persist a list of unique Service Provider ID's belonging to  any created issues.  _This will be used in the "Search by DeviceID" discussed above._

## Implementation

### Queries

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


### Data

__NOTE:__ The Service endpoint and credentials define a "Service Provider".  That is, if CitySourced was servicing both San Jose and Cupertino, but through a different URL, then we would have 2 Service Provider records.  It is expected that, in general, a Service Provider would be unique to Service Area.

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

### Notes

* Initially use a JSON config file.  The file will have the above information in hierarchical order.  Once the system is up and running, and if there is sufficient interest, we can transition to perhaps Postgres or MongoDB.  Postgres has good GeoLoc capabilities baked in.
* Use maps with pointers to structs (`*struct`) to create the "indexes" we need for fast lookups.
