#!/usr/bin/env python
import requests
import json
import argparse
import textwrap
from datetime import datetime, timedelta


def main():
    end_date = datetime.now()
    start_date = end_date - timedelta(days=7)  # 過去7日間のCVEを取得
    api = (
        f"https://services.nvd.nist.gov/rest/json/cves/2.0?resultsPerPage=100&pubStartDate={start_date.isoformat()}Z&pubEndDate={end_date.isoformat()}Z"
    )
    # api = 'https://services.nvd.nist.gov/rest/json/cves/2.0?resultsPerPage=10&startIndex=0'
    response = requests.get(api)
    json_data = json.loads(response.text)

    vulnerabilities = json_data['vulnerabilities']
    print(vulnerabilities)
    descriptions = [item['cve']['descriptions'][0]['value'] for item in json_data['vulnerabilities']]
    ids = [item['cve']['id'] for item in json_data['vulnerabilities']]

    print(descriptions)
    print(ids)

main()
