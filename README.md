# MSDS 434 Final Project


## Overview
### Data
This project uses a subset of 10,000 records from the [NIJ's Recidivism Challenge](https://data.ojp.usdoj.gov/Courts/NIJ-s-Recidivism-Challenge-Full-Dataset/ynf5-u8nk/about_data) full dataset to predict recidivism. Using the API endpoint for the data, the 'getFullDataset()' function concurrently pulls the data in 1,000 row chunks. The data is then split into train and test sets based on the training_sample column, which indicates the rows the competition used for training (80%) and testing (20%).

### Model
The model was created using the XGBoost model in BigQuery ML. BigQuery ML performs automatic preprocessing during traing which consists of feature transformations- categorical variables (strings and booleans) were one-hot-encoded.

> CREATE OR REPLACE MODEL `recidivism.recid_xgb_model`
>       OPTIONS(model_type='BOOSTED_TREE_CLASSIFIER',
>               input_label_cols=['RecidivismWithin3years'])
> AS SELECT * EXCEPT (ID,TrainingSample)
> FROM `recidivism.train_recid`

The model took approximately 5 minutes to complete.

## GitHub Repo Contents
### Backend
This directory contains the code for collecting, preprocessing and loading the data into BigQuery. The code is written in Go and requires BigQuery's client API to connect and load the data. A valid project id and BigQuery dataset id in GCP needs to be available.

### Frontend
This container contains code for various pieces for the frontend of the application.
- [predict_micro.go](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/predict_micro.go)- Defines the function to make and retrieve the prediction from BigQuery.
- [handlers.go](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/handlers.go)- Contains handler functions for the Go server microservice.
- [index.html](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/index.html)- html code for the application's homepage.
- [results.html](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/results.html)- html code to display the prediction results.
- [Dockerfile](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/Dockerfile)- Dockerfile to Dockerize the application.
- [prometheus.yml](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/prometheus.yml)- config file that tells the Prometheus Docker container to scrape our prediction application.
- [main.go](https://github.com/mamaOcoder/msds434_final/blob/main/Frontend/main.go)- contains our main function to set up our Go server microservice. It also utilizes the Prometheus client library to the metrics collected for the application.

### GitHub Actions
Attempt to utilize GitHub Actions for linting, however, all attempts failed. Need to further problem-solve.
