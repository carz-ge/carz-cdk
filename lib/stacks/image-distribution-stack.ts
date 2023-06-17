import {
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
    aws_iam as iam,
    aws_lambda as lambda,
    aws_logs as logs,
    aws_s3 as s3, CfnOutput,
    Duration,
    Environment,
    RemovalPolicy,
    Stack,
    StackProps
} from "aws-cdk-lib";
import {AwsRegion, Stage} from "../config/types";
import {Construct} from "constructs";
import {createS3Bucket} from "../core/create-s3-bucket";
import {createHash} from "crypto";
import {isProd} from "../utils/deployment-utils";
import {FunctionProps} from "aws-cdk-lib/aws-lambda/lib/function";
import {MyCustomResource} from "./my-custom-resource";
import {BehaviorOptions} from "aws-cdk-lib/aws-cloudfront/lib/distribution";


interface ImageDistributionStackProps extends StackProps {
    readonly env: Environment;
    readonly stage: Stage
}

// Parameters of S3 bucket where original images are stored
// CloudFront parameters
const CLOUDFRONT_ORIGIN_SHIELD_REGION = AwsRegion.EU_WEST_1;
// Parameters of transformed images
const S3_TRANSFORMED_IMAGE_EXPIRATION_DURATION = 90;

const S3_TRANSFORMED_IMAGE_CACHE_TTL = 'max-age=31622400';
// Lambda Parameters
const LAMBDA_MEMORY = '1500';
const LAMBDA_TIMEOUT = '60';
const LOG_TIMING = 'false';

type LambdaEnv = {
    originalImageBucketName: string,
    transformedImageBucketName?: any;
    transformedImageCacheTTL: string,
    secretKey: string,
    logTiming: string,
}


export class ImageDistributionStack extends Stack {

    constructor(scope: Construct, id: string, props: ImageDistributionStackProps) {
        super(scope, id, props);

        const imagesBucket = createS3Bucket(this, "carz-images-bucket-" + props.stage, {stage: props.stage});

        const transformedImageBucket = new s3.Bucket(this, 'carz-transformed-image-bucket-' + props.stage, {
            bucketName: 'carz-transformed-image-bucket-' + props.stage,
            removalPolicy:  isProd(props.stage) ? RemovalPolicy.RETAIN : RemovalPolicy.DESTROY,
            autoDeleteObjects: !isProd(props.stage),
            lifecycleRules: [
                {
                    expiration: Duration.days(isProd(props.stage) ? S3_TRANSFORMED_IMAGE_EXPIRATION_DURATION : 5),
                },
            ],
        });

        const secretKey = createHash('md5').update(this.node.addr).digest('hex');

        const imageProcessing = createImageProcessorLambda(this, imagesBucket.bucketName, secretKey, transformedImageBucket.bucketName);


        // Enable Lambda URL
        const imageProcessingURL = imageProcessing.addFunctionUrl({
            authType: lambda.FunctionUrlAuthType.NONE,
        });

        // Leverage a custom resource to get the hostname of the LambdaURL
        const imageProcessingHelper = new MyCustomResource(this, 'customResource', {
            Url: imageProcessingURL.url
        });

        const imageOrigin = new origins.OriginGroup({
            primaryOrigin: new origins.S3Origin(transformedImageBucket, {
                originShieldRegion: CLOUDFRONT_ORIGIN_SHIELD_REGION,
            }),
            fallbackOrigin: new origins.HttpOrigin(imageProcessingHelper.hostname, {
                originShieldRegion: CLOUDFRONT_ORIGIN_SHIELD_REGION,
                customHeaders: {
                    'x-origin-secret-header': secretKey,
                },
            }),
            fallbackStatusCodes: [403],
        });


        // Create a CloudFront Function for url rewrites
        const urlRewriteFunction = new cloudfront.Function(this, 'urlRewrite', {
            code: cloudfront.FunctionCode.fromFile({filePath: 'functions/url-rewrite/index.js',}),
            functionName: `urlRewriteFunction${this.node.addr}`,
        });
        // Creating a custom response headers policy. CORS allowed for all origins.
        const imageResponseHeadersPolicy = new cloudfront.ResponseHeadersPolicy(this, `ResponseHeadersPolicy${this.node.addr}`, {
            responseHeadersPolicyName: 'ImageResponsePolicy',
            corsBehavior: {
                accessControlAllowCredentials: false,
                accessControlAllowHeaders: ['*'],
                accessControlAllowMethods: ['GET'],
                accessControlAllowOrigins: ['*'],
                accessControlMaxAge: Duration.seconds(600),
                originOverride: false,
            },
            // recognizing image requests that were processed by this solution
            customHeadersBehavior: {
                customHeaders: [
                    {header: 'x-aws-image-optimization', value: 'v1.0', override: true},
                    {header: 'consty', value: 'accept', override: true},
                ],
            }
        });
        const imageDeliveryCacheBehaviorConfig: BehaviorOptions = {
            origin: imageOrigin,
            viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            cachePolicy: new cloudfront.CachePolicy(this, `ImageCachePolicy${this.node.addr}`, {
                defaultTtl: Duration.hours(24),
                maxTtl: Duration.days(365),
                minTtl: Duration.seconds(0),
                queryStringBehavior: cloudfront.CacheQueryStringBehavior.all()
            }),
            functionAssociations: [{
                eventType: cloudfront.FunctionEventType.VIEWER_REQUEST,
                function: urlRewriteFunction,
            }],
            responseHeadersPolicy: imageResponseHeadersPolicy
        };


        const imageDelivery = new cloudfront.Distribution(this, 'imageDeliveryDistribution', {
            comment: 'image optimization - image delivery',
            defaultBehavior: imageDeliveryCacheBehaviorConfig
        });

        new CfnOutput(this, 'ImageDeliveryDomain', {
            description: 'Domain name of image delivery',
            value: imageDelivery.distributionDomainName
        });
    }

}


function createImageProcessorLambda(parent: Construct, imagesBucketName: string,
                                    transformedImageBucketName: string,
                                    secretKey: string,) {
    const lambdaEnv: LambdaEnv = {
        originalImageBucketName: imagesBucketName,
        transformedImageCacheTTL: S3_TRANSFORMED_IMAGE_CACHE_TTL,
        secretKey: secretKey,
        logTiming: LOG_TIMING,
        transformedImageBucketName: transformedImageBucketName
    };

    // IAM policy to read from the S3 bucket containing the original images
    const s3ReadOriginalImagesPolicy = new iam.PolicyStatement({
        actions: ['s3:GetObject', 's3:ListBucket'],
        resources: ['arn:aws:s3:::' + imagesBucketName + '/*'],
    });

    // attach iam policy to the role assumed by Lambda
    // write policy for Lambda on the s3 bucket for transformed images
    const s3WriteTransformedImagesPolicy = new iam.PolicyStatement({
        actions: ['s3:PutObject'],
        resources: ['arn:aws:s3:::' + transformedImageBucketName + '/*'],
    });

    // Create Lambda for image processing
    const lambdaProps: FunctionProps = {
        runtime: lambda.Runtime.NODEJS_16_X,
        handler: 'index.handler',
        code: lambda.Code.fromAsset('functions/image-processing'),
        timeout: Duration.seconds(parseInt(LAMBDA_TIMEOUT)),
        memorySize: parseInt(LAMBDA_MEMORY),
        environment: lambdaEnv,
        logRetention: logs.RetentionDays.ONE_DAY,
        initialPolicy: [
            s3ReadOriginalImagesPolicy,
            s3WriteTransformedImagesPolicy
        ]
    };


    return new lambda.Function(parent, 'image-optimization', lambdaProps);
}
