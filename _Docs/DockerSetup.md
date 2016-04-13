# Docker Setup

## Port List

|Program|Instance|Purpose|Port|
|:------|:------:|:------|:----:|
|engine|1|CitySourced Adapter|8080|
|adp_cs|1|CitySourced Adapter <br>San Jose -> CS1 <br>San Jose -> CS2 <br>Santa Clara -> CS3|5001|
|adp_cs|2|CitySourced Adapter <br>Santa Clara -> CS3|5002|
|adp_email|1|Email Adapter (Cupertino, Sunnyvale)|5003|
|cs_sim|1|CS1 - CitySourced Simulator|5051|
|cs_sim|2|CS2 - CitySourced Simulator|5052|
|cs_sim|3|CS3 - CitySourced Simulator|5053|
|monitor|2|System Monitor|5081|

## AWS Setup

|Item|Config|Notes|
|----|----|-----|
|IP|52.34.144.221||
|Port|8080||
|URL|http://52.34.144.221:8080/||

