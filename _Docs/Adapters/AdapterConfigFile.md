## Adapter Configuration File
* This document is an meant to present an overview and outline of the configuration file structure.  The config file is precisely documented in the JSON Schema file at “\_Docs/Adapter/schema\_config.json” file.  **The JSON Schema file is the definitive documentation of the Adapter config files.** It is a [JSON Schema (draft 4)][1] document.  See the “Config File Schema” section below for more information.
* The default filename for an Adapter config is “config.json”, located in the Adapter startup directory.  This filename and/or path can be overridden by the “-config” command line option.

The config file is a JSON file, having 4 major sections:
* Adapter - general setup of the Adapter, including its RPC address.
* Monitor - the address of the System Monitor, if active.
* Service Groups - for Providers having static Service Lists (like Email and CitySourced), this is the list of Service Groups (categories of Services).
* Service Areas - the list of Providers for each geographic area serviced by this Adapter. This is almost always going to be a single Provider… for example, it would be unusual to have CitySourced providing services to San Jose from two API backends.

#### Adapter
General configuration information (name, type), including the RPC address.

|Setting|Description|
|:---|:---|
|name|The name of the adapter.  This must match the name in the Engine configuration file.|
|type|The type of Adapter (e.g. “CitySourced”, or “Email”).  See the JSON Schema file (“schema\_config.json”) for an enumerated list of possible settings.|
|address|The network address for the RPC connection to the Engine.  If the Engine and Adapter are running on the same server, then this can be the port only, e.g. “:5001”.|

#### Monitor
UDP packets representing various operations can be sent to a System Monitor by each Adapter.  

|Setting|Description|
|:---|:---|
|address|Status and request information will be sent to the System Monitor on this channel.  This must match the address the System Monitor is running on.  The System Monitor is optional, and is not required.|

#### Service Groups
The list of Service Groups, used to populate the Service List cache.  This only applies to Adapters with static Service Lists, such as CitySourced or Email.  Other Providers, such as SeeClickFix, provide dynamic, query-able Service Lists - in their case, the Service List cache is dynamically populated.

This is simply an array of strings, like: [“Abandoned”, “Eyesore”, “Graffiti”, “Streets”, “Trash”].

#### Service Areas
This is a set of JSON objects, each listing the Providers servicing the geographic area.  The Provider object contains all information necessary to access the Provider API.  If the Provider does not provide Service List queries, that is they require static Service Lists, then the Service List will be included under the Provider.  See the “schema\_config.json” for more details. 

The following sections are nested JSON objects under “serviceAreas”:

|Setting|Description|
|:---|:---|
|&lt;AreaID&gt;|The areaID is used to create our ServiceID’s.  Example: “SJ” for San Jose; “SC” for Santa Clara.|

The AreaID object contains two fields:

|Setting|Description|
|:---|:---|
|name|The common name for the area.  Example: “San Jose” for “SJ”.|
|providers|The list of Providers servicing this Area, under this Adapter type.  NOTE: it is possible that a city would have multiple types of Providers.|

The Providers list contains objects:

|Setting|Description|
|:---|:---|
|id|Sequential, numerical ID.  NOTE: this should never be changed!  This value is used to build the ServiceIDs and RequestIDs, which may be stored and referenced by a mobile device|
|name|The name of the Provider.  Used in log and error messages, and should be a unique, understandable name.|
|url|The base URL of the Adapter API.  For example, the base URL of the Open311 version of SeeClickFix is “[https://seeclickfix.com/open311/v2/][2]”|
|apiVersion|The API version number for all requests.  For example, CitySourced requires this in their XML request payload.  Optional.|
|key|The provider’s API key.|
|responseType|The default response type, as per the Open311 spec (“realtime”, “blackbox”, etc).|
|services|If the Provider does NOT have query-able Service Lists, then this contains a static list of Services.|

Static Service List, used for Providers who do not provide a Service List query.

|Setting|Description|
|:---|:---|
|id|A number identifying the Service.  This ID must be unique within the Provider, and should never be changed!  If a Service is eliminated, the delete this item, but do not reuse the this ID within this Provider.|
|name|The name of the Service, e.g. “Abandoned Bicycle”, or “Graffiti”.|
|description|The description of the Service.|
|group|The groups the Service should appear within (typically only one). This must be one of the strings from the Service Groups section (see above).|

### Config File Schema
The config file has been documented using [JSON Schema][3].  This file is at “\_Docs/Engine/schema\_config.json”.  

When making changes to the config file, it is highly recommended to validate the changes, both to make sure you have proper JSON, and that the changes match the required structure of the config file.  

Recommended tools:

|Purpose|Tool|
|:---|:---|
|To make sure your changes are valid JSON|[JSON Lint Pro][4]|
|To validate your changes against the expected JSON Schema|[JSON Schema Validator][5]|




[1]:	http://json-schema.org/documentation.html "JSON Schema"
[2]:	https://seeclickfix.com/open311/v2
[3]:	http://json-schema.org/documentation.html
[4]:	http://pro.jsonlint.com/
[5]:	http://www.jsonschemavalidator.net/