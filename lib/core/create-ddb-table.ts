import {CfnOutput, RemovalPolicy, Stack} from "aws-cdk-lib";
import {AttributeType, BillingMode, Table} from "aws-cdk-lib/aws-dynamodb";
import {Stage} from "../config/types";
import {isProd} from "../utils/deployment-utils";


interface CreateDdbTableProps {
    tableName: string;
    partitionKey: string;
    sortKey?: string;
    stage: Stage;
}

export function createDdbTable(stack: Stack, props: CreateDdbTableProps): Table {
    const table = new Table(stack, props.tableName, {
        tableName: props.tableName,
        partitionKey: {
            name: props.partitionKey,
            type: AttributeType.STRING,
        },
        sortKey: props.sortKey ? {
            name: props.sortKey,
            type: AttributeType.STRING,
        } : undefined,
        billingMode: BillingMode.PAY_PER_REQUEST,
        removalPolicy: isProd(props.stage) ? RemovalPolicy.RETAIN : RemovalPolicy.DESTROY,
    })

    new CfnOutput(stack, 'DdbTableName' + props.tableName, {
        value: props.tableName,
        description: `DynamoDB table name for service ${props.tableName}`,
        exportName: 'DdbTableName' + props.tableName

    })
    return table;
}