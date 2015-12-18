## To Do.todo

* Modify CitySourced simulator to return the Request ID and Document ID.
* Outline the Displatch system.

## Log

[2015.12.17 - Thu]

* Router data and indices are working.
* Test OK.
* Saved to GIT.
* Split the test data into two Service Providers for San Jose to test multiple providers.  
* Added providerService to the RouteData struct.  This implements a map of the ServiceID back to the Service Provider.  This provides a fast lookup of the appropriate provider for a New Report, based on the ServiceID.
* Test OK.
* Saved to GIT.

[2015.12.16 - Wed]

* Completely reformatted JSON file... made Services a direct child of Service Provider, and Service Provider is a child element of Service Areas.  Go will load this as a series of maps, which will automatically give us some fast indexing into the data.

[2015.12.15 - Tue]

* Updated ServiceRouter.md documentation.
* Created 

[2015.12.14 - Mon]

* Added design documentation for the Service Router capabilities.  
* Saved to GIT.

[2015.12.10 - Thu]

* Doing other things for a few days.  Coursework, etc.  Back on the case now!
* Current status: "Create" is working for CitySourced API. 
	* Current libraries
		* github.com/ant0ine/go-json-rest/rest - this is working well for the front-end, but does not support XML, and so will not work for the backend.
	* No routing/dispatching is in place yet...  This is probably the next task.

[2015.12.04 - Fri]

* Implemented "napping" REST client.  Not working due to XML.  Starting over using just the HTTP lib.
* Saved to GIT.
* 

[2015.12.01 - Tue]

* 

[2015.11.30 - Mon]

* Wrote up 311 Gateway Proposal.

[2015.10.27 - Fri]

* Reviewed RAML spec with Hassan.  Discussed the Gateway idea.

[2015.10.24 - Tue]

* Revised RAML spec.

[2015.10.23 - Mon]

* First draft of RAML spec.
* Posted to Slack.
