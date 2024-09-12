import ast
import symtable

def extract(src):
    symbol_table = symtable.symtable(src, "", "exec")
    for ident in symbol_table.get_identifiers():
        sym = symbol_table.lookup(ident)
        print(sym.get_name())
        print(sym.is_imported())
