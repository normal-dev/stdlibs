import sys
import inspect
import os, glob
import types
import pathlib
import subprocess
import tempfile
import shutil
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
        ["ansible", "ansible"],
        ["seleniumbase", "SeleniumBase"],
        ["Significant-Gravitas", "AutoGPT"],
        ["ytdl-org", "youtube-dl"]
    ]:
        repo_name = repo[0]
        repo_owner = repo[1]

        print("repo: {}/{}".format(repo_name, repo_owner))
        repository = gh_client.get_repo(repo_name + "/" + repo_owner)
        repos.append(repository)

    return repos

def save_licenses():
    mongo_coll.delete_one({ "_id": "_licenses"})
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
                "repo": ["pylint-dev", "pylint"],
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
            },
            {
                "author": "Michael Mintz",
                "repo": ["seleniumbase", "SeleniumBase"],
                "type": "MIT license"
            },
            {
                "author": "Toran Bruce Richards",
                "repo": ["Significant-Gravitas", "AutoGPT"],
                "type": "MIT license"
            },
            {
                "author": "Daniel Bolton",
                "repo": ["ytdl-org", "youtube-dl"],
                "type": "Unlicense license"
            }
        ]
    })

def save_cat(contribsn, reposn):
    mongo_coll.delete_one({ "_id": "_cat" })
    mongo_coll.insert_one({
        "_id": "_cat",
        "n_contribs": contribsn,
        "n_repos": reposn
    })

auth = Auth.Token(GITHUB_ACCESS_TOKEN_CONTRIBS)
gh_client = Github(auth=auth)
# Disable retrier
gh_client.default_retry = 1

repos = get_repos(gh_client)
contribsn = 0
repo: Repository.Repository
for repo in repos:
    with tempfile.TemporaryDirectory() as tmpdir:
        repo_owner = repo.owner.login
        repo_name = repo.name

        tmpdir = tmpdir + "/" + repo_owner + "_" + repo_name

        subprocess.run([
            "git",
            "clone",
            "-q",
            "--depth=1",
            "--no-tags",
            "--filter=blob:limit=100k",
            repo.clone_url,
            tmpdir
        ])
        # Remove extraneous
        shutil.rmtree(tmpdir + "/.git")

        # Delete existing contributions
        mongo_coll.delete_many({
            "repo_name": repo_name,
            "repo_owner": repo_owner
        })

        contribs = []
        for root, _, filepaths in os.walk(tmpdir):
            for file_path in filepaths:
                if not file_path.endswith(".py"):
                    continue

                file_name = os.path.basename(file_path)
                try:
                    file_content = pathlib.Path(root + "/" + file_path).read_text()
                except:
                    continue
                print("repo: {}/{}: file: {}".format(repo_owner, repo_name, file_path))

                try:
                    locus = extractor.extract(file_content)
                    print("repo: {}/{}: locus: {}".format(repo_owner, repo_name, len(locus)))
                except Exception as error:
                    print("repo: {}/{}: error: {}".format(repo_owner, repo_name, error))
                    continue

                if len(locus) == 0:
                    continue

                contribs.append({
                    "locus": locus,
                    "code": file_content,
                    "filepath": file_path,
                    "filename": file_name,
                    "repo_name": repo_name,
                    "repo_owner": repo_owner
                })
                contribsn += 1

    mongo_coll.insert_many(contribs)

save_licenses()
save_cat(contribsn=contribsn, reposn=len(repos))

mongo_client.close()
gh_client.close()