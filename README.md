# MSDS 434 Final Project

## Data
- This project uses the [NIJ's Recidivism Challenge](https://data.ojp.usdoj.gov/Courts/NIJ-s-Recidivism-Challenge-Full-Dataset/ynf5-u8nk/about_data) full dataset to predict recidivism using a DNN. Using the API endpoint for the data, the 'getFullDataset()' function concurrently pulls the data in 1,000 row chunks. The data is then split into train and test sets based on the training_sample column, which indicates the rows the competition used for training (80%) and testing (20%).