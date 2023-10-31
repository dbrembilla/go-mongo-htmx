#!/bin/bash
mongoimport --type csv --db ages --collection ages_collection --file /data/db/ages.csv --headerline