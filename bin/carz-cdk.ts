#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {ImageDistributionStack} from "../lib/stacks/image-distribution-stack";
import {config} from '../lib/config/cdk-config';
import {Environment} from "aws-cdk-lib";

const app = new cdk.App();

const environment: Environment = {
    account: config.awsAccount,
    region: config.deploymentConfigs[0].region,
};

new ImageDistributionStack(app, 'CarzImageTransformationStack', {
    env: environment,
    stage: config.deploymentConfigs[0].stage,
});
