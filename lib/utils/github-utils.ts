import * as path from "path";
import * as fs from "fs";
import {spawnSync} from "child_process";
import {GithubRepository} from "../config/types";


export const cloneOrPullRepoSync = (repoInfo: GithubRepository, localRepoPath: string) => {
    const repoUrl = getGithubRepoURL(repoInfo.owner, repoInfo.name);

    if (fs.existsSync(localRepoPath)) {
        console.log(`Path ${path} already exists`);
    } else {
        fs.mkdirSync(localRepoPath, {recursive: true});
    }

    try {
        spawnSync("git", ["clone", repoUrl, localRepoPath, '--branch', repoInfo.branch], {stdio: 'inherit'});
        console.log(`Cloned ${repoUrl} to ${localRepoPath}. Branch: ${repoInfo.branch}`);
    } catch (err: unknown) {
        if (err instanceof Error && err.message.includes('already exists and is not an empty directory')) {
            spawnSync("git", ["pull", "origin", repoInfo.branch], {stdio: 'inherit'});
            console.log(`Pulled ${repoUrl} to ${localRepoPath}. Branch: ${repoInfo.branch}`);
        } else {
            throw err;
        }
    }
}


const isWindows = process.platform === "win32";

export function getGithubRepoURL(owner: string, name: string) {
    return isWindows ? `https://github.com/${owner}/${name}` : `git@github.com:/${owner}/${name}.git`;
}
