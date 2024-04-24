# Python example using Flask framework
from flask import Flask
import requests
import time

app = Flask(__name__)


# Load configuration
# Code to load configuration containing hosts and ports

@app.route('/ping')
def ping():
    return 'pong\n'


if __name__ == '__main__':
    # Start HTTP server
    app.run(host='0.0.0.0', port=5001, debug=True, threaded=True)
