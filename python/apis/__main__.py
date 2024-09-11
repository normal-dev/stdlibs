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

def is_capsule(o):
    t = type(o)
    return t.__module__ == 'builtins' and t.__name__ == 'PyCapsule'

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
    typ = get_type(symbol)
    mongo_coll.insert_one({
      "_id" : module_name + "." + name,
      "doc" : "",
      "name" : name,
      "ns" : module_name,
      "type" : typ
    })

