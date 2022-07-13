import { AxiosRequestConfig } from "axios";
import apiClient from ".";
import { IDirectoryInfo, IFileInfo } from "./interfaces";

export interface IGetDirSturectureRequest {
    DirId: number;
    PageSize: number;
    PageOffset: number;
}
export interface IGetDirSturectureResponse {
    Info: IDirectoryInfo;
    SubDirs: IDirectoryInfo[];
    SubFiles: IFileInfo[];
}
export async function getDirStructure(req: IGetDirSturectureRequest, config?: AxiosRequestConfig): Promise<IGetDirSturectureResponse> {
    const {DirId, PageOffset, PageSize} = req;
    const resp = await apiClient.get(`/api/dir/${DirId}?page_index=${PageOffset}&page_size=${PageSize}`, config);
    if (resp.status === 200) {
        return resp.data;
    } else {
        if (resp.data?.message) {
            throw new Error(resp.data.message)
        } else {
            throw new Error('Network Error')
        }
    }
}