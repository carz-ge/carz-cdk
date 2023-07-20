#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {ImageDistributionStack} from "../lib/stacks/image-distribution-stack";
import {config} from '../lib/config/cdk-config';
import {Environment} from "aws-cdk-lib";
import {S3CloudfrontStack} from "../lib/stacks/s3-cloudfront-stack";
import {ScheduledLambdaStack} from "../lib/stacks/scheduled-lambda-stack";

const app = new cdk.App();

const environment: Environment = {
    account: config.awsAccount,
    region: config.deploymentConfigs[0].region,
};

new S3CloudfrontStack(app, 'CarzImageTransformationStack', {
    env: environment,
    stage: config.deploymentConfigs[0].stage,
});


new ScheduledLambdaStack(app, "ScheduledLambdaStack", {
    env: environment,
    stage: config.deploymentConfigs[0].stage,
    serviceName: "stations-fetcher"
})