import csv
import json
import requests

def parse_csv(csv_file):
    data = []
    # open csv
    with open(csv_file, encoding='utf-8') as f:
        # parse csv
        csv_reader = csv.DictReader(f)
        # convert csv rows to dict
        for row in csv_reader:
            data.append(row)

    return json.dumps(data, indent=2)

def send(uname, pwd, csv_file, url):
    data = parse_csv(csv_file)
    json = '{\n\t"uname": "%s",\n\t"pwd": "%s",\n\t"data": ' %(uname, pwd) + data + '}'
    resp = requests.post(url, data=json)
    return resp

if __name__ == "__main__":
    import secret
    import sys

    filename ="data/accel_data_running.csv"

    if '-local' in sys.argv:
        req_url = r'http://localhost:32321/wearables/dashboard'
    else:
        req_url = r'https://ethanvogelsang.com/wearables/dashboard'

    usr = secret.uname
    pwd = secret.pwd
    send(usr, pwd, filename, req_url)
