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
    ["PyGithub", "PyGithub"],
    ["pylint-dev", "astroid"],
    ["pylint-dev", "pylint"],
    ["fastapi", "fastapi"],
    ["tensorflow", "models"],
    ["openai", "whisper"],
    ["ansible", "ansible"]
  ]:
    repo_name = repo[0]
    repo_owner = repo[1]

    print("repo: {}/{}".format(repo_name, repo_owner))
    repository = gh_client.get_repo(repo_owner + "/" + repo_name)
    repos.append(repository)

  return repos

def save_licenses():
  mongo_coll.insert_one({
    "_id": "_licenses",
    "repos": [
      {
        "author": "Vincent Jacques",
        "repo": ["PyGithub", "PyGithub"],
        "type": "GNU General Public License v3.0"
      },
      {
        "author": "Logilab, and astroid contributors",
        "repo": ["pylint-dev", "astroid"],
        "type": "LGPL-2.1 license"
      },
      {
        "author": "Logilab and Pylint contributors",
        "repo": ["pylint-dev", "astroid"],
        "type": "GPL-2.0 license"
      },
      {
        "author": "Sebastián Ramírez",
        "repo": ["fastapi", "fastapi"],
        "type": "MIT license"
      },
      {
        "author": "Google LLC.",
        "repo": ["tensorflow", "models"],
        "type": "Apache License Version 2.0"
      },
      {
        "author": "OpenAI",
        "repo": ["openai", "whisper"],
        "type": "MIT license"
      },
      {
        "author": "Red Hat, Inc.",
        "repo": ["ansible", "ansible"],
        "type": "GPL-3.0 license"
      }
    ]
  })

auth = Auth.Token(GITHUB_ACCESS_TOKEN_CONTRIBS)
gh_client = Github(auth=auth)

repos = get_repos(gh_client)

repo: Repository.Repository
for repo in repos:

  mongo_coll.delete_many({
    "repo_name": repo.name,
    "repo_owner": repo.owner.login
  })

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
        print("file: {}".format(filepath))
        try:
          locus = extractor.extract(f)
          print("locus:", len(locus))
        except Exception as error:
          print("error:", error)
          continue

        if len(locus) is 0:
          continue

        mongo_coll.insert_one({
          "locus": locus,
          "code": f,
          "filepath": file_content.path,
          "filename": file_content.name,
          "repo_name": repo.name,
          "repo_owner": repo.owner.login
        })

save_licenses()

mongo_client.close()
gh_client.close()