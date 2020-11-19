import csv
import json
import requests

def parse_csv(csv_file):
    json_arr = []
    # open csv
    with open(csv_file, encoding='utf-8') as f:
        # parse csv
        csv_reader = csv.DictReader(f)
        # convert csv rows to dict
        for row in csv_reader:
            json_arr.append(row)

    return json.dumps(json_arr, indent=2)

def send(uname, pwd, csv_file, url):
    json_str = '{\n"uname": "%s",\n"pwd": "%s",\n "data": ' %(uname, pwd)
    json_str += parse_csv(csv_file)
    json_str += '\n}'
    resp = requests.post(url, data=json_str)
    return resp
