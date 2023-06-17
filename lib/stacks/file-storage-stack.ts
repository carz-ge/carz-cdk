import { CfnOutput, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib'
import { Construct } from 'constructs'
import * as s3 from 'aws-cdk-lib/aws-s3'
import * as iam from 'aws-cdk-lib/aws-iam'
import {
	AllowedMethods,
	CachePolicy,
	Distribution,
	ViewerProtocolPolicy,
} from 'aws-cdk-lib/aws-cloudfront'
import { S3Origin } from 'aws-cdk-lib/aws-cloudfront-origins'

interface FileStorageStackProps extends StackProps {
	managerAuthenticatedRole: iam.IRole
	appAuthenticatedRole: iam.IRole
	allowedOrigins: string[]
}

export class FileStorageStack extends Stack {

	public readonly distributionDomain: string
	public readonly fileStageBucketName: string

	constructor(scope: Construct, id: string, props: FileStorageStackProps) {
		super(scope, id, props)

		const fileStorageBucket = new s3.Bucket(this, `CarzPublicBucket`, {
			removalPolicy: RemovalPolicy.DESTROY,
			autoDeleteObjects: true,
			cors: [
				{
					allowedMethods: [
						s3.HttpMethods.GET,
						s3.HttpMethods.POST,
						s3.HttpMethods.PUT,
						s3.HttpMethods.DELETE,
					],
					allowedOrigins: props.allowedOrigins,
					allowedHeaders: ['*'],
				},
			],
		})

		const imagesCFDistribution = new Distribution(this, `Distribution`, {
			defaultBehavior: {
				origin: new S3Origin(fileStorageBucket, { originPath: '/public' }),
				cachePolicy: CachePolicy.CACHING_OPTIMIZED,
				allowedMethods: AllowedMethods.ALLOW_GET_HEAD,
				cachedMethods: AllowedMethods.ALLOW_GET_HEAD,
				viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
			},
		})

		const canReadItemFromPublicDirectory = new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: ['s3:GetObject'],
			resources: [`arn:aws:s3:::${fileStorageBucket.bucketName}/public/*`],
		})

		const canReadItemFromProtectedDirectory = new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: ['s3:GetObject'],
			resources: [`arn:aws:s3:::${fileStorageBucket.bucketName}/protected/*`],
		})

		const canReadManyItemsFromCertainDirectories = new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: ['s3:ListBucket'],
			resources: [`arn:aws:s3:::${fileStorageBucket.bucketName}`],
			conditions: {
				StringLike: {
					's3:prefix': ['public/', 'public/*', 'protected/', 'protected/*'],
				},
			},
		})

		const canReadUpdateDeleteFromPublicDirectory = new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: ['s3:PutObject', 's3:GetObject', 's3:DeleteObject'],
			resources: [`arn:aws:s3:::${fileStorageBucket.bucketName}/public/*`],
		})

		const canReadUpdateDeleteFromOwnProtectedDirectory =
			new iam.PolicyStatement({
				effect: iam.Effect.ALLOW,
				actions: ['s3:PutObject', 's3:GetObject', 's3:DeleteObject'],
				resources: [
					`arn:aws:s3:::${fileStorageBucket.bucketName}/protected/\${cognito-identity.amazonaws.com:sub}/*`,
				],
			})

		const canReadUpdateDeleteFromOwnPrivateDirectory = new iam.PolicyStatement({
			effect: iam.Effect.ALLOW,
			actions: ['s3:PutObject', 's3:GetObject', 's3:DeleteObject'],
			resources: [
				`arn:aws:s3:::${fileStorageBucket.bucketName}/private/\${cognito-identity.amazonaws.com:sub}/*`,
			],
		})

		const canReadManyItemsFromPublicProtectedOwnPrivate =
			new iam.PolicyStatement({
				effect: iam.Effect.ALLOW,
				actions: ['s3:ListBucket'],
				resources: [`arn:aws:s3:::${fileStorageBucket.bucketName}`],
				conditions: {
					StringLike: {
						's3:prefix': [
							'public/',
							'public/*',
							'protected/',
							'protected/*',
							'private/${cognito-identity.amazonaws.com:sub}/',
							'private/${cognito-identity.amazonaws.com:sub}/*',
						],
					},
				},
			})

		const managedPolicyForAmplifyUnauth = new iam.ManagedPolicy(
			this,
			`ManagedPolicyForAmplifyApplication`,
			{
				description:
					'Managed policy to allow usage of Storage Library for unauth. For folks that are on the site, but not signed in.',
				statements: [
					canReadItemFromPublicDirectory,
					canReadItemFromProtectedDirectory,
					canReadManyItemsFromCertainDirectories,
				],
				roles: [props.appAuthenticatedRole],
			}
		)

		const managedPolicyForAmplifyAuth = new iam.ManagedPolicy(
			this,
			`ManagedPolicyForAmplifyManager`,
			{
				description:
					'Managed Policy to allow usage of storage library for auth when users are signed in.',
				statements: [
					canReadUpdateDeleteFromPublicDirectory,
					canReadUpdateDeleteFromOwnProtectedDirectory,
					canReadUpdateDeleteFromOwnPrivateDirectory,
					canReadItemFromProtectedDirectory,
					canReadManyItemsFromPublicProtectedOwnPrivate,
				],
				roles: [props.managerAuthenticatedRole],
			}
		)

		this.distributionDomain = imagesCFDistribution.domainName;
		this.fileStageBucketName = fileStorageBucket.bucketName;

		new CfnOutput(this, 'DistributionDomainName', {
			value: imagesCFDistribution.domainName,
		})

		new CfnOutput(this, 'BucketName', {
			value: fileStorageBucket.bucketName,
		})

		new CfnOutput(this, 'BucketRegion', {
			value: this.region,
		})
	}
}
