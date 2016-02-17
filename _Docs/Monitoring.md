# Monitoring

## Engine
### Status

1. Name
2. Status (may have a list of engines that _should_ be online)
3. Address
4. Adapters (list of adapter IDs this Engine is aware of)
4. _Request Count_

### Requests

1. RequestID
2. Status (pending, complete, invalid/error)
3. Elapsed time (log StartTime)
4. Type (Create, Search, etc)
5. AreaID

### AdpRequests

1. ID {RequestID, TransID}
2. Status (pending, complete, invalid/error)
3. Elapsed time (log StartTime)
4. Route

## Adapter
### Status

1. Name
2. Type
3. Status
4. _Transaction Count_

### Requests

1. ID {RequestID, TransID}
2. Status (received, waiting for reply, complete, no response, invalid/error)
3. Elapsed time
4. Route


RPC Call

1. Put request in Native format: structs.NCreateRequest
2. r, e := NewRPCCall(service, areaID string, request interface{}, process func(interface{}) error) (*RPCCall, error)
3. r.Run()

Q?) How does Services handle the RPC responses?