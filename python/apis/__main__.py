import pkgutil
import importlib
import inspect
import sys

def get_symbols(module):
    symbols = []
    for name in dir(module):
        try:
            # TODO: Get symbol type
            symbol = getattr(module, name)
            if callable(symbol) or inspect.ismodule(symbol):
                symbols.append(name)
        except AttributeError:
            pass  # Some members might raise AttributeError on access

    return symbols

def get_stdlib():
    stdlib = {}
    for module_name in sys.stdlib_module_names:
        if module_name.startswith('_'):
          continue
        if module_name == 'antigravity':
          continue

        # TEST
        if module_name != 'math':
          continue
        # TEST

        module = importlib.import_module(module_name)
        for name in dir(module):
          if name.startswith('_'):
            continue

          symbol = getattr(module, name)
          print(name, type(symbol).__name__)

    # mods = pkgutil.iter_modules()
    # for finder, module_name, ispkg in mods:
    #     # TEST
    #     if module_name != 'math':
    #       continue
    #     # TEST

    #     if module_name.startswith('_'):
    #       continue

    #     if module_name == 'antigravity':
    #       continue

    #     # TODO: Skip antrigravity
    #     try:
    #         module_name = importlib.import_module(module_name)
    #         if hasattr(module_name, '__file__') and 'site-packages' not in module_name.__file__:
    #             symbols = get_symbols(module_name)

    #             if symbols:
    #                 stdlib[module_name] = symbols
    #     except (ImportError, AttributeError):
    #         pass  # Handle modules that can't be imported or inspected

    return stdlib

if __name__ == "__main__":
    stdlib = get_stdlib()
    for module, symbols in stdlib.items():
        print(f"Module: {module}")
        print(f"Members: {symbols}\n")
