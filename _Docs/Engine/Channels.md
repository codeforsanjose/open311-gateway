## Channels

### Refresh Service Cache
In services/services.go:

__Channel__

	Type: struct{}
	Capacity: 1
	Input: Refresh() calls cache.refresh()  
			cache.update <- true

	Output: init() go func()
			_, ok := <-cache.update
		if !ok: stop
		else: cache.processRefresh()

### Refresh AreaAdapter List
It is possible that an Area (City) might have more than one Service Provider, requiring calls to more than one Adapter.  For example, San Francisco has a primary system, but they also use a non-profit organization for trash cleanup.

So, when the Service Cache is refreshed, we will also rebuild the map of Adapters by AreaID, so that we always have a way to quickly look up the Adapters for an AreaID.

_The following also avoids import cycles between the "router" and "services" packages._

Once the Service Cache refresh is complete, but before switching caches, the following takes place:

1. In services, `cache.processRefresh()` calls `cache.indexAreaAdapters()`.    
	* _`cache.processRefresh()` is the primary function for updating the Services Cache._  
	* _`cache.indexAreaAdaptersbuilds()` a map of AdapterIDs for each Service Area (`map[AreaID][]AdapterIDs`)_. The map is then sent on channel `router.adapters.chAreaAdp` .
2. In router, `init()` has started a go routine that pulls the `data` off the channel, and calls `adapters.updateAreaAdapters(data)`.  
	* `adapters.updateAreaAdapters()` translates the `map[AreaID][]AdapterIDs` to adapters.Adapters.areaAdapters, a `map[string][]*Adapter`.  

__Channel__

	Type: `map[string][]string`  (`map[AreaID][]AdapterIDs`)
	Capacity: 1
	Input: `services.cache.indexAreaAdapters()`





