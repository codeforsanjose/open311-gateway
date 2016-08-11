# AWS Test Deployment

## API
* A test deployment of the 311 Gateway is on AWS.  It is running at URL:
http://52.34.144.221:80/.
* The test application is openly accessible - it does not require security or authorization.

The API is versioned, currently “v1”.  The following are example API requests (**NOTE: “{baseURL}” is the URL above**):

|URL|Operation|Description|
|:------|:------|:------|
|`{baseURL}/v1/services.json?lat=37.339244&long=-121.883638`|GET|Service List for a location near City Hall in San Jose|
|`{baseURL}/v1/services.json?address=100 First St.&city=San Jose&state=CA`|GET|Service List for the specified address|
|`{baseURL}/v1/requests.json?lat=37.339244&lng=-121.883638&radius=200&comments=1`|GET|Return a list of all requests within 200 meters of lat=37.339244, lng=-121.883638|

NOTE: See the “\_Docs/\_Testing/APITests.sh” document for [cURL][1] examples of each type of API query.

## Configuration
To test a multiple adapter environment, it is configured with two CitySourced Adapter instances providing services for San Jose and Santa Clara, and one Email Adapter providing services for Cupertino and Sunnyvale.  
The following lists all adapters, their config files, and their service lists. (**Note:** _the “Config” column references files relative to the “\_deploy/image” directory in the repo)_:

|Adapter|Instance|Config|ServiceList|
|:------|:------:|:------|:----|
|CitySourced|CS1|adapters/citysourced/config1.json||
|CitySourced|CS2|adapters/citysourced/config1.json||
|Email|EM1|adapters/adapters/email/config.json||

The configuration details  can be found in the “\_Deploy” directory of the project repo.  

**NOTE:** The list of services for each Adapter is listed in “\_Docs/\_Testing/TestAdapterServiceLists.md” file in the repo.

## Usage
The API is meant to follow the [Open311/GeoReport v2][2] spec very closely - please create an Issue on GitHub if you notice a discrepancy.  NOTE: the Gateway API _extends_ the Open311 spec (e.g. Search).

NOTE: the Gateway currently only supports JSON, both in the API requests and in returned data.  If a request specifies “XML”, it will be ignored and JSON used.

Here’s how the application is expected to use the Gateway API:

### Scenarios
#### Create a Request
* Get a list of all services for my current geoloc.
* Choose a service.
* Fill out the new request form.
* Submit

#### Search Before Create
* Search for existing issues at my current geoloc.
* I don’t see the issue, so create it as above.

#### Search For My Previous Issues
* Search by my DeviceID.
* Review all my previously created issues.


### Notes
#### Service List
* This is an HTML **GET** request.
* One of the primary objectives of the Gateway is to provide a uniform service presentation for multiple geographies.  So, when first accessing the Gateway, an application will typically send a request for a Service List for it’s current geolocation.  
* The Service List includes a “group” categorization, providing a coarse grouping of the Services.

#### Create
* Create is an HTML **POST** request.
* There are a number of fields required to create a Request (see the [Open311][3] spec) - these can either be passed as query parameters, or as JSON in the request body.

#### Search
* Search is an HTML **GET** request.
* There are three types of searches:
	1. Search by Lat/Long.
	2. Search by a street address (NOTE: this can be using either separate Address, City, State fields, or a single address string).
	3. Search by Device ID.
* All searches retrieve both open and closed incidents.
* The radius for geographic searches (Lat/Long or Address) is currently limited to between 50 - 200 meters. This is a configuration parameter easily changed.

## Deployment

Details of the AWS deployment are in Quiver / AWS (contact James for these notes).

[1]:	https://curl.haxx.se/
[2]:	http://wiki.open311.org/GeoReport%5C_v2/
[3]:	http://wiki.open311.org/GeoReport%5C_v2/