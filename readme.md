# Go SL Departure Fetcher

This is a MQTT enabled fetcher for departures of SL traffic in Stockholm.
It uses the realtime v4 api from [TrafikLab](https://www.trafiklab.se/api/sl-realtidsinformation-4).

For basic configuration, see `config.example.json`. The `site_id` is filled in using data from the [SL Platsuppslag API](https://www.trafiklab.se/api/sl-platsuppslag). The `direcitons` are pulled from the response of the realtime api. The same goes for `transport_mode`.

## MQTT
Fill in the MQTT broker information, with the server of your choice and the topic you've selected will have data pushed every 30 seconds.
