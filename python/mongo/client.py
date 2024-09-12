from pymongo import MongoClient
import os

def get_client():
   CONNECTION_STRING = os.environ.get('MONGO_DB_URI', 'mongodb://localhost:27017')
   client = MongoClient(CONNECTION_STRING)

   return client
