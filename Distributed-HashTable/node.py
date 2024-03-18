import sys
from flask import Flask, jsonify
import requests

app = Flask(__name__)

class ChordNode:
    def __init__(self, node_id, finger_table_file):
        self.ID = node_id
        self.finger_table = self.load_finger_table(finger_table_file)
        self.lookup_path = []

    def load_finger_table(self, finger_table_file):
        with open(finger_table_file, 'r') as file:
            return [int(line.strip()) for line in file.readlines()]

    def localSuccessorNode(self, key):
        if self.ID < key <= self.finger_table[0]:
            return self.finger_table[0]
        else:
            for i in range(len(self.finger_table) - 1):
                if self.finger_table[i] <= key < self.finger_table[i + 1]:
                    return self.finger_table[i]

    def lookup(self, key):
        key = key % 32  # Modulo arithmetic
        self.lookup_path = []  # Initialize path for this lookup
        successor = self.localSuccessorNode(key)
        if successor == self.ID:
            return self.ID, self.lookup_path  # Key found here
        else:
            self.lookup_path.append(self.ID)  # Add current node to path
            successor_url = f"http://chord-node-service/lookup/{successor}"
            response = requests.get(successor_url, data=str(key))
            # Parse response assuming it contains (responsible_node, path)
            data = response.json()
            responsible_node, path = data['key'], data['path']
            return responsible_node, path + self.lookup_path  # Combine paths

chord_node = None

@app.route('/lookup/<int:key>', methods=['GET'])
def lookup_api(key):
    global chord_node
    try:
        responsible_node, path = chord_node.lookup(key)
        return jsonify({'key': responsible_node, 'path': path})
    except Exception as e:
        return jsonify({'error': str(e)}), 400

def run_server(port=8080):
    global chord_node
    chord_node = ChordNode(int(sys.argv[1]), sys.argv[2])
    app.run(host='0.0.0.0', port=port)

def main():
    if len(sys.argv) != 3:
        print("Error: Incorrect cmd")
        return

    run_server()

if __name__ == "__main__":
    main()
