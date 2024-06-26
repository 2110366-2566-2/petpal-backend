# import libs ----------------------------------------------------------------
import requests
import json
import random
import datetime
import time

from pymongo.mongo_client import MongoClient
from pymongo.server_api import ServerApi

import auth

# connect to MongoDB and drop database ---------------------------------------
USERNAME = 'inwza'
PASSWORD = 'strongpassword'
db_url = f"mongodb://{USERNAME}:{PASSWORD}@localhost:27017/"

client = MongoClient(db_url, server_api=ServerApi('1'))
client.drop_database('petpal')
print('Dropped database petpal')


# set up const ---------------------------------------------------------------
BASE_URL = "http://localhost:8080/"

N_USER = 5
N_PETS = 5

N_SVCP = 5
N_SERVICES = 5

N_ADMIN = 1

USER0_ID = None
SVCP0_ID = None

SEED = 696969
random.seed(SEED)


# register svcp --------------------------------------------------------------
for svcp_idx in range(N_SVCP):
    svcp = {
        "SVCPEmail": f"{svcp_idx}@svcp.com",
        "SVCPPassword": "password",
        "SVCPServiceType": "whatever",
        "SVCPUsername": f"svcp{svcp_idx}"
    }
    response = requests.post(BASE_URL + "register-svcp", json=svcp)
    if response.status_code != 200:
        continue
    else:
        print('registered svcp', svcp_idx)


# create service -------------------------------------------------------------
for svcp_idx in range(N_SVCP):
    # log in
    svcp_token = auth.login(svcp_idx, "svcp")
    if svcp_token == None:
        continue

    # create service
    for service_idx in range(N_SERVICES):
        service = {
            "price": random.randint(100, 1000),
            "serviceDescription": f"this is {service_idx}th service of svcp {svcp_idx}",
            "serviceName": "washing "+random.choice(["dog", "cat", "rabbit"]),
            "serviceType": random.choice(["Healthcare", "Grooming", "Pet walking", "Others"]),
            "timeslots": [
                {
                "endTime": (datetime.datetime.now() + datetime.timedelta(hours=random.randint(1, 10))).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "startTime": (datetime.datetime.now() + datetime.timedelta(hours=-1*random.randint(1, 10))).strftime("%Y-%m-%dT%H:%M:%SZ")
                }
            for _ in range(random.randint(1, 3))]
        }
        response = requests.post(BASE_URL + "service/create", json=service, cookies={"token": svcp_token})
        if response.status_code != 200:
            continue
        else:
            print('\tcreated service', service_idx, 'of svcp', svcp_idx)


# set default bank account ---------------------------------------------------
for svcp_idx in range(N_SVCP):
    # log in
    svcp_token = auth.login(svcp_idx, "svcp")
    if svcp_token == None:
        continue

    default_bank_account = {
        "defaultBankAccountNumber": "1234567890",
        "defaultBank": "KTB",
    }
    response = requests.post(BASE_URL + "serviceproviders/set-default-bank-account", json=default_bank_account, cookies={"token": svcp_token})
    if response.status_code != 200:
        continue
    else:
        print('\tset default bank account of svcp', svcp_idx)


# upload description ----------------------------------------------------------
for svcp_idx in range(N_SVCP):
    # log in
    svcp_token = auth.login(svcp_idx, "svcp")
    if svcp_token == None:
        continue

    description = {
        "description": "this is description of svcp "+str(svcp_idx)
    }
    response = requests.post(BASE_URL + "serviceproviders/upload-description", json=description, cookies={"token": svcp_token})
    if response.status_code != 200:
        continue
    else:
        print('\tuploaded description of svcp', svcp_idx)


# register user --------------------------------------------------------------
for user_idx in range(N_USER):
    user = {
        "address": f"homeland 101 room @{user_idx}",
        "dateOfBirth": (datetime.datetime.now() - datetime.timedelta(days=365*random.randint(20, 50))).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "email": f"{user_idx}@user.com",
        "fullName": f"{user_idx} user",
        "password": "password",
        "phoneNumber": "0999999999",
        "username": f"user{user_idx}"
    }
    response = requests.post(BASE_URL + "register-user", json=user)
    if response.status_code != 200:
        continue
    else:
        print('registered user', user_idx)


# add pets -------------------------------------------------------------------
for user_idx in range(N_USER):
    # log in
    user_token = auth.login(user_idx, "user")
    if user_token == None:
        continue

    for pet_idx in range(N_PETS):
        pet = {
            "age": random.randint(1, 20),
            "behaviouralNotes": random.choice(["good", "bad", "normal", "crazy"]),
            "breed": random.choice(["dog", "cat", "rabbit"]),
            "certificate": random.choice(["yes", "no"]),
            "dietyPreferences": random.choice(["meat", "vege", "both"]),
            "gender": random.choice(["male", "femail"]),
            "healthInformation": "string",
            "name": f"pet{pet_idx} of user{user_idx}",
            "type": random.choice(["dog", "cat", "rabbit"]),
            "vaccinations": random.choice(["yes", "no"]),
        }
        response = requests.post(BASE_URL + "user/pets", json=pet, cookies={"token": user_token})
        if response.status_code != 200:
            continue
        else:
            print('\tadded pet', pet_idx, 'of user', user_idx)

# create booking -------------------------------------------------------------
for user_idx in range(min(N_USER, 2)):
    # get svcp_id, service_id, timeslot_id
    svcp_token = auth.login(user_idx, "svcp")
    if svcp_token == None:
        continue
    svcp_entity = auth.current_entity(svcp_token)
    if svcp_entity == None:
        continue

    svcp_id = svcp_entity['SVCPID']
    for service in svcp_entity['services']:
        service_id = service['serviceID']
        timeslot_id = service['timeslots'][0]['timeslotID']
        
        # log in
        user_token = auth.login(user_idx, "user")
        if user_token == None:
            continue

        booking = {
            "SVCPID": svcp_id,
            "serviceID": service_id,
            "timeslotID": timeslot_id,
        }
        response = requests.post(BASE_URL + "service/booking/create", json=booking, cookies={"token": user_token})
        if response.status_code != 201:
            continue
        else:
            print('\tcreated booking of user', user_idx)
        
        booking_id = response.json()['result']['bookingID']
        
        # sent service issue -------------------------------------------------
        issue = {
            "details": "this is issue of booking",
            "issueType": "service",
            "associatedBookingID": booking_id
        }
        response = requests.post(BASE_URL + "issue", cookies={"token": user_token}, files=issue, data=issue)
        if response.status_code != 200:
            continue
        else:
            print('\tsent issue of booking of user', user_idx)

        # sent refund issue -------------------------------------------------
        if random.random() < 0.5:
            issue = {
                "details": "this is issue of refund",
                "issueType": "refund",
                "associatedBookingID": booking_id
            }
            response = requests.post(BASE_URL + "issue", cookies={"token": user_token}, files=issue, data=issue)
            if response.status_code != 200:
                continue
            else:
                print('\tsent issue of refund of user', user_idx)

# register admin -------------------------------------------------------------
for admin_idx in range(N_ADMIN):
    admin = {
        "email": f"{admin_idx}@admin.com",
        "fullName": f"{admin_idx} admin",
        "password": "password",
        "username": f"admin{admin_idx}"
    }
    response = requests.post(BASE_URL + "register-admin", json=admin)
    if response.status_code != 200:
        continue
    else:
        print('registered admin', admin_idx)

# add system issue -----------------------------------------------------------
for user_idx in range(N_USER):
    # log in
    user_token = auth.login(user_idx, "user")
    if user_token == None:
        continue

    issue = {
        "details": "this is issue of system",
        "issueType": "system",
    }
    response = requests.post(BASE_URL + "issue", cookies={"token": user_token}, files=issue, data=issue)
    if response.status_code != 200:
        continue
    else:
        print('added system issue of user', user_idx)

svcp_token = auth.login(0, "svcp")
if svcp_token is not None:
    svcp_entity = auth.current_entity(svcp_token)
    SVCP0_ID = svcp_entity['SVCPID']

user_token = auth.login(0, "user")
if user_token is not None:
    user_entity = auth.current_entity(user_token)
    USER0_ID = user_entity['id']

booking = client.petpal.booking.find({"userID": USER0_ID, "SVCPID": SVCP0_ID})

# able to see complete and feedback
target_date = {
    'status.paymentStatus': True,
    'status.svcpConfirmed': True,
    'status.userCompleted': True,
    'status.svcpCompleted': True,
    'startTime': (datetime.datetime.now() + datetime.timedelta(days=-10, hours=-1)).strftime("%Y-%m-%dT%H:%M:%SZ"),
    'endTime': (datetime.datetime.now() + datetime.timedelta(days=-10)).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
client.petpal.booking.update_one({"_id": booking[0]['_id']}, {"$set": target_date})

# able to refund
target_date = {
    'status.paymentStatus': True,
    'status.svcpCompleted': True,
    'status.userCompletedTimestamp': (datetime.datetime.now() + datetime.timedelta(days=-1)).strftime("%Y-%m-%dT%H:%M:%SZ"),
    'status.svcpCompletedTimestamp': (datetime.datetime.now() + datetime.timedelta(days=-1)).strftime("%Y-%m-%dT%H:%M:%SZ"),
    'startTime': (datetime.datetime.now() + datetime.timedelta(days=-2, hours=-1)).strftime("%Y-%m-%dT%H:%M:%SZ"),
    'endTime': (datetime.datetime.now() + datetime.timedelta(days=-2)).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
client.petpal.booking.update_one({"_id": booking[1]['_id']}, {"$set": target_date})

client.close()