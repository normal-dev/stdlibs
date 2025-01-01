import importlib.util
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

apis = []
ns = []
nsn = 0
apisn = 0
for mod_name in sys.stdlib_module_names:
    if mod_name.startswith("_") or mod_name == "antigravity":
        continue
    if importlib.util.find_spec(mod_name) is None:
        continue

    print("module:", mod_name)

    mongo_coll.delete_many({
        "ns" : mod_name
    })

    ns.append(mod_name)
    nsn += 1

    module = importlib.import_module(mod_name)
    for ident in dir(module):
        if ident.startswith("_"):
            continue

        print("ident:", ident)

        apisn += 1

        symbol = getattr(module, ident)
        typ = get_type(symbol)
        apis.append({
            "_id" : mod_name + "." + ident,
            "doc" : "",
            "name" : ident,
            "ns" : mod_name,
            "type" : typ
        })

mongo_coll.insert_many(apis)
mongo_coll.insert_one({
    "_id" : "_cat",
    "n_apis" : apisn,
    "n_ns" : nsn,
    "ns" : ns,
    "version": sys.version.split()[0]
})

mongo_client.close()