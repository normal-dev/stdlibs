import importlib
import sys
import inspect
import os
import types
from github import Github
from github import Auth

sys.path.append(os.path.abspath(".."))

from mongo.client import get_client

GITHUB_ACCESS_TOKEN_CONTRIBS = os.getenv("GITHUB_ACCESS_TOKEN_CONTRIBS")
if GITHUB_ACCESS_TOKEN_CONTRIBS is None:
  raise Exception("missing GitHub access token")

mongo_client = get_client()
mongo_db = mongo_client["contribs"]
mongo_coll = mongo_db["python"]


def get_repos(client):
  repos = []
  for repo in [
    ["PyGithub", "PyGithub"]
  ]:
    repo_name = repo[0]
    repo_owner = repo[1]

    repository = gh_client.get_repo(repo_owner + "/" + repo_name)
    repos.append(repository.raw_data)

auth = Auth.Token(GITHUB_ACCESS_TOKEN_CONTRIBS)
gh_client = Github(auth=auth)

repos = get_repos(gh_client)

def find_py_files(dir):


for repo in repos:
  print(repo)

gh_client.close()