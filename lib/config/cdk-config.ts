import {AwsRegion, CdkConfiguration, Stage} from "./types";

export const config: CdkConfiguration = {
    applicationName: process.env.AWS_APP_NAME || "CarzInfra",
    awsAccount: process.env.AWS_ACCOUNT || "907239669915",
    deploymentConfigs: [
        // {
        //     region: AwsRegion.US_EAST_2,
        //     stage: Stage.DEV,
        // },
        {
            region: AwsRegion.EU_WEST_1,
            stage: Stage.PROD,
        }
    ],

}
