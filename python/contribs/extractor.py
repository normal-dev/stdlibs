import symtable
import sys

stdlib = set()
for module_name in sys.stdlib_module_names:
    if module_name.startswith("_") or module_name == "antigravity":
        continue
    if importlib.find_loader(module_name) is None:
        continue

    stdlib.add(module_name)

def find_locus(symbol_table: symtable.symtable, locus: []) -> []:
    for ident in symbol_table.get_identifiers():
        sym = symbol_table.lookup(ident)
        name = sym.get_name()
        if sym.is_imported() and name in stdlib:
            locus.append(name)

    if symbol_table.has_children() == False:
        return

    find_locus(symbol_table=symbol_table.get_children(), locus=locus)

def extract(src):
    symbol_table = symtable.symtable(src, "", "exec")

    locus = []
    find_locus(symbol_table=symbol_table, locus=locus)

    return locus