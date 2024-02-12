import json
import random
import exrex
import os

from pymongo.mongo_client import MongoClient
from pymongo.server_api import ServerApi

# Config this part
CLEAR_PREV = True # clear the collection before inserting new records
N = 500 # number of records to generate
COLLECTION_NAMES = ['user','svcp'] # collection name

# Connect to MongoDB
USERNAME = 'inwza'
PASSWORD = 'strongpassword'
url = f"mongodb://{USERNAME}:{PASSWORD}@localhost:27017/"
#url = "mongodb://localhost:27017"
# Random
SEED = 696969
random.seed(SEED)


client = MongoClient(url, server_api=ServerApi('1'))

try:
    client.admin.command('ping')
    print("Connected successfully!")
except Exception as e:
    print("Unable to connect to the server.")
    print(e)
    exit()

db = client["petpal"]
if CLEAR_PREV:
    for collection_name in COLLECTION_NAMES:
        if collection_name not in db.list_collection_names():
            print(f"Collection {collection_name} not found")
            continue
        db[collection_name].drop()
        print(f"Collection {collection_name} dropped")


print(f"Database : petpal")


def gen_object(id, properties_dict):
    ret = dict()

    for key in properties_dict:
        t = properties_dict[key]['bsonType']
        if t == 'string':
            if 'pattern' in properties_dict[key]:
                limit = 20 if 'maxLength' not in properties_dict[key] else properties_dict[key]['maxLength']
                ret[key] = exrex.getone(properties_dict[key]['pattern'], int(limit))
            else:
                ret[key] = key+str(id)
            if 'minLength' in properties_dict[key] and len(ret[key]) < properties_dict[key]['minLength']:
                ret[key] = ret[key] + "0"*(properties_dict[key]['minLength']-len(ret[key]))
        elif t == 'date':
            ret[key] = f"{random.randint(1, 30)}/{random.randint(1, 12)}/{random.randint(1900, 2021)}"
        elif t == 'int' or t == 'double':
                minimum = 0 if 'minimum' not in properties_dict[key] else properties_dict[key]['minimum']
                maximum = 100 if 'maximum' not in properties_dict[key] else properties_dict[key]['maximum']
                if properties_dict[key]['bsonType'] == 'int':
                    ret[key] = random.randint(minimum, maximum)
                elif properties_dict[key]['bsonType'] == 'double':
                    ret[key] = random.uniform(minimum, maximum)
        elif t == 'bool':
            ret[key] = random.choice([True, False])
        elif t == 'array':
            arr = []
            for i in range(random.randint(1, 6)):
                arr.append(gen_object(1, properties_dict[key]['items']['properties']))
            ret[key] = arr
        else:
            print('error (type='+properties_dict[key]['bsonType']+') not implemented')
    
    return ret

for idx, collection_name in enumerate(COLLECTION_NAMES):
    print(f"Mocking collection {collection_name}")
    collection = db[collection_name]
    
    schema_path = os.path.abspath(f'./{collection_name}_schema.json')
    with open(schema_path) as file:
        print(f"\tReading schema from {schema_path}")
        user_schema = json.load(file)
    
    iter = N if isinstance(N, int) else N[idx]
    print(f"\tGenerating {iter} records -> ", end='')
    for i in range(iter):
        rec = gen_object(i, user_schema['properties'])
        collection.insert_one(rec)
        if i%10 == 0:
            print('.', end='')

    print(f"\n\tDone {collection_name}")