
## Engine Configuration File
* The default filename for the Engine config is “config.json”, located in the Engine startup directory.  This filename and/or path can be overridden by the “-config” command line option.
* The config file is documented and specified in the “\_Docs/Engine/schema\_config.json”.  This is a [JSON Schema (draft 4)][1] document.  See the “Config File Schema” section below for more information.  NOTE the JSON Schema document is the definitive documentation for the config file.

The config file is a JSON file, having 6 major sections:
* Network - the address the Gateway is running on.
* Auxiliary - Additional programs / processes to be started prior to the Gateway.
* Monitor - the address of the System Monitor, if active.
* General - general configuration settings.
* Adapters - a list of the Adapters the Engine will use.
* Areas - a list of the geographic areas serviced by this Gateway instance.

#### Network
The address the Gateway is running on.

|Setting|Description|s
|:---|:---|
|address|The network address the Gateway presents its API on.  This is typically  the port only, e.g. “:80”.|
|protocol|The protocol presented by the API - currently only HTTP is supported.|

#### Auxiliary
This is a list of JSON objects representing any programs that must be started prior to Gateway.  This is currently used to spin up CitySourced Simulators. _See the “Config File Schema” section below for details._  For each object:

|Setting|Description|
|:---|:---|
|name|Used for status and log messages, this is the name given to this program.|
|autostart|Typically set to “true”, this can be used to disable an Auxiliary program.|
|dir|The directory where the program is located.|
|cmd|The filename of the program or process.|
|args|A list of arguments for the program.|

The “args” setting is a list of strings representing any arguments for the program.  If you were to run this process from the command line directly, the “args” setting is all of the command line options, with each string representing anything separated by a space on the command line.  For example, if the command line is “./progA -debug -config testdir/config1.json”, then the args becomes: “[“-debug”, “-config”, “testdir/config1.json”].


#### Monitor
UDP packets representing various operations can be sent to a System Monitor by the Engine, and by each Adapter.  

|Setting|Description|
|:---|:---|
|address|Status and request information will be sent to the System Monitor on this channel.  This must match the address the System Monitor is running on.  The System Monitor is optional, and is not required.|

#### General
General configuration settings.

|Setting|Description|
|:---|:---|
|searchRadiusMin|The minimum search radius.  Any search radius lower than this amount will be reset to this amount.|
|searchRadiusMax|The maximum search radius.  Any search radius greater than this amount will be reset to this amount.|

#### Adapters
This is a set of JSON objects, each representing an Adapter the Engine is expecting to connect to.

|Setting|Description|
|:---|:---|
|type|The type of adapter - see the JSON Schema for enumerated list.|
|address|The address the Adapter will be communicating on, i.e. the RPC address.  For a local instance, this can just be the port number, like “:5001”.|
|startup|A JSON object like Auxiliary above.  |

#### Areas
The geographic area(s) covered by this Gateway instance (i.e. Engine).  This a set of JSON objects, each of which is the primary ID of a City.  

|Setting|Description|
|:---|:---|
|name|The primary name of the city, like “San Jose”, “San Francisco”, etc.|
|aliases|A list of strings of aliases.  These are case sensitive.|

### Config File Schema
The config file has been documented using [JSON Schema][2].  This file is at “\_Docs/Engine/schema\_config.json”.  

When making changes to the config file, it is highly recommended to validate the changes, both to make sure you have proper JSON, and that the changes match the required structure of the config file.  

Recommended tools:

|Purpose|Tool|
|:---|:---|
|To make sure your changes are valid JSON|[JSON Lint Pro][3]|
|To validate your changes against the expected JSON Schema|[JSON Schema Validator][4]|



[1]:	http://json-schema.org/documentation.html "JSON Schema"
[2]:	http://json-schema.org/documentation.html
[3]:	http://pro.jsonlint.com/
[4]:	http://www.jsonschemavalidator.net/