# MSDS 434 Final Project

## Data
- This project uses a subset of 10,000 records from the [NIJ's Recidivism Challenge](https://data.ojp.usdoj.gov/Courts/NIJ-s-Recidivism-Challenge-Full-Dataset/ynf5-u8nk/about_data) full dataset to predict recidivism. Using the API endpoint for the data, the 'getFullDataset()' function concurrently pulls the data in 1,000 row chunks. The data is then split into train and test sets based on the training_sample column, which indicates the rows the competition used for training (80%) and testing (20%).

## Model
- The model was created using the XGBoost model in BigQuery ML. BigQuery ML performs automatic preprocessing during traing which consists of feature transformations- categorical variables (strings and booleans) were one-hot-encoded.

> CREATE OR REPLACE MODEL `recidivism.recid_xgb_model`
>       OPTIONS(model_type='BOOSTED_TREE_CLASSIFIER',
>               input_label_cols=['RecidivismWithin3years'])
> AS SELECT * EXCEPT (ID,TrainingSample)
> FROM `recidivism.train_recid`