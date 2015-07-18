from flask import Flask
import json
import random

app = Flask(__name__)

SHAPES = [
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [True, False, True, True, True, False]
    }),
    json.dumps({
        "Width": 1,
        "Position": [0, 5],
        "Data": [True, True, True, True]
    }),
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [True, True, True, True]
    }),
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [False, True, False, True, True, True]
    }),
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [True, False, True, False, True, True]
    }),
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [True, False, True, True, False, True]
    }),
    json.dumps({
        "Width": 2,
        "Position": [0, 4],
        "Data": [False, True, True, True, True, False]
    })
]


@app.route("/shape", methods=["GET"])
def generate_shape():
    return random.choice(SHAPES)

if __name__ == "__main__":
    app.run(port=8001)
