CREATE TABLE geoinfo(
	ip VARCHAR (200) UNIQUE NOT NULL PRIMARY KEY,
	country_code VARCHAR (200) NOT NULL,
	country VARCHAR (200) NOT NULL,
	city VARCHAR (200) NOT NULL,
	latitude VARCHAR (200) NOT NULL,
	longitude VARCHAR (200) NOT NULL,
	mystery_value VARCHAR (200) NOT NULL
   );