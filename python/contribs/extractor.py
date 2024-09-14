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
    def __init__(self, tree: nodes.Module):
        self.tree = tree

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
            if alias is not "":
                module_ident = alias

            find_locus(module_ident=module_ident, module=self.tree)

def find_locus(module_ident: str, module: nodes.Module):
    node: nodes.Name
    for node in module.nodes_of_class(nodes.Name):
        frame = node.frame()
        inf = node.inferred()
        context = frame.lookup(node.name)
        for n in context[1]:
            if isinstance(n, nodes.Import):
                print(node.lineno)

def extract(code):
    locus = []
    tree = parse(code)

    ImportVisitor(tree).visit(tree)