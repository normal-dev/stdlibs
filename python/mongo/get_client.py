from pymongo import MongoClient
import os

def get_database():
   CONNECTION_STRING = os.environ.get('MONGO_DB_URI', 'mongodb://localhost:27017')
   # Create a connection using MongoClient. You can import MongoClient or use pymongo.MongoClient
   client = MongoClient(CONNECTION_STRING)

   return client
