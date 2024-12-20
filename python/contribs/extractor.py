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
        find_locus(import_node=import_node, tree=self.tree, locus=self.locus)

    def visit_ImportFrom(self, import_node: nodes.Import):
        find_locus(import_node=import_node, tree=self.tree, locus=self.locus)

def resolve_attr(node: nodes.Attribute):
    return node.attrname

def resolve_import_from(import_node: nodes.ImportFrom, modname: str):
    root = import_node.root()
    if isinstance(root, nodes.Module):
        return root.relative_to_absolute_name(
            modname, level=import_node.level
        )

    return None

def find_locus(import_node: nodes.NodeNG, tree: nodes.Module, locus: []):
    name_node: nodes.Name
    for name_node in tree.nodes_of_class(nodes.Name):
        frame = name_node.frame()
        scope = frame.lookup(name_node.name)
        for node in scope[1]:
            if node != import_node:
                continue

            # timedelta()
            if isinstance(import_node, nodes.ImportFrom):
                mod = resolve_import_from(import_node, import_node.modname)
                locus.append({
                    "ident": mod + "." + name_node.name,
                    "line": name_node.lineno
                })

            # types.CodeType
            if isinstance(name_node.parent, nodes.Attribute) and isinstance(import_node, nodes.Import):
                ident = resolve_attr(name_node.parent)
                locus.append({
                    "ident": import_node.real_name(name_node.name) + "." + ident,
                    "line": name_node.lineno
                })

def extract(src):
    locus = []
    tree = parse(src)
    ImportVisitor(tree, locus).visit(tree)

    return locus