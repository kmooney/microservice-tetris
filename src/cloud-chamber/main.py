from flask import Flask, request
app = Flask(__name__)


@app.route("/", methods=["GET", "POST", "PUT"])
def index():
    print "data start============================"
    print request.data
    print "data end=============================="
    return "OK"


if __name__ == "__main__":
    app.run()
