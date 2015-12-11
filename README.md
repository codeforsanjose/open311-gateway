
## Proposal  

There is growing interest from Cities, Jurisdictions and other Organization to participate with community members through mobile devices and the web.  The Open311 standard is just one example.  

There are a number of potential citizen interactions, for example: 

* Opinions expressed by constituents to their representatives.
* Requests for assistance.

A primary interaction, and the one we are focused on, is the reporting of an “issue” of concern - a “request” by a person to bring attention to a non-emergency problem.  The [Open311][1] standard, [CitySourced][2], and [SeeClickFixed][3] all speak to this.

Currently each City or Jurisdiction has it’s own “311” contact system, hopefully with an openly accessible API.  Because there is an access key required to use a Jurisdiction’s API, and because of variations in Jurisdiction’s implementation or use of 3rd party systems (e.g. SeeClickFix), all mobile and web apps are specific to ONE jurisdiction.

* If you live in San Jose, and work in Cupertino, you currently need to 2 mobile apps. 
* If you are in San Francisco, and want to report an issue, you will need to install a 3rd app for SF.

We propose to build a “311 Gateway”, which will provide a uniform, RESTful, standard API interface.  This will allow mobile developers to code against a single API, yet their app will work across a wider geographic area.  

## Terms

Throughout the remainder of this document, a “Jurisdiction” will be a specific logical destination for a Request specific to a geographical location.

* A location may be serviced by more than one Jurisdiction:
	* The City of San Francisco handles many 311 requests, but there is also a non-profit group that provides trash / dumping pickup.  The City of SF may request the 311 Gateway to route “trash / dumping” issues directly to the non-profit.
* A Jurisdiction may correspond to a Web API interface, or it might be a simple email address.
* The “Jurisdictions” is an array of the Jurisdiction(s) for the users’ location coordinates.

## Features

### Front End

* The front-end API will be a standard HTTP / REST interface.
* JSON initially, with XML added if there is sufficient interest.
* An Application Key is required.
* Initial implementation functions:
	* Service List - get a service list from the Jurisdiction(s) for the current location.
	* Create a Request
	* Find Requests
		* Device ID
		* Near a location
	* Upvote a Request
	* Comment on a Request

### Back End

* Provide standard integrations
	* Open311
	* SeeClickFix
	* CitySourced
* Jurisdiction Registry.
	* Jurisdiction access key.
	* Mapping of locations to Jurisdictions.
* Provide Jurisdiction lookups based on current location.
* Provide tagged service lists for registered Jurisdiction.
* Simple email delivery could be used for cities with no formal 311 system in place.

### Data

To reduce risk and keep things simple, the initial implementation will not persist any personal / private information (e.g. mobile Device Ids, user names, etc).  This includes all log files.

Persisted data will include:

* Data required to determine the set of Jurisdictions for the user’s specified location.
* Jurisdiction access codes.
* Jurisdiction email addresses.
	* System maintenance and updates.
	* Destination if the City does not have a 311 API.
	* Email templates.
* Service lists for each Jurisdiction.  These lists will be tagged with our “major categories”, allowing Apps to present Service Lists in a more user friendly manner.
* Applications authorized to use the 311 Gateway, including their Access Keys.

## Implementation

James Haskell and Hassan Schroeder are interested in coding versions in Go and Elixir, with published side-by-side comparisons at project end.

## Organization

This project could have nationwide utility and interest.  An initial search failed to find any similar apps.  

The project will have some minimal costs: hosting/operational, maintenance and enhancements.  There are two ways to cover these costs:  

* On the BackEnd: Jurisdictions could provide a small recurring payment.  _As most cities and municipalities are strapped for funds these days, this may not be viable.  We want cities to participate… _
* On the FrontEnd: A second option is to require Apps to pay a minimal amount.  Homegrown Apps could be monetized via advertising placement.

The 311 Gateway could be spun off to a 501(c)(4) corp to keep finances clean.  Any excess funds could be donated back to CfSJ or other non-profits.

[1]:	http://www.open311.org/
[2]:	http://www.citysourced.com/
[3]:	http://en.seeclickfix.com/