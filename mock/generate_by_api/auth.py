import requests

def login(idx,type):
    svcp = {
        "email": f"{idx}@{type}.com",
        "password": "password",
        "loginType": type
    }
    response = requests.post("http://localhost:8080/login", json=svcp)
    if response.status_code != 200:
        return None 
    else:
        print('logged in', type, idx)
        token = response.json()['AccessToken']
        return token

def current_entity(token):
    response = requests.get("http://localhost:8080/current-entity", cookies={"token": token})
    if response.status_code != 202:
        return None
    else:
        return response.json()
