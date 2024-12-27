import importlib
import sys
import inspect
import os
import types

sys.path.append(os.path.abspath(".."))

from mongo.client import get_client

mongo_client = get_client()
mongo_db = mongo_client["apis"]
mongo_coll = mongo_db["python"]

def get_type(value):
  try:
    return typ(value).__name__
  except:
    if isinstance(value, types.FunctionType):
      return "function"
    elif inspect.isclass(type(value)):
      return "class"
    elif isinstance(value, tuple):
      return "tuple"
    elif isinstance(value, list):
      return "list"
    elif isinstance(value, dict):
      return "dict"
    elif isinstance(value, set):
      return "set"
    elif isinstance(value, enumerate):
      return "enumarate"
    elif isinstance(value, str):
      return "string"
    elif isinstance(value, int):
      return "int"
    elif isinstance(value, float):
      return "float"
    elif isinstance(value, complex):
      return "complex"
    elif isinstance(value, bytes):
      return "bytes"
    elif value is None:
      return "none"
    else:
        raise Exception("can't find type")

ns = []
nsn = 0
apisn = 0
for module_name in sys.stdlib_module_names:
  if module_name.startswith("_") or module_name == "antigravity":
    continue
  if importlib.find_loader(module_name) is None:
    continue

  mongo_coll.delete_many({
    "ns" : module_name
  })

  ns.append(module_name)
  nsn += 1

  module = importlib.import_module(module_name)
  for name in dir(module):
    if name.startswith("_"):
      continue

    apisn += 1

    symbol = getattr(module, name)
    typ = get_type(symbol)
    mongo_coll.insert_one({
      "_id" : module_name + "." + name,
      "doc" : "",
      "name" : name,
      "ns" : module_name,
      "type" : typ
    })

mongo_coll.insert_one({
  "_id" : "_cat",
  "n_apis" : apisn,
  "n_ns" : nsn,
  "ns" : ns,
  "version": sys.version.split()[0]
})

mongo_client.close()