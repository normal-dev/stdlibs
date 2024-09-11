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
    return type(value).__name__
  except:
    if inspect.isbuiltin(value) or inspect.ismethod(value) or isinstance(value, types.FunctionType):
      return "function"
    elif inspect.isclass(value):
      return "class"
    elif inspect.ismodule(value):
      return "module"
    elif isinstance(value, str):
      return "string"
    elif isinstance(value, int):
      return "int"
    elif isinstance(value, float):
      return "float"
    elif isinstance(value, list):
      return "list"
    elif isinstance(value, dict):
      return "dict"
    elif value is None:
      return "none"
    else:
        print(type(value))

mongo_coll.delete_many({})

for module_name in sys.stdlib_module_names:
  if module_name.startswith("_") or module_name is "antigravity":
    continue

  if importlib.find_loader(module_name) is None:
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

