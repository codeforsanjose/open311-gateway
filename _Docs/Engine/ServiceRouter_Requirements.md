# Service Router
The Service Router subsystem is responsible for:  

* Maintaining the list of Service Providers and their Services.
* Retrieving the Services list for the user's location.
* Routing user Requests and Queries to the appropriate Service Provider.

## Create
The user application should identify the appropriate location as soon as possible, and relay this to the Gateway.  For a mobile app, the current geolocation will typically be used, but the user could also specify a different address.  The  location is submitted to the Gateway using the "/services" endpoint.  The list of Services available for the current location will be returned.  The Services list will contain an ID and a description.  It may also contain our own tags (TODO).  
The user will select the appropriate Service, and define the issue.  The issue, with the ServiceID, is submitted using a POST to the "/requests" endpoint. The ServiceID will be used to route the new request to the appropriate Service Provider (i.e. Jurisdiction or Authority).
The Service Provider ID will be returned to the App in the response JSON data.  The Service Provider ID should be persisted in a list of unique Service Provide ID's (see "Search / Device ID" below).

## Search

### Device ID
If the user wishes to review all of their issues, then we need a list of Service Provider ID's for all of their previously created issues.  This will be passed in the Search request GET at the "/requests" endpoint.  The Gateway will send appropriate search requests to each Service Provider for the specified DeviceID, consolidate the returned data, and send it back in the GET response.  If the list of previous Service Providers is not provided, the search will be limited to the current City.


### Location
A search by location will build a list of all Service Providers for the specified location, send search requests to each, consolidate the results, and return that in the GET response.


## App Requirements
Applications using the Gateway must:  
1. Provide access credentials for each request to the Gateway.

To provide the best user experience, applications should also:
1. Persist a list of unique Service Provider ID's (the first two fields of the Service MID) belonging to  any created issues.  _This will be used in the "Search by DeviceID" discussed above._


