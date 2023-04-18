# xmltv

This application serves Electronic Program Guide (EPG) data in XMLTV format.

EPG data is fetched from TV2 Norway. Currently, it has EPG data for norwegian TV channels.

## Demo

I'm hosting this application at https://xmltv.sjurtf.net.
Feel free to use this for your EPG needs.

## Usage

`/channels-norway.xml` - complete list of all available channels

`/{channelId}_{year}-{month}-{day}.xml.gz` - EPG for specific channel on specific day.

Example `/nrk1.nrk.no_2023-04-18.xml.gz`
