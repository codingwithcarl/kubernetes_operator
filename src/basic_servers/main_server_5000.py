import requests
import json
import threading
from flask import Flask, Response

app = Flask(__name__)

# Configuration containing hosts and ports
f = open('config.json')
config = json.load(f)


def ping_other_hosts():
    for host in config["hosts"]:
        url = f"http://{host['host']}:{host['port']}/ping"
        try:
            response = requests.get(url)
            print(f"Ping to {url}: {response.text.strip()}")
        except Exception as e:
            print(f"Failed to ping {url}: {str(e)}")

    # Schedule the next ping after interval seconds
    threading.Timer(config["interval_seconds"], ping_other_hosts).start()


# Start pinging other hosts in background thread
ping_other_hosts()


@app.route('/ping')
def ping():
    return 'pong\n'


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True, threaded=True)
