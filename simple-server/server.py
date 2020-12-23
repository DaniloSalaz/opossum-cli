from flask import Flask
from markupsafe import escape
from flask import request
from flask import jsonify
from flask import json

app = Flask(__name__)

@app.route('/user/')
def hello_wor():
    data = {"Hola":request.args.get('id')}
    response = app.response_class(
        response=json.dumps(data),
        status=200,
        mimetype='application/json'
    )
    return response