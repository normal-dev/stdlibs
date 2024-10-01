import importlib
import sys
import inspect
import os
import types
import pathlib
from github import Github
from github import Auth
from github import Repository
import extractor

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
    repos.append(repository)

  return repos

auth = Auth.Token(GITHUB_ACCESS_TOKEN_CONTRIBS)
gh_client = Github(auth=auth)

repos = get_repos(gh_client)

repo: Repository.Repository
for repo in repos:
  contents = repo.get_contents("")
  while contents:
      file_content = contents.pop(0)
      if file_content.type == "dir":
          contents.extend(repo.get_contents(file_content.path))
      else:
        filepath = pathlib.Path(file_content.name)
        if filepath.suffix != ".py":
          continue

        f = file_content.decoded_content.decode("utf-8")
        locus = extractor.extract(f)
        print(locus)

gh_client.close()