import {Duration, Environment, Stack, StackProps} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {Effect, Policy, PolicyStatement, Role, ServicePrincipal} from 'aws-cdk-lib/aws-iam';
import {Architecture, DockerImageCode, DockerImageFunction, Tracing} from 'aws-cdk-lib/aws-lambda';
import * as path from 'path';
import {Stage} from "../config/types";
import {RetentionDays} from "aws-cdk-lib/aws-logs";
import {ISecret, Secret} from "aws-cdk-lib/aws-secretsmanager";
import {Rule, Schedule} from "aws-cdk-lib/aws-events";
import {LambdaFunction} from "aws-cdk-lib/aws-events-targets";

interface ScheduledLambdaStackProps extends StackProps {
    readonly env: Environment;
    readonly stage: Stage;
    readonly serviceName: string;

}

const getValueFromSecret = (secret: ISecret, key: string): string => {
    return secret.secretValueFromJson(key).unsafeUnwrap()
}

export class ScheduledLambdaStack extends Stack {
    constructor(scope: Construct, id: string, props: ScheduledLambdaStackProps) {
        super(scope, id, props);


        const functionName = `${props.serviceName}-${props.stage.toLowerCase()}-${props.env.region?.toLowerCase()}`;
        console.log("function name: ", functionName);

        const scheduledFunction = new DockerImageFunction(this, functionName, {
            code: DockerImageCode.fromImageAsset(path.join(__dirname,
                '..',
                '..',
                'functions',
                'stations-data-fetcher',
            )),
            memorySize: 1024,
            timeout: Duration.seconds(5),
            environment: {
                "STAGE": props.stage,
            },
            tracing: Tracing.ACTIVE,
            architecture: Architecture.ARM_64,
            logRetention: RetentionDays.THREE_DAYS,
        });

        const secret = Secret.fromSecretCompleteArn(this, `secret-${functionName}`,  "arn:aws:secretsmanager:eu-west-1:907239669915:secret:stations-fetcher-prod-eu-west-1-UNmVWw")
        secret.grantRead(scheduledFunction.role!)

        // need to create role and policy for scheduler to invoke the lambda function
        const schedulerRole = new Role(this, `scheduler-role-${functionName}`, {
            assumedBy: new ServicePrincipal('scheduler.amazonaws.com'),
        });

        new Policy(this, `schedule-policy-${functionName}`, {
            policyName: 'ScheduleToInvokeLambdas',
            roles: [schedulerRole],
            statements: [
                new PolicyStatement({
                    effect: Effect.ALLOW,
                    actions: ['lambda:InvokeFunction'],
                    resources: [scheduledFunction.functionArn],
                }),
            ],
        });


        new Rule(this, `scheduler-rule-${functionName}`, {
            schedule: Schedule.cron({
                year: "*",
                month: "*",
                day: "*",
                hour: "16",
                minute: "30",
            }),
            targets: [new LambdaFunction(scheduledFunction)],
        });

    //     // Create a group for the schedule (maybe you want to add more scheudles to this group the future?)
    //     const group = new CfnScheduleGroup(this, `schedule-group-${functionName}`, {
    //         name: 'SchedulesForLambda',
    //     });
    //
    //
    //     // Creates the schedule to invoke every 5 minutes
    //     new CfnSchedule(this, `schedule-cfn-${functionName}`, {
    //         groupName: group.name,
    //         flexibleTimeWindow: {
    //             maximumWindowInMinutes: 5,
    //             mode: 'FLEXIBLE',
    //         },
    //         scheduleExpression: 'cron(30 16 * * ? *)',
    //         target: {
    //             arn: scheduledFunction.functionArn,
    //             roleArn: schedulerRole.roleArn,
    //         },
    //     });
    }
}