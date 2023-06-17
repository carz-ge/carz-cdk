import {CfnOutput, Stack, RemovalPolicy} from "aws-cdk-lib";
import {BlockPublicAccess, Bucket, BucketEncryption} from "aws-cdk-lib/aws-s3";
import {Stage} from "../config/types";
import {isProd} from "../utils/deployment-utils";


export interface CreateBucketProps {
    stage: Stage;
}

export function createS3Bucket(stack: Stack, bucketName: string, props: CreateBucketProps): Bucket {
    const bucket = new Bucket(stack, bucketName, {
        bucketName: bucketName,
        versioned: true,
        encryption: BucketEncryption.S3_MANAGED,
        blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
        enforceSSL: true,
        removalPolicy: isProd(props.stage) ? RemovalPolicy.RETAIN : RemovalPolicy.DESTROY,
        // lifecycleRules: [
        //     {
        //         enabled: true,
        //         expiration: props.stage === Stage.PROD ? undefined : Duration.days(1),
        //         transitions: [
        //             {
        //                 storageClass: StorageClass.INFREQUENT_ACCESS,
        //                 transitionAfter: Duration.days(30),
        //             },
        //         ],
        //     },
        // ],

    });
    new CfnOutput(stack, 'S3BucketName' + bucketName, {
        value: bucket.bucketName,
        description: `S3 bucket name for service ${bucketName}`,
        exportName: 'S3BucketName' + bucketName
    })

    return bucket;
}