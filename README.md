# FindHotel Coding Challenge

# Notes on the solution

In order to run the code in local clone the .env.dist file to a .env file with for example the following values:
```
PORT=3000
DATABASE_URL=postgres://postgres:root@geolocation_db_1:5432/postgres?sslmode=disable
CSV_URL=https://s3.amazonaws.com/gymbosses/data_dump.csv
```

this will set set the database connection and where to pull the csv from.

after doing this just run: `make run-local` this will run a postgres image and the container with the app.

The CSV file will be downloaded from the CSV_URL, then it will run 10 concurrent gorutines (this can easily be parametrized) parse it, sanitize it and persist it in the db. Since the import of the db is not something we need to care at this moment no effort was done in order to decouple this at the moment.

Once the CSV is persisted it will print in the logs the amount of time it took to parse, persist and also the amount of entries that were duplicated or corrupted/incompleted. 

If I had more time:
- I would debug the error regarding concurrent access to the stats
- I would probably take more time to export this metrics to some external source eg: datadog

### Constraints:
Due to heroku free plan db, im using only a 65k reduced version of the dump that you can find here: `s3://gymbosses/data_dump9k.csv`
you can still use the full version in local setting in the .env `CSV_URL=https://s3.amazonaws.com/gymbosses/data_dump.csv`

The app is running on a free heroku plan, so first time requested it can take some time to respond.

### Metrics in local:
```
With 10 workers 65k:
Total time to parse and persist:  31.1456925s

With 1 worker 65k:
Total time to parse and persist:  1m25.6940387s
```
In order to query the data in production you can do it by executing:

`curl 'https://fh-geolocation.herokuapp.com/api/v1/geoinfo?ip=120.99.153.8'`

or in local:

`curl 'localhost:3000/api/v1/geoinfo?ip=120.99.153.8'`

# Description of the problem

## Geolocation Service

### Overview
You're provided with a CSV file (`data_dump.csv`) that contains raw geolocation data; the goal is to develop a service that imports such data and expose it via an API.

```
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
```

### Requirements
1. Develop a library with two main features:
    * a service that parses the CSV file containing the raw data and persists it in a database;
    * an interface to provide access to the geolocation data (model layer);
1. Develop a REST API that uses the aforementioned library to expose the geolocation data

In doing so:
* define a data format suitable for the data contained in the CSV file;
* sanitize the entries (the file comes from an unreliable source; this means that the entries can be duplicated, may miss some value, the value can not be in the correct format or completely bogus);
* at the end of the import process, return some statistics about the time elapsed, as well as the number of entries accepted/discarded;
* the library should be configurable by an external configuration (particularly with regards to the DB configuration);
* the API layer should implement a single endpoint that, given an IP address, returns information about the IP address' location (i.e. country, city);
* the endpoint should be developed according to the HTTP/1.1 standard;

### Expected outcome and shipping:
* a library that packages the import service and the interface for accessing the geolocation data;
* the REST API application (that uses the aforementioned library) should be Dockerised and the Dockerfile should be included in the solution;
* deploy the project on a cloud platform of your choice (e.g. AWS, Heroku, etc):
    * run a container for the API layer;
    * run any other container that you think necessary;
    * have a database prepared with the already imported data

### Notes
* the file's contents are fake, you don't have to worry about data correctness
* in production the import service would run as part of a scheduled/cron job, but we don't want that part implemented as part of this exercise
* for local/development run a DB container can be included
* you can structure the repository as you see it fit
