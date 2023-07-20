import {config} from "../lib/config/cdk-config";
import * as path from "path";
import {cloneOrPullRepoSync} from "../lib/utils/github-utils";


config.deploymentConfigs.forEach(deploymentConfig => {
    const repoInfo = deploymentConfig.githubRepository;
    if (!repoInfo) {
        return;
    }
    const imageDirectory  =  path.join(__dirname, "../service-repos", repoInfo.name, repoInfo.branch);
    cloneOrPullRepoSync(repoInfo, imageDirectory);
});
