import {ResourcesType} from "@/containers/devices/Devices.types";

export type DpsDataType = {
    deviceId: string | undefined
    resources: ResourcesType[] | undefined
}