import {Stage} from "../config/types";

export function isProd(stage: Stage) {
    return stage === Stage.PROD;
}
