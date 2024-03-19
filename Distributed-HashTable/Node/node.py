#!/bin/python3
from flask import Flask, jsonify, render_template, request, redirect
from flask_cors import CORS 
import os
import logging

logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

class ChordNode:
    def __init__(self, node_id, pred_id, finger_table_file):
        self.ID = int(node_id)
        self.pred = int(pred_id)
        self.finger_table = self.load_finger_table(finger_table_file)
        logger.info("ID:", self.ID, "Pred:", self.pred, "Finger Table:", self.finger_table)

    def load_finger_table(self, finger_table_file):
        finger_table_file_path = os.path.join("/data", finger_table_file)
        with open(finger_table_file_path, 'r') as file:
            return [int(line.strip()) for line in file.readlines()]

    def localSuccessorNode(self, key):
        m = len(self.finger_table)
        if self.pred < key <= self.ID:
            return self.ID
        elif self.ID < key <= self.finger_table[0]:
            return self.finger_table[0]
        else:
            for i in range(0, m):
                if self.finger_table[i % m] <= key < self.finger_table[(i+1) % m] or self.finger_table[i] == key:
                    return self.finger_table[i % m]
            return self.finger_table[-1]
    
    def lookup(self, key, path=None):
        key = key % 32
        if path is None:
            path = []
        successor = self.localSuccessorNode(key)
        path.append(self.ID)
        if successor == self.ID:
            return self.ID, path
        else:
            return successor, path

chord_node = None

@app.route('/lookup/<int:key>', methods=['GET'])
def lookup_api(key):
    global chord_node
    
    try:
        if request.args.get('path'):
            path = request.args.get('path').split(',')
            path = list(map(int, path))
        else:
            path = None
    
        successor, path = chord_node.lookup(key, path)
        if successor == chord_node.ID:
            formatted_path = " -> ".join(map(str, path))
            return jsonify({'key': successor, 'path': formatted_path})
        else:
            target_service_name = f"chord-node-service-{successor}"
            successor_address = f"http://{target_service_name}.default.svc.cluster.local:8080/lookup/{key}"
            path_str = ','.join(map(str, path))
            redirect_url = f"{successor_address}?path={path_str}"
            logger.info(f"Redirect to {redirect_url}")
            return redirect(redirect_url)
        
    except Exception as e:
        logger.error(e)
        return jsonify({'error': str(e)}), 400

def run_server(port=8080):
    global chord_node
    chord_node = ChordNode(int(os.environ['NODE_ID']), int(os.environ['PREV_NODE_ID']), os.environ['FINGER_TABLE_FILE'])
    app.run(host='0.0.0.0', port=port)

def main():
    run_server()
# docker buildx build --platform=linux/amd64 -t raomilind04/chord-node-app:latest-amd64 .
if __name__ == "__main__":
    main()