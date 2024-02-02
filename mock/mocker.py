import json
import random
import exrex
import os

from pymongo.mongo_client import MongoClient
from pymongo.server_api import ServerApi

# Config this part
CLEAR_PREV = True # clear the collection before inserting new records
N = 1000 # number of records to generate
COLLECTION_NAME = 'user' # collection name

# Connect to MongoDB
USERNAME = 'inwza'
PASSWORD = 'strongpassword'
url = f"mongodb://{USERNAME}:{PASSWORD}@localhost:27017/"

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
    db[COLLECTION_NAME].drop()
    print(f"Collection {COLLECTION_NAME} dropped")

user_collection = db[COLLECTION_NAME]

print(f"Database : petpal\nCollection : {COLLECTION_NAME}")

# read schema json file
schema_path = os.path.abspath(f'./mock/{COLLECTION_NAME}_schema.json')
with open(schema_path) as file:
    print(f"Reading schema from {schema_path}")
    user_schema = json.load(file)

def gen_object(id, properties_dict):
    ret = dict()

    for key in properties_dict:
        match properties_dict[key]['bsonType']:
            case 'string':
                if 'date' in key: # date special case
                    ret[key] = f"{random.randint(1, 30)}/{random.randint(1, 12)}/{random.randint(1900, 2021)}"
                else: # normal string
                    if 'pattern' in properties_dict[key]:
                        limit = 20 if 'maxLength' not in properties_dict[key] else properties_dict[key]['maxLength']
                        ret[key] = exrex.getone(properties_dict[key]['pattern'], int(limit))
                    else:
                        ret[key] = key+str(id)

                    if 'minLength' in properties_dict[key] and len(ret[key]) < properties_dict[key]['minLength']:
                        ret[key] = ret[key] + "0"*(properties_dict[key]['minLength']-len(ret[key]))
            case 'int':
                minimum = 0 if 'minimum' not in properties_dict[key] else properties_dict[key]['minimum']
                maximum = 100 if 'maximum' not in properties_dict[key] else properties_dict[key]['maximum']
                ret[key] = random.randint(minimum, maximum)
            case 'array':
                arr = []
                for i in range(random.randint(1, 6)):
                    arr.append(gen_object(1, properties_dict[key]['items']['properties']))
                ret[key] = arr
            case _:
                print('error : '+properties_dict[key]['bsonType'])
    
    return ret

print(f"Generating {N} records")
for i in range(N):
    rec = gen_object(i, user_schema['properties'])
    user_collection.insert_one(rec)
    if i%10 == 0:
        print('.', end='')

print("\nDone all")