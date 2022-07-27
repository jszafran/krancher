# Krancher

Go app for crunching the data for employee engagement surveys.

## Running the program

Here's an example command along with parameters required to run the application:

```shell
krancher \
-data resources/itest_data_2x.csv \
-schema resources/itest_schema.json \
-org_structure resources/itest_org.csv \
-workload resources/itest_all_cuts.json
```

where 
* `-data` is a path to CSV file containing the survey data
* `-schema` is a path to a JSON file containing survey schema (metadata about survey questions/demographics)
* `-org_structure` is a path to a CSV file containing org structure
* `-workload` is a path to a JSON file that defines a workload to be processed (cuts by which the data will be computed)
## Running tests

Use `make test` for running application tests.
