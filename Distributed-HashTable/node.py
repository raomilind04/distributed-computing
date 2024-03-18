import sys
from flask import Flask, jsonify
import requests

app = Flask(__name__)

class ChordNode:
    def __init__(self, node_id, pred_id, finger_table_file):
        self.ID = int(node_id)
        self.pred = int(pred_id)
        self.finger_table = self.load_finger_table(finger_table_file)
        print("ID:", self.ID, "Pred:", self.pred, "Finger Table:", self.finger_table)

    def load_finger_table(self, finger_table_file):
        with open(finger_table_file, 'r') as file:
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
    
    def lookup(self, key):
        key = key % 32
        self.lookup_path = []
        successor = self.localSuccessorNode(key)
        if successor == self.ID:
            return self.ID, self.lookup_path 
        else:
            self.lookup_path.append(self.ID) 
            print("Sending to next node: ", successor)
            successor_url = f"http://chord-node-service/lookup/{successor}"
            response = requests.get(successor_url, data=str(key))
            data = response.json()
            responsible_node, path = data['key'], data['path']
            return responsible_node, path + self.lookup_path

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
    chord_node = ChordNode(int(sys.argv[1]), int(sys.argv[2]), sys.argv[3])
    app.run(host='0.0.0.0', port=port)

def main():
    if len(sys.argv) != 4:
        print("Error: Incorrect cmd")
        return

    run_server()

if __name__ == "__main__":
    main()
