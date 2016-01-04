# Multiple Dispatch 

The Gateway supports multiple Providers for any location. So to implement functionality like Search, we need to have a mechanism to query multiple Provider backends, and accumulate and consolidate the queries into one response.

## Constraints

* Multiple back ends.  Each back-end has it's own data structs.
* Everything in Go is strongly typed.  Using interface{} is limited.  We will not be able to pass "generic" data to a back-end processor... In other words, the core processing for each back-end is unique and isolated.
* We can pass functions, and function closures around... this seems like a good way to go.

## Ideas

* Draft skeleton functions, starting at the request/search.go.

## Implementation

1. Create an empty SearchResp struct. All results will be merged back into this struct, one result at a time (use a buffered channel). 
2. Get the list of pertinent Providers.
3. For each Provider:
	* Launch SearchReq.processXX()

### processXX(resp chan)

1. Map request.SearchReq() struct to native back-end struct.
2. Run the HTTP request within a Context.
3. If the request timed-out, return an error as such.
4. MERGE the translated results back into SearchResp struct.
5. Marshal back to JSON and return the response.

Idea: processXX() function sends a func() back on the "resp" channel.  This function is a closure that merges the back-end results back into SearchResp.

OR: send an interface{} back on the channel (pointer to a result struct), and use type assertion to determine the type of struct.  Or maybe send an interface type having the appropriate methods to convert / merge back?




---
## Search
### Search by Location
There are two searches by location:
1. Search by lat/lng.
2. Search by address.

**Recipe**  
1. Get the City.
2. Get the list of service Providers for the City.
3. Send a "search by location" to each Provider.
4. Merge the results.
5. Return results.


### Search by Device ID

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




