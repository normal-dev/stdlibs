import importlib
import sys
import inspect
import os

sys.path.append(os.path.abspath(".."))

from mongo.client import get_client

mongo_client = get_client()
mongo_db = mongo_client["apis"]
mongo_coll = mongo_db["python"]

def get_type(value):
  try:
    return type(value).__name__
  except:
    if inspect.isbuiltin(value) or inspect.ismethod(value):
      return "function"
    elif inspect.isclass(value):
      return "class"
    elif inspect.ismodule(value):
      return "module"
    else:
      return type(value)

mongo_coll.delete_many({})

for module_name in sys.stdlib_module_names:
  if module_name.startswith("_") or module_name == "antigravity":
    continue

  if module_name != "math":
    continue

  module = importlib.import_module(module_name)
  for name in dir(module):
    if name.startswith("_"):
      continue

    symbol = getattr(module, name)
    type = get_type(symbol)
    mongo_coll.insert_one({
      "_id" : module_name + "." + name,
      "doc" : "",
      "name" : name,
      "ns" : module_name,
      "type" : type
    })

