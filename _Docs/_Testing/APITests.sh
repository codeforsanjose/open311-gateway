# ============== SERVICE LISTS =====================================

# San Jose lat/long
curl -X "GET" "http://52.34.144.221:80/v1/services.json?lat=37.339244&long=-121.883638" \
  -H "Content-Type: application/json" \
  -d "{}"

# FAIL: Morgan Hill lat/long (not a serviced city)
curl -X "GET" "http://52.34.144.221:80/v1/services.json?lat=37.124975&long=-121.662096"

# San Jose address (address, city, state)
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address=100%20First%20St.&city=San%20Jose&state=CA"

# San Jose address string
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=100%20First%20St.%2C%20San%20Jose%2C%20CA%2092222"

# Sunnyvale address string
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=235%20Olson%20Way%2C%20Sunnyvale%2C%20CA%2094086"

# Cupertino address string
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=20955%20Stevens%20Creek%20Blvd%2C%20Cupertino%2C%20CA%2095014"

# FAIL: Campbell address string (not a serviced city)
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=267%20E%20Campbell%20Ave%2C%20Campbell%2C%20CA%2095008"

# FAIL: Invalid zip code (6 digits)
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=100%20First%20St.%2C%20San%20Jose%2C%20CA%20922222"

# FAIL: Address string for Morgan Hill (not a serviced city)
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address_string=17575%20Peak%20Ave%2C%20Morgan%20Hill%2C%20CA%20"

# Santa Clara address
curl -X "GET" "http://52.34.144.221:80/v1/services.json?address=100%20First%20St.&city=Santa%20Clara&state=CA"


# ============== CREATE =====================================

# CS1-SJ-1-20 (Graffiti Removal) - lat/long in San Jose
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "CS1-SJ-1-20",
  "lat": "37.339244",
  "long": "-121.883638",
  "device_id": "123456789",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "isAnonymous": "true",
  "description": "Gang signs spray painted on the building."
}'

# CS1-SJ-2-50 (Potholes) - lat/long & address string in San Jose
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "api_key": "xyz",
  "jurisdiction_id": "city.gov",
  "service_code": "CS1-SJ-2-58",
  "lat": "37.361780",
  "long": "-121.902293",
  "address_string": "1234 5th street, san jose, CA",
  "first_name": "John",
  "last_name": "Smith",
  "phone": "111111111",
  "description": "A large sinkhole is destroying the street",
  "media_url": "http://farm3.static.flickr.com/2002/2212426634_5ed477a060.jpg",
  "attribute[WHISPAWN]": "123456",
  "attribute[WHISDORN]": "COISL001"
}'

# CS1-SJ-1-20 (Graffiti Removal) - Address string in San Jose
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "CS1-SJ-1-20",
  "address_string": "322 E Santa Clara St, San Jose, CA 95112",
  "email": "jameskhaskell@gmail.com",
  "device_id": "123456789",
  "first_name": "James",
  "last_name": "Haskell",
  "phone": "4445556666",
  "description": "Gang signs spray painted on the building.",
  "isAnonymous": "true"
}'

# CS1-SJ-2-64 (Yard Waste Removal) - lat/long in San Jose
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "CS1-SJ-2-64",
  "service_name": "Yard Waste Removal",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "lat": "37.339244",
  "long": "-121.883638",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "There are piles of trash in my front yard!"
}'

# FAIL: Attempt to create a San Jose Service in Cupertino
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "CS1-SJ-1-39",
  "service_name": "Parking Issue",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "lat": "37.323614",
  "long": "-122.039439",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "There are lots of illegally parked cars around here."
}'

# CS1-SC-1-3 (Abandoned Home) - lat/long in Santa Clara
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "CS1-SC-1-3",
  "service_name": "Abandoned Home",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "lat": "37.350038",
  "long": "-121.946097",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "This home is full of vagrant gangstas."
}'

# EM1-CU-1-10 (Gang Activity) - lat/long in Cupertino
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "EM1-CU-1-10",
  "service_name": "Gang Activity",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "lat": "37.323614",
  "long": "-122.039439",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "Please take care of this Gang Activity!!!"
}'

# FAIL: Attempt to create a Cupertino service in Santa Clara.
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "EM1-CU-1-10",
  "service_name": "Gang Activity",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "lat": "37.350038",
  "long": "-121.946097",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "Please take care of this Gang Activity!!!"
}'

# EM1-CU-2-30 (Illegal Dumping / Trash) - address string in Cupertino
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "EM1-CU-2-30",
  "service_name": "Illegal Dumping",
  "device_id": "123456789",
  "address_string": "10630 S De Anza Blvd, Cupertino, CA 95014",
  "device_type": "iPhone",
  "device_model": "6",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "There\'s a big pile of trash at the end of my block."
}'

# EM1-SUN-1-80 (Sidewalks) - address in Sunnyvale
curl -X "POST" "http://52.34.144.221:80/v1/requests.json" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d $'{
  "service_code": "EM1-SUN-1-80",
  "service_name": "Sidewalks",
  "device_id": "123456789",
  "device_type": "iPhone",
  "device_model": "6",
  "address": "235 Olson Way",
  "city": "Sunnyvale",
  "state": "CA",
  "first_name": "James",
  "last_name": "Haskell",
  "email": "jameskhaskell@gmail.com",
  "phone": "4445556666",
  "description": "The sidewalk in front of my house is cracked and dangerous."
}'

# ============== SEARCH =====================================

# Near SJ City Hall, radius 200m
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?lat=37.339244&lng=-121.883638&radius=200&comments=1"

# Near 101 & 880, San Jose, radius 200m
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?lat=37.361422&lng=-121.900405&radius=200&comments=1"

# Near Santa Clara Courthouse
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?lat=37.350038&lng=-121.946097&radius=200&comments=1"

# Near Costco on Hofstetter, San Jose
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?lat=37.391009&lng=-121.884839&radius=200&comments=0"

# FAIL: Morgan Hill
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?lat=37.125372&lng=-121.662452&radius=200&comments=0"

# DeviceID: 123456789
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?did=123456789&dtype=IPHONE&lat=37.339244&lng=-121.883638&radius=200"

# RequestID: CS1-SJ-1-100
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?rid=CS1-SJ-1-100&lat=37.339244&lng=-121.883638&radius=200"

# RequestID: CS1-SJ-1-101
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?rid=CS1-SJ-2-101&lat=37.339244&lng=-121.883638&radius=200"

# RequestID: CS1-SJ-1-102
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?rid=CS2-SC-1-102&lat=37.339244&lng=-121.883638&radius=200"

# FAIL: invalid query parameter ("deviceID" should be "did")
curl -X "GET" "http://52.34.144.221:80/v1/requests.json?deviceID=12345"
