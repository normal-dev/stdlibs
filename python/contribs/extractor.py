import sys
import ast
import importlib
from astroid import parse
from astroid import nodes
from astroid import helpers
from astroid import util
import astroid

stdlib = set()
for mod_name in sys.stdlib_module_names:
    if mod_name.startswith("_") or mod_name == "antigravity":
        continue
    if importlib.util.find_spec(mod_name) is None:
        continue

    stdlib.add(mod_name)

class ImportVisitor():
    def __init__(self, tree: nodes.Module, locus: []):
        self.tree = tree
        self.locus = locus

    def visit(self, node: astroid.node_classes.NodeNG):
        func_name = "visit_" + node.__class__.__name__
        func = getattr(self, func_name, self.visit_generic)
        return func(node)

    def visit_generic(self, node: astroid.NodeNG) -> None:
        for child in node.get_children():
            self.visit(child)

    def visit_Import(self, import_node: nodes.Import):
        # If no import name is part of stdlib, return
        if not any([n[0].split(".")[0] in stdlib for n in import_node.names]):
            return None

        find_locus(import_node=import_node, tree=self.tree, locus=self.locus)

    def visit_ImportFrom(self, import_node: nodes.ImportFrom):
        # If import name is not part of stdlib, return
        mod_name = import_node.modname.split(".")[0]
        if mod_name not in stdlib:
            return

        find_locus(import_node=import_node, tree=self.tree, locus=self.locus)

# Returns the first level import with the name nodes name
def resolve_qual_import(import_name: str, name_node: nodes.Name, from_import=False):
    # "datetime.time()"
    if "." in import_name:
        spl = import_name.split(".", 1)
        return spl[0] + "." + spl[1]

    return import_name + "." + name_node.name

# "import sys.stdlib_module_names as mod_names, ast.parse, datetime as d"
def resolve_import(import_node: nodes.Import, name_node: nodes.Name):
    if isinstance(name_node.parent, nodes.Attribute):
        mod_name = import_node.real_name(name_node.name)
        return mod_name + "." + name_node.parent.attrname

    "(sys.stdlib_module_namen, mod_names), (ast.parse, None)"
    for names in import_node.names:
        match names:
            # "(sys, None) | (sys, s)"
            case (x, None) | (x, name_node.name):
                mod_name = names[0].split(".")[0]
                if mod_name not in stdlib:
                    continue

                return resolve_qual_import(names[0], name_node)

    return None

# "from typing import Any, ClassVar, Union"
def resolve_import_from(import_node: nodes.ImportFrom, name_node: nodes.Name):
    return resolve_qual_import(import_node.modname, name_node)

# Search for import and name node match
def find_locus(import_node: nodes.NodeNG, tree: nodes.Module, locus: []):
    name_node: nodes.Name
    # For every "ast.Name"
    for name_node in tree.nodes_of_class(nodes.Name):
        # Get scope of frame of name
        scope = name_node.frame().lookup(name_node.name)
        for n in scope[1]:
            if n != import_node:
                continue

            # "from foo.bar import baz"
            if isinstance(import_node, nodes.ImportFrom):
                ident = resolve_import_from(import_node, name_node)
                locus.append({
                    "ident": ident,
                    "line": name_node.lineno
                })

            # "import foo"
            if isinstance(import_node, nodes.Import):
                ident = resolve_import(import_node, name_node)
                if ident is None:
                    continue
                locus.append({
                    "ident": ident,
                    "line": name_node.lineno
                })

def extract(src):
    locus = []
    tree = parse(src)
    ImportVisitor(tree, locus).visit(tree)

    return locus