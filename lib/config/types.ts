export enum AwsRegion {
    US_EAST_1 = "us-east-1",
    US_EAST_2 = "us-east-2",
    EU_WEST_1 = "eu-west-1"
}

export enum Stage {
    PROD = "prod",
    DEV = "dev",
}

export type DeploymentConfig = {
    region: AwsRegion
    stage: Stage
    githubRepository?: GithubRepository
}

export type GithubRepository = {
    name: string
    owner: string
    branch: string
}

export type CdkConfiguration = {
    applicationName: string
    deploymentConfigs: DeploymentConfig[]
    awsAccount: string
}
