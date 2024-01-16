from flask import Flask
from flask_graphql import GraphQLView
from mongoengine import connect
import os

from schema import schema

DATABASE = 'flask-mongodb-graphene'
PASSWORD = os.environ.get("MONGODB_PASSWORD")

client = connect(DATABASE, host=f'mongodb+srv://mongograph:{PASSWORD}@clusterjj-gazky.mongodb.net/?ssl=true&ssl_cert_reqs=CERT_NONE', alias='default')

app = Flask(__name__)
app.debug = True

app.add_url_rule('/graphql', view_func=GraphQLView.as_view('graphql', schema=schema, graphiql=True))

if __name__ == '__main__':
    app.run(port=5002)