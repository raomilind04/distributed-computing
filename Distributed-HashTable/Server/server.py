#!/bin/python3
from flask import Flask, jsonify, render_template, request, redirect
from flask_cors import CORS
import os
import requests
import logging


logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

@app.route('/lookup/<int:key>/<int:starting_node>', methods=['GET'])
def lookup_api(key, starting_node): 
    try:
        starting_service_name = f"chord-node-service-{starting_node}" 
        starting_address = f"http://{starting_service_name}.default.svc.cluster.local:8080/lookup/{key}"
        logger.info(f"starting_address: {starting_address}")

        response = requests.get(starting_address, allow_redirects=True)

        if response.status_code == 200:
            return response.json()
        else:
            logger.error(f"Failed to get response from {starting_address}")
            return jsonify({'error': f"Failed to get response from {starting_address}"}), 500

    except Exception as e:
        logger.error(e)
        return jsonify({'error': str(e)}), 400
    
@app.route('/')
def index():
    return render_template('index.html')

def main():
    app.run(host='0.0.0.0', port=8080)
# docker buildx build --platform=linux/amd64 -t raomilind04/chord-node-app:server-latest-amd64 .
if __name__ == '__main__':
    main()
