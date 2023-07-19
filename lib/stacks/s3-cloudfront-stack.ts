import {Environment, Stack, StackProps} from "aws-cdk-lib";
import {Construct} from "constructs";
import {CloudFrontToS3, CloudFrontToS3Props} from "@aws-solutions-constructs/aws-cloudfront-s3";
import {Stage} from "../config/types";
import {createS3Bucket} from "../core/create-s3-bucket";


interface S3CloudfrontStackProps extends StackProps {
    readonly env: Environment;
    readonly stage: Stage
}

export class S3CloudfrontStack extends Stack {
    constructor(scope: Construct, id: string, props: S3CloudfrontStackProps) {
        super(scope, id, props);

        const cloudFrontToS3Props: CloudFrontToS3Props = {
            // existingBucketObj: createS3Bucket(this, "carz-images-pros", props)
        }
        const cloudFrontToS3 = new CloudFrontToS3(this, 'carz-cloudfront-s3', cloudFrontToS3Props);

    }
}
