import sys
import ast
import importlib
from astroid import parse
from astroid import nodes
from astroid import helpers
from astroid import util
import astroid

stdlib = set()
for module_name in sys.stdlib_module_names:
    if module_name.startswith("_") or module_name == "antigravity":
        continue
    if importlib.find_loader(module_name) is None:
        continue

    stdlib.add(module_name)

class ImportVisitor():
    def __init__(self, tree: nodes.Module, locus: list[dict]):
        self.tree = tree
        self.locus = locus

    def visit(self, node: astroid.node_classes.NodeNG):
        func_name = "visit_" + node.__class__.__name__
        func = getattr(self, func_name, self.visit_generic)
        return func(node)

    def visit_generic(self, node: astroid.NodeNG) -> None:
        for child in node.get_children():
            self.visit(child)

    def visit_Import(self, node: nodes.Import):
        for n in node.names:
            name = n[0]
            alias = n[1]
            module_ident = name
            if alias != None:
                module_ident = alias

            find_locus(import_node=node, tree=self.tree, locus=self.locus)

def find_locus(import_node: nodes.Import, tree: nodes.Module, locus: list[dict]):
    name_node: nodes.Name
    for name_node in tree.nodes_of_class(nodes.Name):
        frame = name_node.frame()
        scope = frame.lookup(name_node.name)
        for node in scope[1]:
            if node != import_node:
                continue

            locus.append({
                "ident": import_node + "." + name_node,
                "line": name_node.lineno
            })

def extract(code):
    locus = list[dict]
    tree = parse(code)

    ImportVisitor(tree, locus).visit(tree)