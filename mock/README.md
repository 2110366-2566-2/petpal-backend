# Use to generate mock data in 'petpal' database

## Usage

1. install `exrex` lib in python
    ```bash
    pip install exrex
    ```
1. [config](#config) variable in `mocker.py`
1. run `mocker.py`

## Configuration
- `N` : number of records to be generated
- `COLLECTION_NAME` : name of the collection
    - **Note** must create schema file first (`COLLECTION_NAME_schema.json`)
- `CLEAR_PREV` : clear previous data in the collection
- `SEED` : seed for random data generation

**Note** in case that you change **username** and **password** of DB, you must change in `mocker.py` too.