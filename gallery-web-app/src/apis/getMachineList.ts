import apiClient from ".";
import { IRootDirectoryInfo } from "./interfaces";

export interface IGetMachineListRequest {

}
export interface IGetMachineListResponse {
    machines: IRootDirectoryInfo[];
}

export async function getMachineList(req?: IGetMachineListRequest): Promise<IGetMachineListResponse> {
    const resp = await apiClient.get("/api/machines");
    if (resp.status === 200) {
        return resp.data;
    } else if (resp.data?.message){
        throw new Error (resp.data.message);
    } else {
        throw new Error('Network Error')
    }
}