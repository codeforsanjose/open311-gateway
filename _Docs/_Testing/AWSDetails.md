
## AWS Test Deployment Details

### Port List
The following is the port configuration for the Test Deployment.

|Program|Instance|Purpose|Port|
|:------|:------:|:------|:----:|
|engine|MAIN|CitySourced Adapter|80|
|adp\_cs|CS1|CitySourced Adapter <br>San Jose -\> CS1 <br>San Jose -\> CS2 <br>Santa Clara -\> CS3|5001|
|adp\_cs|CS2|CitySourced Adapter <br>Santa Clara -\> CS3|5002|
|adp\_email|EM1|Email Adapter (Cupertino, Sunnyvale)|5003|
|cs\_sim|1|CitySourced Simulator mimicking a San Jose interface, connected via the CS1 Adapter|5051|
|cs\_sim|2|CitySourced Simulator mimicking a San Jose interface, connected via the CS2 Adapter|5052|
|cs\_sim|3|CitySourced Simulator mimicking Santa Clara, connected via the CS1 Adapter|5053|
|monitor|2|System Monitor|5081|
